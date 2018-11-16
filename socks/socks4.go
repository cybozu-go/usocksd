package socks

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
)

func (s *Server) handleSOCKS4(ctx context.Context, conn net.Conn, cmdByte byte) net.Conn {
	var responseData [8]byte
	responseData[1] = byte(Status4Rejected)
	fields := well.FieldsFromContext(ctx)
	fields[log.FnType] = "access"
	fields[log.FnProtocol] = SOCKS4.String()
	fields["client_addr"] = conn.RemoteAddr().String()

	errFunc := func(msg string, err error) net.Conn {
		conn.Write(responseData[:])
		if err != nil {
			fields[log.FnError] = err.Error()
		}
		s.Logger.Error(msg, fields)
		return nil
	}

	command := commandType(cmdByte)
	fields["command"] = command.String()
	if command != CmdConnect {
		return errFunc("command not supported", nil)
	}

	var payload [6]byte
	_, err := io.ReadFull(conn, payload[:])
	if err != nil {
		return errFunc("failed to read port/ip", err)
	}

	port := int(binary.BigEndian.Uint16(payload[0:2]))
	ip := payload[2:6]
	socks4a := (ip[0] == 0 && ip[1] == 0 && ip[2] == 0 && ip[3] != 0)
	username, err := readUntilNull(conn)
	if err != nil {
		return errFunc("failed to read username", err)
	}
	r := &Request{
		Version:  SOCKS4,
		Command:  command,
		Port:     port,
		Username: username,
		Conn:     conn,
		ctx:      ctx,
	}
	if socks4a {
		hostname, err := readUntilNull(conn)
		if err != nil {
			return errFunc("failed to read hostname", err)
		}
		r.Hostname = hostname
		fields["dest_host"] = hostname
	} else {
		r.IP = net.IPv4(ip[0], ip[1], ip[2], ip[3])
		fields["dest_host"] = r.IP.String()
	}

	if s.Auth != nil && !s.Auth.Authenticate(r) {
		return errFunc("authentication failure", nil)
	}

	if s.Rules != nil && !s.Rules.Match(r) {
		return errFunc("ruleset mismatch", nil)
	}

	destConn, err := s.dial(ctx, r, "tcp4")
	if err != nil {
		return errFunc("dial to destination failed", err)
	}

	responseData[1] = byte(Status4Granted)
	copy(responseData[2:8], payload[:])

	_, err = conn.Write(responseData[:])
	if err != nil {
		destConn.Close()
		return errFunc("failed to write response", err)
	}

	fields["dest_addr"] = destConn.RemoteAddr().String()
	fields["src_addr"] = destConn.LocalAddr().String()
	s.Logger.Info("proxy starts", fields)
	return destConn
}

func readUntilNull(conn net.Conn) (string, error) {
	var buf []byte
	var data [1]byte

	for {
		_, err := conn.Read(data[:])
		if err != nil {
			return "", err
		}
		if data[0] == 0 {
			return string(buf), nil
		}
		buf = append(buf, data[0])
	}
}

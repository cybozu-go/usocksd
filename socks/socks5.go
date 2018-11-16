package socks

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/cybozu-go/log"
	"github.com/cybozu-go/netutil"
	"github.com/cybozu-go/well"
)

func (s *Server) handleSOCKS5(ctx context.Context, conn net.Conn, nauth byte) net.Conn {
	r := &Request{
		Version: SOCKS5,
		Conn:    conn,
		ctx:     ctx,
	}
	if !s.negotiateAuth(r, int(nauth)) {
		return nil
	}
	if !s.readAddress(r) {
		return nil
	}

	if s.Logger.Enabled(log.LvDebug) {
		s.Logger.Debug("request info", map[string]interface{}{
			"request": r,
		})
	}

	response := makeSOCKS5Response(r)
	fields := well.FieldsFromContext(ctx)
	fields[log.FnType] = "access"
	fields[log.FnProtocol] = SOCKS5.String()
	fields["client_addr"] = conn.RemoteAddr().String()
	fields["command"] = r.Command.String()
	if len(r.Hostname) > 0 {
		fields["dest_host"] = r.Hostname
	} else {
		fields["dest_host"] = r.IP.String()
	}

	errFunc := func(msg string) net.Conn {
		conn.Write(response)
		s.Logger.Error(msg, fields)
		return nil
	}

	if r.Command != CmdConnect {
		response[1] = byte(Status5CommandNotSupported)
		return errFunc("command not supported")
	}

	if s.Rules != nil {
		if !s.Rules.Match(r) {
			response[1] = byte(Status5DeniedByRuleset)
			return errFunc("ruleset mismatch")
		}
	}

	destConn, err := s.dial(ctx, r, "tcp")
	if err != nil {
		fields[log.FnError] = err.Error()
		switch {
		case netutil.IsNetworkUnreachable(err):
			response[1] = byte(Status5NetworkUnreachable)
		case netutil.IsConnectionRefused(err):
			response[1] = byte(Status5ConnectionRefused)
		case netutil.IsNoRouteToHost(err):
			response[1] = byte(Status5HostUnreachable)
		}
		return errFunc("dial to destination failed")
	}

	response[1] = byte(Status5Granted)
	_, err = conn.Write(response)
	if err != nil {
		destConn.Close()
		fields[log.FnError] = err.Error()
		return errFunc("failed to write response")
	}

	fields["dest_addr"] = destConn.RemoteAddr().String()
	fields["src_addr"] = destConn.LocalAddr().String()
	s.Logger.Info("proxy starts", fields)
	return destConn
}

func hasAuth(t authType, methods []byte) bool {
	for _, m := range methods {
		if t == authType(m) {
			return true
		}
	}
	return false
}

func (s *Server) negotiateAuth(r *Request, nauth int) bool {
	logError := func(msg string, err error) {
		fields := well.FieldsFromContext(r.ctx)
		fields[log.FnType] = "access"
		fields[log.FnProtocol] = SOCKS5.String()
		fields["client_addr"] = r.Conn.RemoteAddr().String()
		if err != nil {
			fields[log.FnError] = err.Error()
		}
		s.Logger.Error(msg, fields)
	}

	methods := make([]byte, nauth)
	_, err := io.ReadFull(r.Conn, methods)
	if err != nil {
		logError("failed to read auth methods", err)
		return false
	}

	var response [2]byte
	response[0] = byte(SOCKS5)

	if hasAuth(AuthBasic, methods) {
		response[1] = byte(AuthBasic)
		_, err = r.Conn.Write(response[:])
		if err != nil {
			logError("failed to negotiate auth method", err)
			return false
		}

		// basic auth response (failure)
		response[0] = 0x01
		response[1] = 0xff

		var preamble [2]byte
		_, err = io.ReadFull(r.Conn, preamble[:])
		if err != nil {
			r.Conn.Write(response[:])
			logError("failed to read username/password", err)
			return false
		}
		if preamble[0] != 0x01 {
			r.Conn.Write(response[:])
			logError("invalid auth version", nil)
			return false
		}

		usernameLength := int(preamble[1])
		if usernameLength > 0 {
			username := make([]byte, usernameLength)
			_, err := io.ReadFull(r.Conn, username)
			if err != nil {
				r.Conn.Write(response[:])
				logError("failed to read username", err)
				return false
			}
			r.Username = string(username)
		}

		var oneByte [1]byte
		_, err = io.ReadFull(r.Conn, oneByte[:])
		if err != nil {
			r.Conn.Write(response[:])
			logError("failed to read password length", err)
			return false
		}
		if oneByte[0] > 0 {
			password := make([]byte, oneByte[0])
			_, err := io.ReadFull(r.Conn, password)
			if err != nil {
				r.Conn.Write(response[:])
				logError("failed to read password", err)
				return false
			}
			r.Password = string(password)
		}

		if s.Auth != nil && !s.Auth.Authenticate(r) {
			r.Conn.Write(response[:])
			logError("authentication failure", nil)
			return false
		}

		response[1] = 0x00
		_, err = r.Conn.Write(response[:])
		if err != nil {
			logError("failed to write authentication result", err)
			return false
		}

		return true
	}

	if hasAuth(AuthNo, methods) {
		if s.Auth != nil && !s.Auth.Authenticate(r) {
			// No authentication method still need to be checked
			// by s.Auth if given.
			goto FAIL
		}

		response[1] = byte(AuthNo)
		_, err := r.Conn.Write(response[:])
		if err != nil {
			logError("failed to negotiate auth method", err)
			return false
		}
		return true
	}

	// unacceptable authentication method
FAIL:
	response[1] = 0xff
	r.Conn.Write(response[:])
	logError("no acceptable auth methods", nil)
	return false
}

func (s *Server) readAddress(r *Request) bool {
	logError := func(msg string, err error) {
		fields := well.FieldsFromContext(r.ctx)
		fields[log.FnType] = "access"
		fields[log.FnProtocol] = SOCKS5.String()
		fields["client_addr"] = r.Conn.RemoteAddr().String()
		if err != nil {
			fields[log.FnError] = err.Error()
		}
		s.Logger.Error(msg, fields)
	}

	var addrData [4]byte
	_, err := io.ReadFull(r.Conn, addrData[:])
	if err != nil {
		logError("failed to read address", err)
		return false
	}
	if version(addrData[0]) != SOCKS5 {
		logError("request is not SOCKS5", nil)
		return false
	}
	r.Command = commandType(addrData[1])

	switch addressType(addrData[3]) {
	case AddrIPv4:
		var ipv4Data [4]byte
		_, err := io.ReadFull(r.Conn, ipv4Data[:])
		if err != nil {
			logError("failed to read address", err)
			return false
		}
		r.IP = net.IPv4(ipv4Data[0], ipv4Data[1], ipv4Data[2], ipv4Data[3])
	case AddrIPv6:
		ipv6Data := make([]byte, 16)
		_, err := io.ReadFull(r.Conn, ipv6Data)
		if err != nil {
			logError("failed to read address", err)
			return false
		}
		r.IP = net.IP(ipv6Data)
	case AddrDomain:
		var nameLen [1]byte
		_, err := io.ReadFull(r.Conn, nameLen[:])
		if err != nil {
			logError("failed to read address", err)
			return false
		}
		if nameLen[0] > 0 {
			name := make([]byte, nameLen[0])
			_, err := io.ReadFull(r.Conn, name)
			if err != nil {
				logError("failed to read address", err)
				return false
			}
			r.Hostname = string(name)
		}
	default:
		logError("unknown address type", nil)
		return false
	}

	var portData [2]byte
	_, err = io.ReadFull(r.Conn, portData[:])
	if err != nil {
		logError("failed to read port number", err)
		return false
	}
	r.Port = int(binary.BigEndian.Uint16(portData[:]))

	return true
}

func makeSOCKS5Response(r *Request) []byte {
	responseLen := 6
	var at addressType
	switch {
	case len(r.Hostname) > 0:
		responseLen += len(r.Hostname) + 1
		at = AddrDomain
	case r.IP.To4() == nil:
		responseLen += 16
		at = AddrIPv6
	default:
		responseLen += 4
		at = AddrIPv4
	}
	response := make([]byte, responseLen)
	response[0] = byte(SOCKS5)
	response[1] = byte(Status5Failure)
	response[3] = byte(at)
	switch at {
	case AddrDomain:
		response[4] = byte(len(r.Hostname))
		copy(response[5:5+len(r.Hostname)], r.Hostname)
	case AddrIPv4:
		ipv4 := r.IP.To4()
		copy(response[4:8], []byte(ipv4))
	case AddrIPv6:
		ipv6 := r.IP.To16()
		copy(response[4:20], []byte(ipv6))
	}
	binary.BigEndian.PutUint16(response[len(response)-2:], uint16(r.Port))
	return response
}

package socks

import (
	"context"
	"net"
)

// Request is a struct to represent a request from SOCKS client.
//
// Authenticator, RuleSet, and Dialer can use Context and SetContext
// to associate any value with the request, and to cancel lengthy
// operations.
type Request struct {
	// Version is either SOCKS4 or SOCKS5
	Version version

	// Hostname is the destination DNS hostname.
	// If this is empty, IP is set to the destination address.
	Hostname string

	// Command is the requested command.
	Command commandType

	// IP is the destination IP address.
	// This may not be set if Hostname is not empty.
	IP net.IP

	// Port is the destination port number.
	Port int

	// Username is user name string for authentication.
	// Username may be empty when no authencation is requested.
	Username string

	// Password is password string for authentication.
	Password string

	// Conn is the connection from the client.
	Conn net.Conn

	ctx context.Context
}

// Context returns the request context.
func (r *Request) Context() context.Context {
	return r.ctx
}

// SetContext sets the request context.
func (r *Request) SetContext(ctx context.Context) {
	r.ctx = ctx
}

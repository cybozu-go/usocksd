package socks

type version byte

// SOCKS versions.
const (
	SOCKS4 = version(0x04)
	SOCKS5 = version(0x05)
)

func (v version) String() string {
	switch v {
	case SOCKS4:
		return "SOCKS4/4a"
	case SOCKS5:
		return "SOCKS5"
	}
	return ""
}

type commandType byte

// SOCKS commands.
const (
	CmdConnect = commandType(0x01)
	CmdBind    = commandType(0x02)
	CmdUDP     = commandType(0x03)
)

func (c commandType) String() string {
	switch c {
	case CmdConnect:
		return "connect"
	case CmdBind:
		return "bind"
	case CmdUDP:
		return "UDP associate"
	}
	return ""
}

type addressType byte

// SOCKS address types.
const (
	AddrIPv4   = addressType(0x01)
	AddrDomain = addressType(0x03)
	AddrIPv6   = addressType(0x04)
)

func (at addressType) String() string {
	switch at {
	case AddrIPv4:
		return "IPv4"
	case AddrDomain:
		return "Domain name"
	case AddrIPv6:
		return "IPv6"
	}
	return ""
}

type authType byte

// SOCKS authentication types.
const (
	AuthNo     = authType(0x00)
	AuthGSSAPI = authType(0x01)
	AuthBasic  = authType(0x02)
)

func (a authType) String() string {
	switch a {
	case AuthNo:
		return "no auth"
	case AuthGSSAPI:
		return "GSSAPI"
	case AuthBasic:
		return "basic"
	}
	return ""
}

type socks4ResponseStatus byte

// SOCKS4 response status codes.
const (
	Status4Granted     = socks4ResponseStatus(0x5a)
	Status4Rejected    = socks4ResponseStatus(0x5b)
	Status4NoIdentd    = socks4ResponseStatus(0x5c)
	Status4InvalidUser = socks4ResponseStatus(0x5d)
)

func (s socks4ResponseStatus) String() string {
	switch s {
	case Status4Granted:
		return "granted"
	case Status4Rejected:
		return "rejected"
	case Status4NoIdentd:
		return "no identd"
	case Status4InvalidUser:
		return "invalid user"
	}
	return ""
}

type socks5ResponseStatus byte

// SOCKS5 response status codes.
const (
	Status5Granted             = socks5ResponseStatus(0x00)
	Status5Failure             = socks5ResponseStatus(0x01)
	Status5DeniedByRuleset     = socks5ResponseStatus(0x02)
	Status5NetworkUnreachable  = socks5ResponseStatus(0x03)
	Status5HostUnreachable     = socks5ResponseStatus(0x04)
	Status5ConnectionRefused   = socks5ResponseStatus(0x05)
	Status5TTLExpired          = socks5ResponseStatus(0x06)
	Status5CommandNotSupported = socks5ResponseStatus(0x07)
	Status5AddressNotSupported = socks5ResponseStatus(0x08)
)

func (s socks5ResponseStatus) String() string {
	switch s {
	case Status5Granted:
		return "granted"
	case Status5Failure:
		return "failure"
	case Status5DeniedByRuleset:
		return "now allowed"
	case Status5NetworkUnreachable:
		return "network unreachable"
	case Status5HostUnreachable:
		return "host unreachable"
	case Status5ConnectionRefused:
		return "connection refused"
	case Status5TTLExpired:
		return "TTL expired"
	case Status5CommandNotSupported:
		return "command not supported"
	case Status5AddressNotSupported:
		return "address type not supported"
	}
	return ""
}

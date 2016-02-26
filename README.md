[![GoDoc](https://godoc.org/github.com/cybozu-go/usocksd?status.png)][godoc]
[![Build Status](https://travis-ci.org/cybozu-go/usocksd.png)](https://travis-ci.org/cybozu-go/usocksd)

Micro SOCKS5 server
===================

`usocksd` is a SOCKS5 server written in Go.
It is based on [armon/go-socks5][armon] server framework.

Usage
-----

`usocksd [-h] [-f CONFIG]`

The default configuration file path is `/usr/local/etc/usocksd.toml`.

`usocksd` does not have *daemon* mode.  Use systemd or upstart to
run it on your background.

Install
-------

Use Go 1.5 or better.

```
go get github.com/cybozu-go/usocksd
go install github.com/cybozu-go/usocksd/cmd/usocksd
```

Configuration file format
-------------------------

`usocksd.toml` is a [TOML][] file.
All fields are optional.

```
[log]
file = "/path/to/file"
level = "info"                     # crit, error, warn, info, debug

[incoming]
port = 1080
addresses = ["127.0.0.1"]          # List of listening IP addresses
allow_from = ["10.0.0.0/8"]        # CIDR network or IP address

[outgoing]
allow_sites = [                    # List of FQDN to be granted.
    "www.amazon.com",              # exact match
    ".google.com",                 # subdomain match
]
deny_sites = [                     # List of FQDN to be denied.
    ".2ch.net",                    # subdomain match
    "bad.google.com",              # deny a domain of *.google.com
    "",                            # "" matches non-FQDN (IP) requests.
]
deny_ports = [22, 25]              # Black list of outbound ports
addresses = ["12.34.56.78"]        # List of source IP addresses
dnsbl_domain = "some.dnsbl.org"    # to exclude black listed IP addresses
```

License
-------

[MIT](https://opensource.org/licenses/MIT)

[armon]: https://github.com/armon/go-socks5/
[TOML]: https://github.com/toml-lang/toml
[godoc]: https://godoc.org/github.com/cybozu-go/usocksd

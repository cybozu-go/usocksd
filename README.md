[![GoDoc](https://godoc.org/github.com/cybozu-go/usocksd?status.png)][godoc]
[![Build Status](https://travis-ci.org/cybozu-go/usocksd.png)](https://travis-ci.org/cybozu-go/usocksd)

Micro SOCKS5 server
===================

**usocksd** is a SOCKS5 server written in Go.
It is based on [armon/go-socks5][armon] server framework.

Features
--------

* Multiple external IP addresses

    usocksd can be configured to use multiple external IP addresses
    for outgoing connections.

    usocksd keeps using the same external IP address for a client
    as much as possible.  This means usocksd can proxy passive FTP
    connections reliably.

    Moreover, you can use a [DNSBL][] service to exclude dynamically
    from using some undesirable external IP addresses.

* White- and black- list of sites

    usocksd can be configured to grant access to the sites listed
    in a white list, and/or to deny access to the sites listed in a
    black list.

    usocksd can block connections to specific TCP ports, too.

Usage
-----

`usocksd [-h] [-f CONFIG]`

The default configuration file path is `/usr/local/etc/usocksd.toml`.

usocksd does not have *daemon* mode.  Use systemd or upstart to
run it on your background.

Install
-------

Use Go 1.5 or better.

```
go get github.com/cybozu-go/usocksd/cmd/usocksd
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

Tuning
------

If you see usocksd consumes too much CPU, try setting [`GOGC`][GOGC] to higher value, say **300**.

License
-------

[MIT](https://opensource.org/licenses/MIT)

[armon]: https://github.com/armon/go-socks5/
[DNSBL]: https://en.wikipedia.org/wiki/DNSBL
[TOML]: https://github.com/toml-lang/toml
[godoc]: https://godoc.org/github.com/cybozu-go/usocksd
[GOGC]: https://golang.org/pkg/runtime/#pkg-overview

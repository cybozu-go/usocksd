[![GitHub release](https://img.shields.io/github/release/cybozu-go/usocksd.svg?maxAge=60)][releases]
[![GoDoc](https://godoc.org/github.com/cybozu-go/usocksd?status.svg)][godoc]
[![CircleCI](https://circleci.com/gh/cybozu-go/usocksd.svg?style=svg)](https://circleci.com/gh/cybozu-go/usocksd)
[![Go Report Card](https://goreportcard.com/badge/github.com/cybozu-go/usocksd)](https://goreportcard.com/report/github.com/cybozu-go/usocksd)
[![License](https://img.shields.io/github/license/cybozu-go/usocksd.svg?maxAge=2592000)](LICENSE)

Micro SOCKS server
==================

**usocksd** is a SOCKS server written in Go.

[`usocksd/socks`](https://godoc.org/github.com/cybozu-go/usocksd/socks)
is a general purpose SOCKS server library.  usocksd is built on it.

Features
--------

* Support for SOCKS4, SOCKS4a, SOCK5

    * Only CONNECT is supported (BIND and UDP associate is missing).

* Graceful stop & restart

    * On SIGINT/SIGTERM, usocksd stops gracefully.
    * On SIGHUP, usocksd restarts gracefully.

* Access log

    Thanks to [`cybozu-go/log`](https://github.com/cybozu-go/log),
    usocksd can output access logs in structured formats including
    JSON.

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

Install
-------

Use Go 1.7 or better.

```
go get -u github.com/cybozu-go/usocksd/...
```

Usage
-----

`usocksd [-h] [-f CONFIG]`

The default configuration file path is `/etc/usocksd.toml`.

In addition, `usocksd` implements [the common spec](https://github.com/cybozu-go/well#specifications) from [`cybozu-go/well`](https://github.com/cybozu-go/well).

usocksd does not have *daemon* mode.  Use systemd to run it on your background.

Configuration file format
-------------------------

`usocksd.toml` is a [TOML][] file.
All fields are optional.

```
[log]
filename = "/path/to/file"         # default to stderr.
level = "info"                     # critical, error, warning, info, debug
format = "plain"                   # plain, logfmt, json

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

[releases]: https://github.com/cybozu-go/usocksd/releases
[DNSBL]: https://en.wikipedia.org/wiki/DNSBL
[TOML]: https://github.com/toml-lang/toml
[godoc]: https://godoc.org/github.com/cybozu-go/usocksd
[GOGC]: https://golang.org/pkg/runtime/#pkg-overview

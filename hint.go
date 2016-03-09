package usocksd

import (
	"hash/fnv"

	socks5 "github.com/cybozu-go/go-socks5"
	"golang.org/x/net/context"
)

type rewrite struct{}

// add context a hint for choosing an outgoing IP address.
func (r rewrite) Rewrite(ctx context.Context, req *socks5.Request) (context.Context, *socks5.AddrSpec) {
	hash := fnv.New32a()
	hash.Write(req.RemoteAddr.IP)
	hash.Write(req.DestAddr.IP)
	return context.WithValue(ctx, "hint", hash.Sum32()), req.DestAddr
}

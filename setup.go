package vpnrouting

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("vpnrouting", setup) }

func setup(c *caddy.Controller) error {
	routing := New()
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "address":
				if !c.NextArg() {
					return c.ArgErr()
				}
				routing.ResolverIP = c.Val()
			case "port":
				if !c.NextArg() {
					return c.ArgErr()
				}
				val, err := strconv.Atoi(c.Val())
				if err != nil {
					val = 3000
				}
				routing.ResolverPort = int(val)
			default:
				if c.Val() != "}" {
					return c.Errf("unknown property '%s'", c.Val())
				}
			}
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		routing.Next = next
		return routing
	})

	return nil
}

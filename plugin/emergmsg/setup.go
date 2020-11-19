package emergmsg

import (
	"errors"
	"regexp"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("emergmsg", setup) }

func validDelimeter(delim string) bool {
	validDelimRegex := regexp.MustCompile("^[a-z][a-z0-9]+")
	return validDelimRegex.MatchString(delim)
}

type setupParameters struct {
	Delimeter string
	RedisAddr string
	RedisKey  string
}

func setup(c *caddy.Controller) error {
	p := &setupParameters{}

	for c.Next() {
		// Delimeter
		if !c.NextArg() {
			return plugin.Error("emergmsg", c.ArgErr())
		}

		if !validDelimeter(c.Val()) {
			return plugin.Error("emergmsg", errors.New("delimiter is not valid"))
		}
		p.Delimeter = c.Val()

		// Redis address
		if !c.NextArg() {
			return plugin.Error("emergmsg", c.ArgErr())
		}
		p.RedisAddr = c.Val()

		// Redis key to write to
		if !c.NextArg() {
			return plugin.Error("emergmsg", c.ArgErr())
		}
		p.RedisKey = c.Val()
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		em, err := New(next, p.Delimeter, p.RedisAddr, p.RedisKey)
		if err != nil {
			plugin.Error("emergmsg", err)
		}
		return em
	})

	return nil
}

package emergmsg

import (
	"context"
	"fmt"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/go-redis/redis/v8"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("emergmsg")

type Emergmsg struct {
	Next  plugin.Handler
	delim string
	rdb   *redis.Client
	key   string
}

// Pull out the message
func (e *Emergmsg) parseMsg(name string) string {
	idx := strings.LastIndex(name, "."+e.delim+".")
	if idx < 0 {
		return ""
	}
	msg := name[:idx]
	return msg
}

func (e *Emergmsg) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	a := &dns.Msg{}
	a.SetReply(r)
	a.Authoritative = true
	defer w.WriteMsg(a)

	if len(r.Question) == 0 {
		return dns.RcodeNameError, nil
	}
	q := r.Question[0]
	name := q.Name
	log.Debugf("Received emergency request: %+v", name)
	msg := e.parseMsg(name)
	if msg == "" {
		return dns.RcodeNameError, nil
	}
	err := e.rdb.RPush(ctx, e.key, msg).Err()
	if err != nil {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	return dns.RcodeNameError, nil
}

func (e *Emergmsg) Name() string { return "emergmsg" }

func New(next plugin.Handler, delim, addr, key string) (*Emergmsg, error) {
	if delim == "" {
		return nil, fmt.Errorf("delim cannot empty")
	}
	if addr == "" {
		return nil, fmt.Errorf("addr cannot empty")
	}
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	em := &Emergmsg{
		Next:  next,
		delim: delim,
		rdb:   rdb,
		key:   key,
	}
	return em, nil
}

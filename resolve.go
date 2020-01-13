package main

import (
	"net"
	"regexp"

	"github.com/miekg/dns"
)

func resolveName(name string, ret chan<- net.IP) {
	cl := dns.Client{}
	a := dns.Msg{}
	var (
		ok  bool
		ans *dns.A
	)
	t := dns.TypeA
	ip := regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`)
	for ans == nil {
		if ip.Match([]byte(name)) {
			logger.Info(name + " is ip address, not domain")
			ret <- nil
			return
		}
		a.SetQuestion(name, t)
		res, _, err := cl.Exchange(&a, ds+":53")
		if err != nil || len(res.Answer) == 0 {
			ret <- nil
			return
		}
		ans, ok = res.Answer[0].(*dns.A)
		if !ok {
			cn, cok := res.Answer[0].(*dns.CNAME)
			if !cok {
				ret <- nil
				return
			}
			name = cn.Target
			t = dns.TypeCNAME
		}
	}
	ret <- ans.A
}

package main

import (
	"net"

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
	d := ds
	n := name
	for ans == nil {
		if ip.Match([]byte(n)) {
			logger.Info(name + " is ip address, not domain")
			ret <- nil
			return
		}
		if n != name {
			d = "8.8.8.8"
		}
		a.SetQuestion(n, t)
		res, _, err := cl.Exchange(&a, d+":53")
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
			n = cn.Target
		}
	}
	ret <- ans.A
}

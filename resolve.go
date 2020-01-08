package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
)

func resolveName(name string, ret chan<- net.IP) {
	d := os.Getenv("DNS_SERVER")
	log.Println(d)
	cl := dns.Client{}
	a := dns.Msg{}
	var (
		ok  bool
		ans *dns.A
	)
	t := dns.TypeA
	for ans == nil {
		a.SetQuestion(name, t)
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
			name = cn.Target
			t = dns.TypeCNAME
		}
	}
	ret <- ans.A
}

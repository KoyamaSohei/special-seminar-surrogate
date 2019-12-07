package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
)

func serveEcho() {
	i := os.Getenv("TARGET_IP")

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		if r.Question[0].Qtype != dns.TypeA {
			return
		}
		m.SetReply(r)
		rr := &dns.A{
			Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP(i),
		}
		m.Answer = append(m.Answer, rr)
		w.WriteMsg(m)
	})

	go func() {
		server := &dns.Server{Addr: "[::]:53", Net: "udp", TsigSecret: nil, ReusePort: false}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

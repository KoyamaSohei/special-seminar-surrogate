package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/miekg/dns"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Signal (%s) received, stopping\n", s)
}

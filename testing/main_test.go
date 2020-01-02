package main

import (
	"log"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	cl := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion("example.com.", dns.TypeA)
	res, _, err := cl.Exchange(&m, "localhost:53")
	if err != nil {
		log.Fatal(err)
	}
	a := res.Answer[0].(*dns.A)
	assert.Equal(t, "12.34.56.78", a.A.String())
}

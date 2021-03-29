// Package vpnrouting implements a plugin that returns details about the resolving
// querying it.
package vpnrouting

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

const name = "vpnrouting"

// Vpnrouting is a plugin that returns your IP address, port and the protocol used for connecting
// to CoreDNS.
type Vpnrouting struct {
	Next             plugin.Handler
	ResolverIP       string
	ResolverPort     int
	ResolverGrpcPort int
}

// New default values
func New() *Vpnrouting {
	return &Vpnrouting{
		ResolverIP:       "127.0.0.1",
		ResolverPort:     3000,
		ResolverGrpcPort: 8888,
	}
}

// Name implements the Handler interface.
func (rout Vpnrouting) Name() string { return name }

// ResolverHostname d
func (rout Vpnrouting) ResolverHostname() string {
	return "http://" + rout.ResolverIP + ":" + fmt.Sprint(rout.ResolverGrpcPort)
}

func (rout Vpnrouting) ResolveGrpcAddress() string {
	return rout.ResolverIP + ":" + fmt.Sprint(rout.ResolverPort)
}

// ServeDNS implements the plugin.Handler interface.
func (rout Vpnrouting) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	ip := state.IP()
	var rr dns.RR

	geo := geoLookup(ip)
	resolver, err := rout.resolve(state.Name(), geo)

	if err != nil {
		return dns.RcodeServerFailure, err
	}

	if resolver.IP == "" {
		return plugin.NextOrFailure(rout.Name(), rout.Next, ctx, w, r)
	}

	switch state.Family() {
	case 1:
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.(*dns.A).A = net.ParseIP(resolver.IP).To4()
	case 2:
		rr = new(dns.AAAA)
		rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
		rr.(*dns.AAAA).AAAA = net.ParseIP(resolver.IP)
	}
	srv := new(dns.SRV)
	srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
	if state.QName() == "." {
		srv.Hdr.Name = "_" + state.Proto() + state.QName()
	}
	port, _ := strconv.Atoi(state.Port())
	srv.Port = uint16(port)
	srv.Target = "."

	a.Extra = []dns.RR{srv}
	a.Answer = []dns.RR{rr}

	w.WriteMsg(a)

	return dns.RcodeSuccess, nil
}

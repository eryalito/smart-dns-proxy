package dns

import (
	"log"
	"net"

	"github.com/eryalito/smart-dns-proxy/internal/data"
	"github.com/miekg/dns"
)

type Server struct {
	Addr    string
	Net     string
	Querier *data.Querier
}

func (s *Server) Start() error {
	dns.HandleFunc(".", s.dnsHandler)
	server := &dns.Server{Addr: s.Addr, Net: s.Net}
	return server.ListenAndServe()
}

func (s *Server) dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	// Forward query to 8.8.8.8
	c := new(dns.Client)
	resp, _, err := c.Exchange(r, "8.8.8.8:53")
	if err != nil {
		log.Printf("DNS forward error: %v", err)
		dns.HandleFailed(w, r)
		return
	}
	data := s.Querier.GetData()
	// Check if the answered at least once of the Answered IPs is blocked
	present := false
	blocked := false
	provider := ""
	for _, answer := range resp.Answer {
		if answer.Header().Rrtype == dns.TypeA {
			ip := answer.(*dns.A).A.String()
			for _, element := range data.Elements {
				present = true
				provider = element.Provider
				if element.IP == ip {
					if element.Blocked {
						blocked = true
						break
					}
				}
			}
			if blocked {
				break
			}
		}
	}

	// If the IP is not present on the list, just forward the response
	if !present {
		w.WriteMsg(resp)
		return
	}

	// If none IP is blocked, just forward the response
	if !blocked {
		w.WriteMsg(resp)
		return
	}

	// If it's blocked, but no provider found there's nothing to do, forward the request
	if provider == "" {
		log.Println("No provider found for blocked IP")
		w.WriteMsg(resp)
		return
	}

	// If present, blocked and provider detected, try to find a clean IP and send that

	for _, element := range data.Elements {
		if element.Provider == provider && !element.Blocked {
			resp.Answer = []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:   r.Question[0].Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    300,
					},
					A: net.ParseIP(element.IP),
				},
			}
			break
		}
	}
	w.WriteMsg(resp)
}

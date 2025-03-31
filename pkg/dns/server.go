package dns

import (
	"baderkha-no-dns/pkg/dns/blocklist"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

var server = sync.OnceValue(func() *dns.Server {
	dns.HandleFunc(".", handleRequest)
	server := &dns.Server{Addr: ":53", Net: "udp"}
	return server
})

var AdsBlocked int64

// start the dns server
func Start() {
	// Start DNS server
	err := PreStartSteps()
	if err != nil {
		log.Fatal(err)
	}
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println(" got interrupt signal: ", <-sigChan)
		fmt.Println("running stop steps")
		Stop()
		if err != nil {
			fmt.Println(err.Error())
		}
		os.Exit(0)
	}()
	err = server().ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start DNS server: %v\n", err)
	}

}

func Stop() error {
	server().Shutdown()
	return PreStopSteps()
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Prepare the response
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	for _, question := range r.Question {
		tn := time.Now()
		if blocklist.Checker().Has(question.Name) {
			atomic.AddInt64(&AdsBlocked, 1)
			fmt.Println("This dns is blocked[", question.Name, "]")
			// If it's an ad domain, return a fake IP (e.g., 0.0.0.0)
			rr := &dns.A{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP("0.0.0.0"),
			}
			m.Answer = append(m.Answer, rr)
		} else {

			// Forward the non-blocked request to Google DNS (8.8.8.8)
			response, err := forwardQueryToDNSServer(question.Name)
			if err != nil {
				log.Printf("Failed to forward DNS request for %s: %v", question.Name, err)
				continue
			}

			fmt.Println("Forwarding DNS[", question.Name, "]", "is [", response[0].String(), "]", time.Since(tn))

			m.Answer = append(m.Answer, response...)
		}
	}

	// Send the response back to the client
	w.WriteMsg(m)
}
func forwardQueryToDNSServer(domain string) ([]dns.RR, error) {
	client := new(dns.Client)
	client.Timeout = 5 * time.Second // Set a timeout for the DNS query

	// Prepare the DNS message with the domain query
	message := new(dns.Msg)
	message.SetQuestion(domain, dns.TypeA)

	// List of DNS servers to try (5 different servers for redundancy)
	servers := []string{
		"1.1.1.1:53",         // Cloudflare
		"8.8.8.8:53",         // Google DNS
		"8.8.4.4:53",         // Google DNS (Secondary)
		"9.9.9.9:53",         // Quad9 DNS
		"149.112.112.112:53", // DNS.WATCH
	}

	var answers []dns.RR

	// Try multiple DNS servers
	for _, server := range servers {
		resp, _, err := client.Exchange(message, server)
		if err != nil {
			// Log error and continue to the next server
			fmt.Printf("Error querying DNS server %s for %s: %v\n", server, domain, err)
			continue
		}

		// If we get a response, parse and return it
		for _, answer := range resp.Answer {
			if aRecord, ok := answer.(*dns.A); ok {
				answers = append(answers, aRecord)
			}
		}

		// If we got any answers, return them
		if len(answers) > 0 {
			return answers, nil
		}
	}

	// If no servers responded successfully, return an error
	return nil, fmt.Errorf("unable to forward query for %s after trying all DNS servers", domain)
}

package dns

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"
)

var adDomains = map[string]struct{}{
	"ads.google.com.": {},
	// Add more known ad domains here
}

var server = sync.OnceValue(func() *dns.Server {
	dns.HandleFunc(".", handleRequest)
	server := &dns.Server{Addr: ":53", Net: "udp"}
	return server
})

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
		_, blocked := adDomains[strings.ToLower(question.Name)]
		fmt.Println(question, blocked)

		if blocked {
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

			m.Answer = append(m.Answer, response...)
		}
	}

	// Send the response back to the client
	w.WriteMsg(m)
}

func forwardQueryToDNSServer(domain string) ([]dns.RR, error) {
	client := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion(domain, dns.TypeA)

	// Forward the query to Google's DNS server (8.8.8.8)
	server := "8.8.8.8:53"
	resp, _, err := client.Exchange(message, server)
	if err != nil {
		return nil, fmt.Errorf("error forwarding query to %s: %v", server, err)
	}

	// If there are A records, return them
	var answers []dns.RR
	for _, answer := range resp.Answer {
		if aRecord, ok := answer.(*dns.A); ok {
			answers = append(answers, aRecord)
		}
	}

	return answers, nil
}

// The code is adapted by @JosueEncinar from https://github.com/blackhat-go/bhg/blob/master/ch-5/subdomain_guesser/main.go
// For @six2dez's reconfwt
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/miekg/dns"
)

type empty struct{}

type result struct {
	IPAddress string
	Hostname  string
}

func lookupA(fqdn, resolver string) string {
	var m dns.Msg
	var in *dns.Msg
	var err error
	count := 0
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	for {
		in, err = dns.Exchange(&m, resolver)
		if err != nil || len(in.Answer) < 1 {
			count += 1
			if count >= 5 {
				return ""
			}
		} else {
			break
		}
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			return a.A.String() // Only one ip is desired
		}
	}
	return ""
}

func lookupCNAME(fqdn, resolver string) []string {
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, resolver)
	if err != nil || len(in.Answer) < 1 {
		return nil
	}
	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, c.Target)
		}
	}
	return fqdns
}

func lookup(fqdn, resolver string) result {
	var cfqdn = fqdn
	for {
		cnames := lookupCNAME(cfqdn, resolver)
		if cnames == nil {
			break
		} else if len(cnames) > 0 {
			cfqdn = cnames[0]
		}
	}
	ip := lookupA(cfqdn, resolver)
	if ip != "" {
		return result{IPAddress: ip, Hostname: fqdn}
	}
	return result{}
}

func worker(tracker chan empty, fqdns chan string, gather chan result, resolver string) {
	for fqdn := range fqdns {
		result := lookup(fqdn, resolver)
		if result.Hostname != "" {
			gather <- result
		}
	}
	var e empty
	tracker <- e
}

func main() {
	var (
		flDomains  = flag.String("d", "", "List of domains to resolve.")
		flThreads  = flag.Int("t", 100, "The amount of threads to use.")
		flResolver = flag.String("r", "8.8.8.8:53", "Resolver DNS to use.")
	)
	flag.Parse()

	if *flDomains == "" {
		fmt.Println("-d are required")
		os.Exit(1)
	}

	var results []result

	fqdns := make(chan string, *flThreads)
	gather := make(chan result)
	tracker := make(chan empty)

	fh, err := os.Open(*flDomains)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)

	for i := 0; i < *flThreads; i++ {
		go worker(tracker, fqdns, gather, *flResolver)
	}

	go func() {
		for r := range gather {
			results = append(results, r)
		}
		var e empty
		tracker <- e
	}()

	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s", scanner.Text())
	}

	close(fqdns)
	for i := 0; i < *flThreads; i++ {
		<-tracker
	}
	close(gather)
	<-tracker

	for _, r := range results {
		fmt.Println(r.Hostname + " " + r.IPAddress)
	}
}

package main

import (
	"baderkha-no-dns/pkg/dns/blocklist"
	"fmt"
	"log"
	"time"
)

func main() {
	slc, err := blocklist.LoadStorage()
	if err != nil {
		log.Fatal(err)
	}
	q := "www.appnerve.com"
	t := time.Now()
	ok := blocklist.Checker().Has(q)
	fmt.Println(time.Since(t), len(slc), ok)
}

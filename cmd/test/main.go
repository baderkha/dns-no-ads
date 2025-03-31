package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	t := time.Now()
	kv := make(map[string]struct{})
	dat, err := ioutil.ReadFile(filepath.Join("resources", "dns", "blocklist", "block-list.txt"))
	if err != nil {
		log.Fatal(err)
	}
	all := strings.Split(string(dat), "\n")
	_ = all
	for _, v := range all {
		kv[v] = struct{}{}
	}
	fmt.Println(time.Since(t), len(all), len(kv))
}

package blocklist

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

const ExpectedFileRows = 400000

var Checker = sync.OnceValue(func() Store {
	return NewBinarytreeStore()
})

type Store interface {
	Has(domain string) bool
}

type BinarySearchStore struct {
	dat []string
}

func (b *BinarySearchStore) Has(domain string) bool {
	domainClean := strings.TrimRight(domain, ".")
	_, ok := slices.BinarySearch(b.dat, domainClean)
	return ok
}

func NewBinarytreeStore() *BinarySearchStore {
	slc, err := LoadStorage()
	if err != nil {
		log.Fatal(err)
	}
	return &BinarySearchStore{
		dat: slc,
	}
}

func LoadStorage() ([]string, error) {
	slc := make([]string, 0, ExpectedFileRows)
	f, err := os.Open(filepath.Join("resources", "dns", "blocklist", "block-list.txt"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024*10) // set max token size to 1MB
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		slc = append(slc, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return slc, nil
}

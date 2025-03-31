package blocklist

import (
	"baderkha-no-dns/pkg/fs"
	"os"
	"path/filepath"
)

const ExpectedFileRows = 300000

type Store interface {
	Has(domain string) bool
	HasM(domain ...string) []bool
}

type InMemSimpleKV struct {
	simpleKv map[string]struct{}
}

func NewInMemSimpleKV() *InMemSimpleKV {
	return &InMemSimpleKV{
		simpleKv: make(map[string]struct{}, ExpectedFileRows),
	}
}

func LoadStorage(kv map[string]struct{}) error {

	f, err := os.Open(filepath.Join("resources", "dns", "blocklist", "block-list.txt"))
	if err != nil {
		return err
	}
	defer f.Close()
	err = fs.LineByLineLg(f, func(line string) error {
		kv[line] = struct{}{}
		return nil
	})
	return err
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func DownloadAndFilterDomains(url, folder, filename string) error {
	// Create folder if it doesn't exist
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("creating folder: %w", err)
	}

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", resp.Status)
	}

	// Process and filter domains
	domains := []string{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Expect lines like: 0.0.0.0 domain.tld
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "0.0.0.0" {
			domains = append(domains, fields[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning response: %w", err)
	}

	// Sort domains alphabetically
	sort.Strings(domains)

	// Write output
	outPath := filepath.Join(folder, filename)
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer outFile.Close()

	for _, domain := range domains {
		_, _ = outFile.WriteString(domain + "\n")
	}

	return nil
}

func main() {
	expectedFileOut := filepath.Join("resources", "dns", "blocklist", "block-list.txt")
	os.Remove(expectedFileOut)
	dbPath := filepath.Join("resources", "db", "block-list.db")

	err := DownloadAndFilterDomains("https://raw.githubusercontent.com/hagezi/dns-blocklists/refs/heads/main/hosts/multi.txt", filepath.Join("resources", "dns", "blocklist"), "block-list.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done! Data saved to:", dbPath)
}

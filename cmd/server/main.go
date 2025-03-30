package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	// Set DNS via netsh
	dns1 := "127.0.0.1"
	interfaceName := "Ethernet 4"
	// The netsh command as a single string to be passed to PowerShell
	psCommand := fmt.Sprintf(`Start-Process netsh -ArgumentList 'interface ipv4 set dnsserver "%s" static %s' -Verb runAs`, interfaceName, dns1)
	fmt.Println(psCommand)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to run elevated netsh command:", err)
		return
	}

	fmt.Println("DNS server added successfully (if you accepted UAC).")
}
func getActiveConnection() (string, error) {
	// Command to list network interfaces
	cmd := exec.Command("netsh", "interface", "show", "interface")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run netsh command: %v", err)
	}

	// Search for the active connection
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Look for the line that mentions "Connected" status
		if strings.Contains(line, "Connected") {
			// The interface name should be in the first column
			fmt.Println(line)
			columns := strings.Split(line, " ")

			if len(columns) > 0 {
				fmt.Println(columns)
				return columns[4], nil // The interface name is the first column
			}
		}
	}

	return "", fmt.Errorf("no active connection found")
}

package dns

import (
	"fmt"
	"net"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var _ dnsSetter = &WinDnsSet{}

type WinDnsSet struct {
}

// GetDefaultDevice implements DnsSetter.
func (w *WinDnsSet) GetDefaultDevice() (*NetworkDevice, error) {
	// Step 1: use GetBestInterfaceEx to get default interface index
	ip := net.ParseIP("8.8.8.8").To4()
	if ip == nil {
		return nil, fmt.Errorf("invalid IP")
	}

	sa := &windows.SockaddrInet4{Addr: [4]byte{ip[0], ip[1], ip[2], ip[3]}}
	var ifaceIndex uint32
	if err := windows.GetBestInterfaceEx(sa, &ifaceIndex); err != nil {
		return nil, fmt.Errorf("GetBestInterfaceEx failed: %w", err)
	}

	// Step 2: get the MAC address of that interface
	var targetMAC string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if iface.Index == int(ifaceIndex) {
			targetMAC = iface.HardwareAddr.String()
			break
		}
	}

	if targetMAC == "" {
		return nil, fmt.Errorf("no MAC found for interface index %d", ifaceIndex)
	}

	// Step 3: get adapter GUID from GetAdaptersAddresses
	var size uint32
	windows.GetAdaptersAddresses(windows.AF_UNSPEC, 0, 0, nil, &size)
	buf := make([]byte, size)
	addr := (*windows.IpAdapterAddresses)(unsafe.Pointer(&buf[0]))
	err = windows.GetAdaptersAddresses(windows.AF_UNSPEC, 0, 0, addr, &size)
	if err != nil {
		return nil, fmt.Errorf("GetAdaptersAddresses failed: %w", err)
	}

	for a := addr; a != nil; a = a.Next {
		mac := net.HardwareAddr(a.PhysicalAddress[:a.PhysicalAddressLength]).String()
		if mac == targetMAC {
			guid := windows.BytePtrToString(a.AdapterName)
			name := windows.UTF16PtrToString(a.FriendlyName)
			return &NetworkDevice{
				Description:      name,
				DeviceIdentifier: guid, // This is the adapter's GUID used in the registry
				MacAddress:       mac,
			}, nil
		}
	}

	return nil, fmt.Errorf("default adapter not found")
}

func (w *WinDnsSet) ListDevices() ([]*NetworkDevice, error) {
	var size uint32
	windows.GetAdaptersAddresses(windows.AF_UNSPEC, 0, 0, nil, &size)

	buf := make([]byte, size)
	addr := (*windows.IpAdapterAddresses)(unsafe.Pointer(&buf[0]))

	err := windows.GetAdaptersAddresses(windows.AF_UNSPEC, 0, 0, addr, &size)
	if err != nil {
		return nil, fmt.Errorf("GetAdaptersAddresses failed: %w", err)
	}

	var result []*NetworkDevice

	for a := addr; a != nil; a = a.Next {
		mac := net.HardwareAddr(a.PhysicalAddress[:a.PhysicalAddressLength]).String()
		guid := windows.BytePtrToString(a.AdapterName)
		name := windows.UTF16PtrToString(a.FriendlyName)

		result = append(result, &NetworkDevice{
			Description:      name,
			DeviceIdentifier: guid, // This is the adapter's GUID used in the registry
			MacAddress:       mac,
		})
	}

	return result, nil
}

// SetDnsForDevice implements DnsSetter.
func (w *WinDnsSet) SetDnsForDevice(n *NetworkDevice, dns []string) error {
	adapterGUID := n.DeviceIdentifier
	// Build the registry path
	path := fmt.Sprintf(`SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\%s`, adapterGUID)

	// Open the key
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Join DNS servers with newline for REG_MULTI_SZ format
	dnsValue := strings.Join(dns, ",")

	// Set the NameServer value
	err = key.SetStringValue("NameServer", dnsValue)
	if err != nil {
		return fmt.Errorf("failed to set DNS: %w", err)
	}
	return nil
}

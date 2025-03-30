package dns

import (
	"runtime"
	"sync"
)

var Setter = sync.OnceValue(func() dnsSetter {
	switch {
	case runtime.GOOS == "windows":
		return &WinDnsSet{}
	}
	return nil

})

type NetworkDevice struct {
	DeviceIdentifier string
	Description      string
	Meta             any
	MacAddress       string
}

type dnsSetter interface {
	SetDnsForDevice(n *NetworkDevice, dns []string) error
	GetDefaultDevice() (*NetworkDevice, error)
	ListDevices() ([]*NetworkDevice, error)
}

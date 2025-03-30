package dns

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

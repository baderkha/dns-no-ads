package dns

import "fmt"

func PreStartSteps() error {

	localDns := []string{"127.0.0.1"}

	dnsDeviceId, err := Setter().GetDefaultDevice()
	if err != nil {
		fmt.Println("could not get default device", err)
	}

	return Setter().SetDnsForDevice(dnsDeviceId, localDns)

}
func PreStopSteps() error {
	fallbackDns := []string{"8.8.8.8", "8.8.4.4"}
	dnsDeviceId, err := Setter().GetDefaultDevice()
	if err != nil {
		fmt.Println("could not get default device", err)
	}
	return Setter().SetDnsForDevice(dnsDeviceId, fallbackDns)
}

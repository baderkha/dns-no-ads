package dns

import "fmt"

func PreStartSteps() {

	localDns := []string{"127.0.0.1"}

	dnsDeviceId, err := Setter().GetDefaultDevice()
	if err != nil {
		fmt.Println("could not get default device", err)
	}

	Setter().SetDnsForDevice(dnsDeviceId, localDns)

}
func PreStopSteps() {
	fallbackDns := []string{"8.8.8.8", "8.8.4.4"}
	dnsDeviceId, err := Setter().GetDefaultDevice()
	if err != nil {
		fmt.Println("could not get default device", err)
	}
	Setter().SetDnsForDevice(dnsDeviceId, fallbackDns)
}

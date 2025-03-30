package dns

import (
	"context"
	"fmt"

	"github.com/taigrr/systemctl"
)

func PreStartSteps() error {

	err := systemctl.Stop(context.Background(), "systemd-resolved", systemctl.Options{
		UserMode: false,
	})
	fmt.Println("ok")
	return err
}

func PreStopSteps() error {
	err := systemctl.Start(context.Background(), "systemd-resolved", systemctl.Options{
		UserMode: false,
	})
	fmt.Println(err)
	return err
}

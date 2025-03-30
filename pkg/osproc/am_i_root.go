package osproc

import (
	"log"
	"os"
	"os/user"
	"runtime"
	"sync"
)

type IsRootFunc func() bool

func IsRootWindows() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func IsRootLinux() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

var isRootInit = sync.OnceValue(func() IsRootFunc {
	switch runtime.GOOS {
	case "windows":
		return IsRootWindows
	case "linux":
		return IsRootLinux
	}

	return nil
})

var IsRoot = isRootInit()

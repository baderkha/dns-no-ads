package osproc

import (
	"os"
	"runtime"
	"sync"
)

type IsRootFunc func() bool

func IsRootWindows() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

var isRootInit = sync.OnceValue(func() IsRootFunc {
	switch runtime.GOOS {
	case "windows":
		return IsRootWindows
	}
	return nil
})

var IsRoot = isRootInit()

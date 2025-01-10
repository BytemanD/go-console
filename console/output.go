package console

import (
	"fmt"
	"sync"
)

var outputLock *sync.Mutex
var enableLog bool

func withOutputLock(outputFunc func()) {
	outputLock.Lock()
	defer outputLock.Unlock()

	fmt.Print("\033[2K\r")
	outputFunc()
}

func Printf(format string, v ...any) {
	withOutputLock(
		func() {
			fmt.Printf(format, v...)
			pbrManager.Output()
		})
}

func Println(a ...any) {
	withOutputLock(func() {
		fmt.Println(a...)
		pbrManager.Output()
	})
}

func Print(a ...any) {
	withOutputLock(func() {
		fmt.Print(a...)
		pbrManager.Output()
	})
}

func PkgEnableLog() {
	enableLog = true
}
func PkgDisablePkgLog() {
	enableLog = false
}
func init() {
	outputLock = &sync.Mutex{}
	enableLog = false
}

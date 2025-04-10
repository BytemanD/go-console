package console

import (
	"fmt"

	"github.com/samber/lo"
)

var outputSync = lo.Synchronize()

var enableLog bool

func withOutputLock(outputFunc func()) {
	outputSync.Do(func() {
		fmt.Print("\033[2K\r")
		outputFunc()
		defaultPbrGroup.Output()
	})
}

func Printf(format string, v ...any) {
	withOutputLock(
		func() {
			fmt.Printf(format, v...)
		})
}

func Println(a ...any) {
	withOutputLock(func() {
		fmt.Println(a...)
		defaultPbrGroup.Output()
	})
}

func Print(a ...any) {
	withOutputLock(func() {
		fmt.Print(a...)
		defaultPbrGroup.Output()
	})
}

func PkgEnableLog() {
	enableLog = true
}
func PkgDisablePkgLog() {
	enableLog = false
}
func init() {
	// outputLock = &sync.Mutex{}
	enableLog = false
}

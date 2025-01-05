package console

import (
	"fmt"
	"sync"
)

var outputLock *sync.Mutex

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

func init() {
	outputLock = &sync.Mutex{}
}

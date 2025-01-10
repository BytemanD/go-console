package console

import (
	"testing"
	"time"
)

func TestPbrs(t *testing.T) {
	PkgEnableLog()
	EnableLogDebug()
	SetLogFile("/tmp/go_console.log")
	Debug("start tasks")
	go func() {
		for i := 0; i < 10; i++ {
			Println("pring task ", i)
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		for i := 0; i < 60; i++ {
			Debug("log %d", i)
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		for i := 0; i < 60; i++ {
			DebugS("log with struct", "i", i)
			time.Sleep(time.Millisecond * 20)
		}
	}()

	pbr := NewPbr(100, "foo")
	go func() {
		for i := 0; i < 100; i++ {
			pbr.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()
	pbr2 := NewPbr(100, "bar")
	go func() {
		for i := 0; i < 100; i++ {
			pbr2.Increment()
			time.Sleep(time.Millisecond * 20)
		}
	}()

	pbr3 := NewPbr(100, "bar3 无进度")
	pbr3.Done()
	WaitAllPbrs()
	PkgDisablePkgLog()
	Println("====== start new tasks ========")

	pbr4 := NewPbr(100, "bar4")
	go func() {
		for i := 0; i < 100; i++ {
			pbr4.Increment()
			time.Sleep(time.Millisecond * 20)
		}
	}()
	pbr5 := NewPbr(100, "bar5")
	go func() {
		for i := 0; i < 100; i++ {
			pbr5.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()
	WaitAllPbrs()

	doneNum := GetPbrNum()
	if doneNum != 0 {
		t.Errorf("expected done num be 0, but got %d", doneNum)
	}
}

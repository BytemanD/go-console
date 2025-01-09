package console

import (
	"testing"
	"time"
)

func TestPbrs(t *testing.T) {
	EnableLogDebug()
	SetLogFile("/tmp/go_console.log")
	go func() {
		for i := 0; i < 10; i++ {
			Println("pring task ", i)
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		for i := 0; i < 60; i++ {
			Info("log %d", i)
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		for i := 0; i < 60; i++ {
			InfoS("log with struct", "i", i)
			time.Sleep(time.Millisecond * 20)
		}
	}()

	go func() {
		pbr := NewPbr(100, "foo")
		for i := 0; i < 100; i++ {
			pbr.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		pbr2 := NewPbr(100, "bar")
		for i := 0; i < 100; i++ {
			pbr2.Increment()
			time.Sleep(time.Millisecond * 20)
		}
	}()

	pbr3 := NewPbr(100, "bar3 无进度")
	pbr3.Done()
	WaitPbrs()
	doneNum := GetPbrNum()
	if doneNum != 0 {
		t.Errorf("expected done num be 0, but got %d", doneNum)
	}
}

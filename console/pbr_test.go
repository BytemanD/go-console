package console

import (
	"fmt"
	"math/rand/v2"
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

	pbr := NewPbr(100, "默认的进度条1")
	go func() {
		for i := 0; i < 100; i++ {
			pbr.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()
	pbr2 := NewPbr(100, "默认的进度条2")
	go func() {
		for i := 0; i < 100; i++ {
			pbr2.Increment()
			time.Sleep(time.Millisecond * 20)
		}
	}()

	pbr3 := NewPbr(100, "默认的进度条3 无进度")
	pbr3.Done()

	WaitAllPbrs()

	PkgDisablePkgLog()
	Println("====== start new tasks ========")

	pbr4 := NewPbr(100, "默认进度条4")
	go func() {
		for i := 0; i < 100; i++ {
			pbr4.Increment()
			time.Sleep(time.Millisecond * 20)
		}
	}()
	pbr5 := NewPbr(100, "默认进度条5")
	go func() {
		for i := 0; i < 100; i++ {
			pbr5.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()
	WaitAllPbrs()
	Println("====== start new tasks ========")
	doneNum := GetPbrNum()
	if doneNum != 0 {
		t.Errorf("expected done num be 0, but got %d", doneNum)
	}

	pbr6 := NewPbrWithTheme(100, "简单的进度条", THEME_SIMPLE)
	go func() {
		for i := 0; i < 100; i++ {
			pbr6.Increment()
			time.Sleep(time.Millisecond * 10)
		}
	}()
	pbr7 := NewPbrWithTheme(100, "自定义进度条", CustomeTheme("*"))
	go func() {
		for i := 0; i < 100; i++ {
			pbr7.Increment()
			time.Sleep(time.Millisecond * 20)
		}
	}()
	WaitAllPbrs()
}
func TestPbrsCreatedSync(t *testing.T) {
	themeChars := []string{"#", "*", ">", "o", "@", "~"}
	createNewPbr := func(i int) {
		pbr := NewPbrWithTheme(100, fmt.Sprintf("并发创建的进度条%d", i), CustomeTheme(themeChars[i]))
		for range 100 {
			pbr.Increment()
			time.Sleep(time.Millisecond * time.Duration(rand.IntN(20)))
		}
	}
	for i := range 5 {
		go createNewPbr(i + 1)
	}
	time.Sleep(time.Second * 2)
	WaitAllPbrs()
}

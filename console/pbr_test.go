package console

import (
	"fmt"
	"testing"
	"time"
)

func doSomething(pbr *Pbr, interval int) {
	for range pbr.Total {
		pbr.Increment()
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}

func TestPbrs(t *testing.T) {
	PkgEnableLog()
	EnableLogDebug()
	SetLogFile("/tmp/go_console.log")
	Debug("start tasks")
	go func() {
		for i := range 10 {
			Println("pring task ", i)
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		for i := range 60 {
			Debug("log %d", i)
			time.Sleep(time.Millisecond * 10)
		}
	}()
	go func() {
		for i := range 60 {
			DebugS("log with struct", "i", i)
			time.Sleep(time.Millisecond * 20)
		}
	}()

	pbr := NewPbr(100, "默认的进度条1")
	pbr2 := NewPbr(100, "默认的进度条2")
	go doSomething(pbr, 10)
	go doSomething(pbr2, 20)

	pbr3 := NewPbr(100, "默认的进度条3 无进度")
	pbr3.Done()

	WaitAllProgressBar()

	PkgDisablePkgLog()
	Println("====== start new tasks ========")

	pbr4 := NewPbr(100, "默认进度条4")
	pbr5 := NewPbr(100, "默认进度条5")
	go doSomething(pbr4, 20)
	go doSomething(pbr5, 10)

	WaitAllProgressBar()
	Println("====== start new tasks ========")
	doneNum := GetPbrNum()
	if doneNum != 0 {
		t.Errorf("expected done num be 0, but got %d", doneNum)
	}

	pbr6 := NewPbrWithTheme(100, "简单的进度条", THEME_SIMPLE)
	go doSomething(pbr6, 10)
	pbr7 := NewPbrWithTheme(100, "自定义进度条", CustomeTheme("*"))
	go doSomething(pbr7, 20)

	WaitAllProgressBar()
}

func TestPbrsCreatedSync(t *testing.T) {
	themeChars := []string{"#", "*", ">", "o", "@", "~"}
	for i := range 5 {
		pbr := NewPbrWithTheme(100, fmt.Sprintf("并发创建的进度条%d", i), CustomeTheme(themeChars[i]))
		go doSomething(pbr, i*10)
	}
	time.Sleep(time.Second * 2)
	WaitAllProgressBar()
}

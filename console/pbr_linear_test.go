package console

import (
	"fmt"
	"testing"
	"time"
)

func doSomething(pbr *ProgressLinear, interval int) {
	for range pbr.Total {
		pbr.Increment()
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
	SuccessS("complted", "title", pbr.Title)
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

	pbr := NewProgressLinear(100, "默认的进度条1")
	pbr2 := NewProgressLinear(100, "默认的进度条2")
	go doSomething(pbr, 10)
	go doSomething(pbr2, 20)

	pbr3 := NewProgressLinear(100, "默认的进度条3(无进度)")
	pbr3.ForceDone()

	WaitAllProgressBar()

	PkgDisablePkgLog()
	Println("====== start new tasks ========")

	pbr4 := NewProgressLinear(100, "默认进度条4")
	pbr5 := NewProgressLinear(100, "默认进度条5")
	go doSomething(pbr4, 20)
	go doSomething(pbr5, 10)

	WaitAllProgressBar()
	Println("====== start new tasks ========")
	doneNum := ProgressCount()
	if doneNum != 0 {
		t.Errorf("expected done num be 0, but got %d", doneNum)
	}

	pbr6 := NewProgressLinear(100, "简单的进度条", THEME_CONCISE)
	go doSomething(pbr6, 10)
	pbr7 := NewProgressLinear(100, "自定义进度条", *NewProgressBarTheme("*", ""))
	go doSomething(pbr7, 20)

	WaitAllProgressBar()
}

func TestPbrsSync(t *testing.T) {
	EnableLogDebug()
	themes := []ProgressBarTheme{THEME_LIGHT, THEME_CONCISE, THEME_EMOJI_AIRPLANE, THEME_EMOJI_HEART}
	for i := range len(themes) {
		pbr := NewProgressLinear(100, fmt.Sprintf("并发创建的进度条%d", i), themes[i])
		go doSomething(pbr, (i+1)*10)
	}
	WaitAllProgressBar()
}

func TestPbrsSyncWithCustomTemplate(t *testing.T) {
	themeChars := []string{"#", "*", ">", "o", "@", "~"}
	for i := range len(themeChars) {
		theme := NewProgressBarTheme(themeChars[i], "{{.Title}} [{{.Percent}}%]: {{.Progress}}")
		pbr := NewProgressLinear(100, fmt.Sprintf("并发创建的进度条%d", i), *theme)
		go doSomething(pbr, i*10)
	}
	WaitAllProgressBar()
}

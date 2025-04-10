package console

import (
	"fmt"
	"testing"
	"time"
)

func doSomethingForCircle(pbr *ProgressCircle, interval int) {
	for range 10 {
		pbr.Increment()
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
	pbr.ForceDone()
}

func TestPbrCircleSync(t *testing.T) {
	themes := []ProgressCircleTheme{THEME_CIRCLE_LIGHT, THEME_CIRCLE_CONCISE}
	for i := range len(themes) {
		pbr := NewProgressCircle(fmt.Sprintf("并发创建的进度环%d", i), themes[i])
		go doSomethingForCircle(pbr, i+1)
	}
	WaitAllProgressBar()
}
func TestPbrsMix(t *testing.T) {
	themes := []ProgressCircleTheme{THEME_CIRCLE_LIGHT, THEME_CIRCLE_CONCISE}
	for i := range len(themes) {
		pbr := NewProgressCircle(fmt.Sprintf("并发创建的进度环%d", i), themes[i])
		go doSomethingForCircle(pbr, (i+1)*100)
	}
	// WaitAllProgressBar()
	themes2 := []ProgressBarTheme{THEME_LIGHT, THEME_CONCISE}
	for i := range len(themes2) {
		pbr := NewProgressLinear(100, fmt.Sprintf("并发创建的进度条%d", i), themes2[i])
		go doSomething(pbr, (i+1)*10)
	}
	WaitAllProgressBar()
}

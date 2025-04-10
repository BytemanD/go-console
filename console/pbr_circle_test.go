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
	// PkgEnableLog()
	// EnableLogDebug()
	themes := []ProgressCircleTheme{THEME_CIRCLE_LIGHT, THEME_CIRCLE_CONCISE}
	for i := range len(themes) {
		pbr := NewProgressCircle(fmt.Sprintf("并发创建的进度环%d", i), themes[i])
		go doSomethingForCircle(pbr, (i+1)*100)
	}
	WaitAllProgressBar()
}

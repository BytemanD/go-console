package console

import (
	"fmt"
	"os"
	"sync"

	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
	"golang.org/x/term"
)

type PbrGroup struct {
	items       []ProgressBar
	titleLength int
	waitGroup   *sync.WaitGroup
}

func (m *PbrGroup) Add(pbr ProgressBar) {
	m.items = append(m.items, pbr)
	m.titleLength = max(m.titleLength, runewidth.StringWidth(pbr.GetTitle()))
}

func (m *PbrGroup) Reset() {
	m.items = []ProgressBar{}
	m.titleLength = 0
}

func (m *PbrGroup) Output() {
	if len(m.items) == 0 {
		return
	}
	//获取终端宽度
	progressWidth := 100
	if width, _, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
		progressWidth = lo.Min([]int{200, width - 8})
	}
	// 计算进度条长度
	progressWidth = progressWidth - m.titleLength - 10

	fmt.Print("\033[2K\r")
	completed := lo.Reduce(m.items, func(agg int, pbr ProgressBar, _ int) int {
		fmt.Println("\033[2K\r", pbr.Render(m.titleLength, progressWidth))
		if pbr.IsDone() {
			return agg + 1
		} else {
			return agg
		}
	}, 0)

	if completed < len(m.items) {
		fmt.Printf("\033[%dA", len(m.items))
	} else {
		m.Reset()
	}
}

var defaultPbrGroup *PbrGroup

func ProgressCount() int {
	return len(defaultPbrGroup.items)
}

func WaitAllProgressBar() {
	defaultPbrGroup.waitGroup.Wait()
}

func init() {
	defaultPbrGroup = &PbrGroup{
		items:     []ProgressBar{},
		waitGroup: &sync.WaitGroup{},
	}
}

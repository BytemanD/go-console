package console

import (
	"fmt"
	"os"
	"sync"

	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
	"golang.org/x/term"
)

var pbrSync = lo.Synchronize()

type ProgressLinear struct {
	Total     int
	Title     string
	completed int
	forceDone bool
	theme     ProgressBarTheme
}

func (p *ProgressLinear) Increment() {
	p.IncrementN(1)
}
func (p *ProgressLinear) IncrementN(n int) {
	if p.IsDone() {
		return
	}
	p.completed += n
	withOutputLock(func() {})
	if p.completed >= p.Total {
		defaultPbrGroup.waitGroup.Done()
	}
}
func (p *ProgressLinear) IsDone() bool {
	return p.completed >= p.Total || p.forceDone
}
func (p *ProgressLinear) ForceDone() {
	defer defaultPbrGroup.waitGroup.Done()

	p.forceDone = true
	if enableLog {
		DebugS("force progress done", "pkg", "console", "title", p.Title)
	}
}

func (p *ProgressLinear) Percent() float64 {
	return float64(p.completed*100) / float64(p.Total)
}
func (p *ProgressLinear) Render(titleWidth, progressWidth int) string {
	return p.theme.Render(p.Title, p.Percent(), titleWidth, progressWidth)
}

func NewProgressLinear(total int, title string, opt ...ProgressBarTheme) *ProgressLinear {
	theme := THEME_LIGHT
	if len(opt) > 0 {
		theme = opt[0]
	}
	var pbr *ProgressLinear
	pbrSync.Do(func() {
		if enableLog {
			DebugS("new progress bar", "pkg", "console", "title", title)
		}
		pbr = &ProgressLinear{Total: total, Title: title, theme: theme}
		defaultPbrGroup.waitGroup.Add(1)
		defaultPbrGroup.Add(pbr)
	})
	return pbr
}

type PbrGroup struct {
	items       []*ProgressLinear
	titleLength int
	waitGroup   *sync.WaitGroup
}

func (m *PbrGroup) Add(pbr *ProgressLinear) {
	m.items = append(m.items, pbr)
	m.titleLength = max(m.titleLength, runewidth.StringWidth(pbr.Title))
}

func (m *PbrGroup) Reset() {
	m.items = []*ProgressLinear{}
	m.titleLength = 0
}

func (m *PbrGroup) Output() {
	if len(m.items) == 0 {
		return
	}
	//获取终端宽度
	progressWidth := 100
	if width, _, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
		progressWidth = width
	}
	// 计算进度条长度
	progressWidth = progressWidth - m.titleLength - 10

	fmt.Print("\033[2K\r")
	completed := lo.Reduce(m.items, func(agg int, pbr *ProgressLinear, _ int) int {
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
		items:     []*ProgressLinear{},
		waitGroup: &sync.WaitGroup{},
	}
}

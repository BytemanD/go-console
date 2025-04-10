package console

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
	"golang.org/x/term"
)

// var pbrMu *sync.Mutex
var pbrWaitGroup *sync.WaitGroup

var pbrSync = lo.Synchronize()

type Pbr struct {
	Total     int
	Title     string
	completed int
	done      bool
	theme     PbrThemeInterface
}

func (p *Pbr) Increment() {
	p.IncrementN(1)
}
func (p *Pbr) IncrementN(n int) {
	if p.IsDone() {
		return
	}
	p.completed += n
	if p.completed >= p.Total {
		p.Done()
	}
	withOutputLock(func() {
		pbrManager.Output()
	})
}
func (p *Pbr) IsDone() bool {
	return p.done
}
func (p *Pbr) Done() {
	if p.IsDone() {
		return
	}
	p.done = true
	if enableLog {
		DebugS("task done", "pkg", "console", "title", p.Title)
	}
	pbrWaitGroup.Done()
}

func (p *Pbr) Percent() float64 {
	return float64(p.completed*100) / float64(p.Total)
}
func NewPbr(total int, title string) *Pbr {
	return NewPbrWithTheme(total, title, THEME_DEFAULT)
}
func NewPbrWithTheme(total int, title string, theme PbrThemeInterface) *Pbr {
	var pbr *Pbr
	pbrSync.Do(func() {
		if enableLog {
			DebugS("new pbr", "pkg", "console", "title", title)
		}
		pbr = &Pbr{Total: total, Title: title, theme: theme}
		pbrWaitGroup.Add(1)
		pbrManager.Add(pbr)
	})
	return pbr
}

type PbrThemeInterface interface {
	Render(pbr Pbr, titleLength, progrssLength int) string
}

type PbrTheme struct {
	Char string
}

func (t PbrTheme) fixTitle(pbr Pbr, titleLength int) string {
	return runewidth.FillRight(pbr.Title, titleLength) + ":"
}

func (t PbrTheme) Render(pbr Pbr, titleLength, progrssLength int) string {
	// 计算百分比
	percent := pbr.Percent()
	fixedProgressLength := int(percent) * progrssLength / 100
	progressStr := strings.Repeat(t.Char, max(fixedProgressLength, 0))
	if pbr.IsDone() {
		progressStr = color.GreenString(progressStr)
	} else {
		progressStr = color.RedString(progressStr)
	}
	progressStr += color.WhiteString(strings.Repeat(t.Char, progrssLength-int(fixedProgressLength)))
	return fmt.Sprintf("%s %s %3.2f%%", t.fixTitle(pbr, titleLength), progressStr, percent)
}

var (
	// 标题 : ━━━━━━━━━━━━━━━━━━━━━━ 100.00%
	THEME_DEFAULT = PbrTheme{Char: "━"}
	// 标题 : ====================== 100.00%
	THEME_SIMPLE = PbrTheme{Char: "="}
)

type Manager struct {
	items []*Pbr
	// theme       PbrThemeInterface
	titleLength int
}

func (m *Manager) Add(pbr *Pbr) {
	m.items = append(m.items, pbr)
	m.titleLength = max(m.titleLength, runewidth.StringWidth(pbr.Title))
}

func (m *Manager) Reset() {
	m.items = []*Pbr{}
	m.titleLength = 0
}

func (m *Manager) Output() {
	if len(m.items) == 0 {
		return
	}
	//获取终端宽度
	progrssLength := 100
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err == nil {
		progrssLength = width
	}
	// 计算进度条长度
	progrssLength = progrssLength - m.titleLength - 10

	fmt.Print("\033[2K\r")
	completed := 0
	for _, pbr := range m.items {
		fmt.Print("\033[2K\r")
		fmt.Println(pbr.theme.Render(*pbr, m.titleLength, progrssLength))
		if pbr.IsDone() {
			completed += 1
		}
	}
	if completed < len(m.items) {
		for range m.items {
			fmt.Print("\033[1A")
		}
	} else {
		m.Reset()
	}
}

var pbrManager *Manager

func GetPbrNum() int {
	return len(pbrManager.items)
}

func WaitAllProgressBar() {
	pbrWaitGroup.Wait()
}

func CustomeTheme(c string) PbrTheme {
	return PbrTheme{Char: c}
}

func init() {
	pbrWaitGroup = &sync.WaitGroup{}

	pbrManager = &Manager{}
	pbrManager.Reset()
}

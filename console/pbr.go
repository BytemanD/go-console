package console

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

var pbrMu *sync.Mutex
var pbrChan chan bool

type Pbr struct {
	Total     int
	Title     string
	completed int
	done      bool
	mu        *sync.Mutex
}

func (p *Pbr) Increment() {
	p.completed += 1
	withOutputLock(func() {
		pbrManager.Output()
	})
}
func (p *Pbr) IncrementN(n int) {
	p.completed += n
	withOutputLock(func() {
		pbrManager.Output()
	})
}
func (p *Pbr) IsDone() bool {
	if p.done {
		return true
	}
	return p.completed >= p.Total
}
func (p *Pbr) Done() {
	p.done = true
}

func (p *Pbr) Percent() float64 {
	return float64(p.completed*100) / float64(p.Total)
}
func NewPbr(total int, title string) *Pbr {
	pbrMu.Lock()
	defer pbrMu.Unlock()

	pbr := &Pbr{Total: total, Title: title, mu: &sync.Mutex{}}
	pbrManager.Add(pbr)
	return pbr
}

type PbrThemeInterface interface {
	Render(pbr Pbr, titleLength, progrssLength int) string
}

type PbrDefaultTheme struct {
	Char string
}

func (t PbrDefaultTheme) fixTitle(pbr Pbr, titleLength int) string {
	return runewidth.FillRight(pbr.Title, titleLength) + ":"
}

func (t PbrDefaultTheme) Render(pbr Pbr, titleLength, progrssLength int) string {
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

type Manager struct {
	items       []*Pbr
	theme       PbrThemeInterface
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
	pbrMu.Lock()
	defer pbrMu.Unlock()

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
		fmt.Println(m.theme.Render(*pbr, m.titleLength, progrssLength))
		if pbr.IsDone() {
			completed += 1
		}
	}
	if completed < len(m.items) {
		for i := 0; i < len(m.items); i++ {
			fmt.Print("\033[1A")
		}
	} else {
		m.Reset()
		pbrChan <- true
	}
}

var pbrManager *Manager

func SetPbrTheme(theme PbrThemeInterface) {
	pbrManager.theme = theme
}

func GetPbrNum() int {
	return len(pbrManager.items)
}

func WaitPbrs() {
	<-pbrChan
}

func init() {
	pbrMu = &sync.Mutex{}
	pbrChan = make(chan bool)

	pbrManager = &Manager{theme: PbrDefaultTheme{Char: "━"}}
	pbrManager.Reset()
}

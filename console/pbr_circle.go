package console

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
)

type ProgressCircleTheme struct {
	Chars []string
	index int
	tmpl  *template.Template
}

func (p *ProgressCircleTheme) Render(title string, titleWidth int, progressWidth int) string {
	defer func() { p.index += 1 }()

	if p.index >= len(p.Chars) {
		p.index = 0
	}
	var buf bytes.Buffer
	if err := p.tmpl.Execute(&buf, pbrData{
		Title:    runewidth.FillRight(title, titleWidth),
		Progress: p.Chars[p.index],
	}); err != nil {
		return fmt.Sprintf("<!parse template failed: %s>", err)
	}
	return buf.String()

}

type ProgressCircle struct {
	Title     string
	forceDone bool
	theme     ProgressCircleTheme
}

func (p *ProgressCircle) GetTitle() string {
	return p.Title
}
func (p *ProgressCircle) Increment() {
	p.IncrementN(1)
}
func (p *ProgressCircle) IncrementN(n int) {
	if p.IsDone() {
		return
	}
	withOutputLock(func() {})
}
func (p *ProgressCircle) IsDone() bool {
	return p.forceDone
}
func (p *ProgressCircle) ForceDone() {
	if p.IsDone() {
		return
	}
	defer defaultPbrGroup.waitGroup.Done()

	p.forceDone = true
	if enableLog {
		DebugS("force progress done", "pkg", "console", "title", p.Title)
	} else {
		withOutputLock(func() {})
	}
}

func (p *ProgressCircle) Render(titleWidth, progressWidth int) string {
	return p.theme.Render(p.Title, titleWidth, progressWidth)
}

func NewProgressCircleTheme(chars []string, templateArg ...string) *ProgressCircleTheme {
	tmplStr := lo.CoalesceOrEmpty(lo.FirstOrEmpty(templateArg), DEFAULT_PBR_TEMPLATE)
	tmpl, err := template.New("pbr").Parse(tmplStr)
	if err != nil {
		tmpl, _ = template.New("pbr").Parse(DEFAULT_PBR_TEMPLATE)
	}
	return &ProgressCircleTheme{Chars: chars, tmpl: tmpl}
}

func NewProgressCircle(title string, opt ...ProgressCircleTheme) *ProgressCircle {
	theme := THEME_CIRCLE_LIGHT
	if len(opt) > 0 {
		theme = opt[0]
	}
	var pbr *ProgressCircle
	pbrSync.Do(func() {
		if enableLog {
			DebugS("new progress circle", "pkg", "console", "title", title)
		}
		pbr = &ProgressCircle{Title: title, theme: theme}
		defaultPbrGroup.waitGroup.Add(1)
		defaultPbrGroup.Add(pbr)
	})
	return pbr
}

var (
	// 标题 : ⠇ 100.00%
	THEME_CIRCLE_LIGHT   = *NewProgressCircleTheme(strings.Split("⠇⠋⠉⠙⠸⠴⠤⠦", ""), DEFAULT_PBR_TEMPLATE)
	THEME_CIRCLE_CONCISE = *NewProgressCircleTheme(strings.Split("|/-\\", ""), DEFAULT_PBR_TEMPLATE)
)

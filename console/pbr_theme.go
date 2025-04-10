package console

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
)

var DEFAULT_PBR_TEMPLATE = "{{.Title}}: {{.Progress}} {{.Percent}}%"

type pbrData struct {
	Title    string
	Progress string
	Percent  string
}

type ProgressBarTheme struct {
	Char      string
	topChar   string
	tmpl      *template.Template
	fillRight bool
}

func (p *ProgressBarTheme) SetTopchar(topChar string) *ProgressBarTheme {
	p.topChar = topChar
	p.fillRight = true
	return p
}
func (p *ProgressBarTheme) Topchar() string {
	return lo.CoalesceOrEmpty(p.topChar, p.Char)
}
func (p ProgressBarTheme) Render(title string, percent float64, titleWidth int, progressWidth int) string {
	fixedProgressLength := int(percent) * progressWidth / 100
	progressStr := strings.Repeat(p.Char, max(fixedProgressLength, 0)) + p.topChar
	if percent >= 100 {
		progressStr = color.GreenString(progressStr)
	} else {
		progressStr = color.RedString(progressStr)
	}
	if p.fillRight {
		progressStr += color.WhiteString(strings.Repeat(" ", progressWidth-int(fixedProgressLength)))
	} else {
		progressStr += color.WhiteString(strings.Repeat(p.Char, progressWidth-int(fixedProgressLength)))
	}

	var buf bytes.Buffer
	if err := p.tmpl.Execute(&buf, pbrData{
		Title:    runewidth.FillRight(title, titleWidth),
		Progress: progressStr,
		Percent:  runewidth.FillLeft(fmt.Sprintf("%.2f", percent), 6),
	}); err != nil {
		return fmt.Sprintf("<!parse template failed: %s>", err)
	}
	return buf.String()
}

func NewProgressBarTheme(char string, templateArg ...string) *ProgressBarTheme {
	tmplStr := lo.CoalesceOrEmpty(lo.FirstOrEmpty(templateArg), DEFAULT_PBR_TEMPLATE)
	tmpl, err := template.New("pbr").Parse(tmplStr)
	if err != nil {
		tmpl, _ = template.New("pbr").Parse(DEFAULT_PBR_TEMPLATE)
	}
	return &ProgressBarTheme{Char: char, tmpl: tmpl}
}

var (
	// 标题 : ━━━━━━━━━━━━━━━━━━━━━━ 100.00%
	THEME_LIGHT = *NewProgressBarTheme("━", DEFAULT_PBR_TEMPLATE)
	// 标题 : ====================== 100.00%
	THEME_CONCISE = *NewProgressBarTheme("=", DEFAULT_PBR_TEMPLATE)
	// 标题 : ☁️☁️☁️☁️🛬             100.00%
	THEME_EMOJI_AIRPLANE = *NewProgressBarTheme("☁️", DEFAULT_PBR_TEMPLATE).SetTopchar("🛬")
	// 标题 : ❤️❤️❤️❤️❤️             100.00%
	THEME_EMOJI_HEART = *NewProgressBarTheme("❤️", DEFAULT_PBR_TEMPLATE).SetTopchar("❤️")
)

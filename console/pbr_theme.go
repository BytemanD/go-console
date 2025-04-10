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
	Char string
	tmpl *template.Template
}

func (p ProgressBarTheme) Render(title string, percent float64, titleWidth int, progressWidth int) string {
	fixedProgressLength := int(percent) * progressWidth / 100
	progressStr := strings.Repeat(p.Char, max(fixedProgressLength, 0))
	if percent >= 100 {
		progressStr = color.GreenString(progressStr)
	} else {
		progressStr = color.RedString(progressStr)
	}
	progressStr += color.WhiteString(strings.Repeat(p.Char, progressWidth-int(fixedProgressLength)))

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

func NewProgressBarTheme(char string, templateArg ...string) ProgressBarTheme {
	tmplStr := lo.CoalesceOrEmpty(lo.FirstOrEmpty(templateArg), DEFAULT_PBR_TEMPLATE)
	tmpl, err := template.New("pbr").Parse(tmplStr)
	if err != nil {
		tmpl, _ = template.New("pbr").Parse(DEFAULT_PBR_TEMPLATE)
	}
	return ProgressBarTheme{Char: char, tmpl: tmpl}
}

var (
	// 标题 : ━━━━━━━━━━━━━━━━━━━━━━ 100.00%
	THEME_LIGHT = NewProgressBarTheme("━", DEFAULT_PBR_TEMPLATE)
	// 标题 : ====================== 100.00%
	THEME_CONCISE = NewProgressBarTheme("=", DEFAULT_PBR_TEMPLATE)
)

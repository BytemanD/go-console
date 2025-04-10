package console

type ProgressLinear struct {
	Total     int
	Title     string
	completed int
	forceDone bool
	theme     ProgressBarTheme
}

func (p ProgressLinear) GetTitle() string {
	return p.Title
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

func (p *ProgressLinear) Percent() float64 {
	return float64(p.completed*100) / float64(p.Total)
}
func (p ProgressLinear) Render(titleWidth, progressWidth int) string {
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

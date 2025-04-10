package console

import "github.com/samber/lo"

var pbrSync = lo.Synchronize()

type ProgressBar interface {
	IsDone() bool
	GetTitle() string
	Render(titleWidth, progressWidth int) string
}

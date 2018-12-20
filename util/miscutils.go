package util

import (
	"github.com/schollz/progressbar"
)

var ProgressBarOptions = progressbar.OptionSetTheme(progressbar.Theme{
	Saucer:        "▒",
	SaucerHead:    "▒",
	SaucerPadding: " ",
	BarStart:      "|",
	BarEnd:        "|",
})

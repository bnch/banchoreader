package banchoreader

import (
	"github.com/fatih/color"
)

var green *color.Color
var yellow *color.Color
var blue *color.Color
var def *color.Color

func colors() {
	green = color.New(color.Bold, color.FgGreen)
	yellow = color.New(color.FgYellow)
	blue = color.New(color.FgBlue)
	def = color.New(color.Reset)
}

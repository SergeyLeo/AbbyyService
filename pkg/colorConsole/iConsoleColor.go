package colorConsole

import (
	"github.com/fatih/color"
)

type ColorMessage uint8

const (
	_ ColorMessage = 1 << iota
	ColorTypeInfo
	ColorTypeInfoEffect
	ColorTypeWarning
	ColorTypeError
)

var wConsole *color.Color
var eConsole *color.Color

func InstanceColor(colorType ColorMessage) *color.Color {
	switch colorType {
	case ColorTypeInfo:
		return color.New(color.FgGreen)
	case ColorTypeInfoEffect:
		return color.New(color.BgGreen, color.FgBlack, color.Bold)
	case ColorTypeWarning:
		if wConsole == nil {
			return color.New(color.FgYellow)
		} else {
			return wConsole
		}
	case ColorTypeError:
		if eConsole == nil {
			return color.New(color.FgRed)
		} else {
			return eConsole
		}
	}
	return color.New()
}

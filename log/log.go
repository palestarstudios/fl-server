package log

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Info    = color.New(color.FgBlue)
	Success = color.New(color.FgGreen)
	Warning = color.New(color.FgYellow)
	Error   = color.New(color.FgRed, color.Bold)

	Player = color.New(color.FgMagenta)
	System = color.New(color.FgWhite)
)

func Log(prefix *color.Color, tag string, msg string) {
	prefix.Printf("[%s] ", tag)
	fmt.Println(msg)
}

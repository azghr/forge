package main

import (
	"fmt"
	"os"
)

type colorCode int

const (
	colorReset  colorCode = 0
	colorRed    colorCode = 31
	colorGreen  colorCode = 32
	colorYellow colorCode = 33
	colorBlue   colorCode = 34
	colorCyan   colorCode = 36
)

func colorize(c colorCode, s string) string {
	if os.Getenv("NO_COLOR") != "" {
		return s
	}
	return fmt.Sprintf("\033[%dm%s\033[0m", c, s)
}

func green(s string) string  { return colorize(colorGreen, s) }
func red(s string) string    { return colorize(colorRed, s) }
func yellow(s string) string { return colorize(colorYellow, s) }
func blue(s string) string   { return colorize(colorBlue, s) }
func cyan(s string) string   { return colorize(colorCyan, s) }

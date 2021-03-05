package main

import (
	"math/rand"

	"github.com/fatih/color"
)

var availableColors = []*color.Color{
	color.New(color.FgCyan),
	color.New(color.FgBlue),
	color.New(color.FgRed),
	color.New(color.FgGreen),
	color.New(color.FgYellow),
	color.New(color.FgMagenta),
}

func randomColor() *color.Color {
	i := rand.Int31n(int32(len(availableColors)))
	return availableColors[i]
}

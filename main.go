package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		// Creates a new window
		w := app.NewWindow(
			app.Title("Egg Timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)
		//ops are the operations from the UI
		var ops op.Ops

		//startButton is the clickable widget
		var startButton widget.Clickable

		//th defines the material design stype (Theme)
		th := material.NewTheme(gofont.Collection())

		// listin for events in the window we made
		for e := range w.Events() {
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				btn := material.Button(th, &startButton, "Start the timer")
				btn.Layout(gtx)
				e.Frame(gtx.Ops)
			}

		}
		os.Exit(0)

	}()
	app.Main()
}

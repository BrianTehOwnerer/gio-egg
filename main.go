package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"os"
	"time"
)

// Define the progress variables, a channel and a variable
var progressIncrementer chan float32
var progress float32

func main() {
	// Setup a separate channel to provide ticks to increment progress
	progressIncrementer = make(chan float32)
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressIncrementer <- 0.004
		}
	}()

	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Egg timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
			app.MinSize(unit.Dp(300), unit.Dp(400)),
			app.MaxSize(unit.Dp(800), unit.Dp(1000)),
		)
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

type C = layout.Context
type D = layout.Dimensions

func draw(w *app.Window) error {
	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable

	// is the egg boiling?
	var boiling bool

	// th defines the material design style
	th := material.NewTheme(gofont.Collection())

	for {
		select {
		// listen for events in the window.
		case e := <-w.Events():

			// detect what type of event
			switch e := e.(type) {

			// this is sent when the application should re-render.
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				// Let's try out the flexbox layout concept
				if startButton.Clicked() {
					boiling = !boiling
				}

				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					layout.Rigid(
						func(gtx C) D {
							circle := clip.Ellipse{
								//Hard coding the X cordinate
								//Min: image.Pt(80, 0),
								//Max: image.Pt(320, 240),
								//Softcodeing the x cord
								Min: image.Pt(gtx.Constraints.Max.X/2-120, 0),
								Max: image.Pt(gtx.Constraints.Max.X/2+120, 240),
							}.Op(gtx.Ops)
							color := color.NRGBA{R: 200, A: 255}
							paint.FillShape(gtx.Ops, color, circle)
							d := image.Point{Y: 400}
							return layout.Dimensions{Size: d}
						}), layout.Rigid(
						func(gtx C) D {
							bar := material.ProgressBar(th, progress)
							return bar.Layout(gtx)
						},
					),
					layout.Rigid(
						func(gtx C) D {
							// We start by defining a set of margins
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(35),
								Left:   unit.Dp(35),
							}
							// Then we lay out within those margins ...
							return margins.Layout(gtx,
								// ...the same function we earlier used to create a button
								func(gtx C) D {
									var text string
									if !boiling {
										text = "Start the timer"
									} else {
										text = "Pause the timer"
									}
									btn := material.Button(th, &startButton, text)
									return btn.Layout(gtx)
								},
							)
						},
					),
				)
				e.Frame(gtx.Ops)

			// this is sent when the application is closed.
			case system.DestroyEvent:
				return e.Err
			}

		// listen for events in the incrementor channel
		case p := <-progressIncrementer:
			if boiling && progress < 1 {
				progress += p
				w.Invalidate()
			}
		}
	}
}

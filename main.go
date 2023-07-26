package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Define the progress variables, a channel and a variable
var progressIncrementer chan bool
var progress float32

func main() {
	// Setup a separate channel to provide ticks to increment progress
	progressIncrementer = make(chan bool)
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressIncrementer <- true
		}
	}()

	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Egg timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
			app.MaxSize(unit.Dp(600), unit.Dp(800)),
			app.MinSize(unit.Dp(300), unit.Dp(400)),
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
	var boilDurationInput widget.Editor
	var boilDuration float32

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

					//reset the boil
					if progress >= 1 {
						progress = 0
					}

					//here we read from the input box
					inputString := boilDurationInput.Text()
					inputString = strings.TrimSpace(inputString)
					inputFloat, _ := strconv.ParseFloat(inputString, 32)
					boilDuration = float32(inputFloat)
					boilDuration = boilDuration / (1 - progress)
				}

				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,

					layout.Rigid(
						func(gtx C) D {
							// Draw a custom path, shaped like an egg
							var eggPath clip.Path
							op.Offset(image.Pt(gtx.Dp(200), gtx.Dp(150))).Add(gtx.Ops)
							eggPath.Begin(gtx.Ops)
							// Rotate from 0 to 360 degrees
							for deg := 0.0; deg <= 360; deg++ {

								// Egg math (really) at this brilliant site. Thanks!
								// https://observablehq.com/@toja/egg-curve
								// Convert degrees to radians
								rad := deg / 360 * 2 * math.Pi
								// Trig gives the distance in X and Y direction
								cosT := math.Cos(rad)
								sinT := math.Sin(rad)
								// Constants to define the eggshape
								a := 110.0
								b := 150.0
								d := 20.0
								// The x/y coordinates
								x := a * cosT
								y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
								// Finally the point on the outline
								p := f32.Pt(float32(x), float32(y))
								// Draw the line to this point
								eggPath.LineTo(p)
							}
							// Close the path
							eggPath.Close()

							// Get hold of the actual clip
							eggArea := clip.Outline{Path: eggPath.End()}.Op()

							// Fill the shape
							// color := color.NRGBA{R: 255, G: 239, B: 174, A: 255}
							color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1 - progress)), A: 255}
							paint.FillShape(gtx.Ops, color, eggArea)

							d := image.Point{Y: 375}
							return layout.Dimensions{Size: d}
						},
					),
					layout.Rigid(
						func(gtx C) D {
							ed := material.Editor(th, &boilDurationInput, "sec")

							boilDurationInput.SingleLine = true
							boilDurationInput.Alignment = text.Middle
							if boiling && progress < 1 {
								boilRemain := (1 - progress) * boilDuration
								//formatting to make things look nice, using the
								// multiply by 10 devide by 10 trick
								inputStr := fmt.Sprintf("%.1f", math.Round(math.Abs(float64(boilRemain)*10))/10)
								//update the text inside of the input box
								boilDurationInput.SetText(inputStr)
							}

							//layouts for the text box
							margins := layout.Inset{
								Top:    unit.Dp(0),
								Right:  unit.Dp(170),
								Bottom: unit.Dp(40),
								Left:   unit.Dp(170),
							}
							border := widget.Border{
								Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
								CornerRadius: unit.Dp(3),
								Width:        unit.Dp(2),
							}
							return margins.Layout(gtx,
								func(gtx C) D {
									return border.Layout(gtx, ed.Layout)
								})
						},
					),
					//Progress Bar happens here
					layout.Rigid(
						func(gtx C) D {
							bar := material.ProgressBar(th, progress)
							return bar.Layout(gtx)
						},
					),
					//this is the bounding box for the start button
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
										text = "Start"
									}
									if boiling && progress < 1 {
										text = "Stop"
									}
									if boiling && progress >= 1 {
										text = "Egg Boiled"
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
		case <-progressIncrementer:
			if boiling && progress < 1 {
				progress += float32(1.0 / 25.0 / math.Abs(float64(boilDuration)))
				if progress >= 1 {
					progress = 1
				}
				w.Invalidate()
			}
		}
	}
}

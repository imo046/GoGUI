package src

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// TODO: use default state to return params after clicking "Finished"
type State struct {
	progress   float32
	duration   float32
	processing bool
	text       string
}

func Draw(w *app.Window) error {

	//operations from the UI
	var ops op.Ops

	//clickable widget
	var startButton widget.Clickable

	//editor editorText widget
	var durationInput widget.Editor

	//progress params
	var progressInc = make(chan float32)
	defer close(progressInc)

	var progress float32

	var duration float32

	//runs in the background
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressInc <- 0.004
		}
	}()

	//state
	var processing bool

	//th defines the Material Design style
	th := material.NewTheme(gofont.Collection())

	//event listener, gets us the channel
	for {

		select {

		case e := <-w.Events():

			//detect event type
			switch e := e.(type) {

			case system.FrameEvent:
				//specify new graphical context
				gtx := layout.NewContext(&ops, e)

				//event
				if startButton.Clicked() {
					userInput := durationInput.Text()
					userInput = strings.TrimSpace(userInput)
					floatInput, err := strconv.ParseFloat(userInput, 32)
					if err != nil {
						log.Fatal(err)
					}
					duration = float32(floatInput)
					duration = duration / (1 - progress)

					processing = !processing
					//TODO: revert to default at click event
					//if processing && progress >= 1 {
					//}
				}
				layout.Flex{
					//Vertical alignment, from top to bottom
					Axis: layout.Vertical,

					//Empty space at the start (top)
					Spacing: layout.SpaceStart,
				}.Layout(gtx,

					//Figure
					layout.Rigid(
						func(gtx C) D {
							circle := Circle{
								Min: image.Point{gtx.Constraints.Max.X/2 - 120, 0},
								Max: image.Point{gtx.Constraints.Max.X/2 + 120, 240},
							}
							circle.ChangeColor(progress)
							return circle.Draw(gtx)
						},
					),
					//EditorText input field
					layout.Rigid(

						func(gtx C) D {
							// Wrap the editor in material design
							ed := TextField{th, &durationInput, "sec"}
							durationInput.SingleLine = true
							durationInput.Alignment = text.Middle

							if processing && progress < 1 {
								remainTime := (1 - progress) * duration
								inputStr := fmt.Sprintf("%.1f", math.Round(float64(remainTime)*10)/10)
								//update
								durationInput.SetText(inputStr)
							}

							// Define insets
							margins := layout.Inset{
								Top:    unit.Dp(0),
								Right:  unit.Dp(170),
								Bottom: unit.Dp(40),
								Left:   unit.Dp(170),
							}
							// Borders
							border := widget.Border{
								Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
								CornerRadius: unit.Dp(3),
								Width:        unit.Dp(2),
							}
							// Draw them
							return margins.Layout(gtx,
								func(gtx C) D {
									return border.Layout(gtx, ed.Draw)
								},
							)
						},
					),
					//Progress bar
					layout.Rigid(
						func(gtx C) D {
							return ProgressBar{th: th, progress: progress}.Draw(gtx)
						},
					),

					//Button
					layout.Rigid(
						func(gtx C) D {

							var editorText string
							if !processing {
								editorText = "Start"
							}
							if processing && progress < 1 {
								editorText = "Stop"
							}
							if processing && progress >= 1 {
								editorText = "Finished"
							}

							//Define margins around the button with margins
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(35),
								Left:   unit.Dp(35),
							}
							//Layout those margins
							return margins.Layout(gtx,
								//Within margins, define and layout the button
								func(gtx C) D {
									//place button (Button is a widget, something which returns its own dimensions)
									return Button{th: th, btn: &startButton, text: editorText}.Draw(gtx)
								},
							)
						},
					),

					//Empty spacer, goes after the button (under it on the flex box)
					layout.Rigid(
						Spacer{Width: unit.Dp(25), Height: unit.Dp(25)}.Draw,
					),
				)
				//run operations
				e.Frame(gtx.Ops)

			case system.DestroyEvent:
				return e.Err
			}
		//TODO: rewrite to separate progress channel
		case p := <-progressInc:
			if processing && progress < 1 {
				progress += p
				//redraw
				w.Invalidate()
			}

		}

	}
	//return nil
}

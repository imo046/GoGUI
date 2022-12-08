package src

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"log"
	"strconv"
	"strings"
	"time"
)

func Draw(w *app.Window) error {

	//operations from the UI
	var ops op.Ops

	//clickable widget
	var startButton widget.Clickable

	//editor text widget
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
				}
				layout.Flex{
					//Vertical alignment, from top to bottom
					Axis: layout.Vertical,

					//Empty space at the start (top)
					Spacing: layout.SpaceStart,
				}.Layout(gtx,

					//Figure
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							circle := Circle{
								Min: image.Point{gtx.Constraints.Max.X/2 - 120, 0},
								Max: image.Point{gtx.Constraints.Max.X/2 + 120, 240},
							}
							circle.ChangeColor(progress)
							return circle.Draw(gtx)
						},
					),
					//placeholder for the text input field
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{}
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

							var text string
							if !processing {
								text = "Start"
							} else {
								text = "Stop"
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
									return Button{th: th, btn: &startButton, text: text}.Draw(gtx)
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

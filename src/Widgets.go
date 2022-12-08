package src

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	color2 "image/color"
)

type Widget interface {
	Draw(gtx C) D
}

type Spacer struct {
	Height unit.Dp
	Width  unit.Dp
}

type Button struct {
	th   *material.Theme
	btn  *widget.Clickable
	text string
}

type ProgressBar struct {
	th       *material.Theme
	progress float32
}

type Circle struct {
	Min image.Point
	Max image.Point
}

func (c Circle) Draw(gtx C) D {
	circleObj := clip.Ellipse{
		c.Min,
		c.Max,
	}.Op(gtx.Ops)
	color := color2.NRGBA{R: 200, A: 255}
	paint.FillShape(gtx.Ops, color, circleObj)
	d := image.Point{Y: 400}
	return layout.Dimensions{
		Size: d,
	}
}

func (b ProgressBar) Draw(gtx C) D {
	barObj := material.ProgressBar(b.th, b.progress)
	return barObj.Layout(gtx)
}

func (b Button) Draw(gtx C) D {
	btnObj := material.Button(b.th, b.btn, b.text)
	return btnObj.Layout(gtx)
}

func (s Spacer) Draw(gtx C) D {
	spacerObj := layout.Spacer{Width: s.Width, Height: s.Height}
	return spacerObj.Layout(gtx)
}

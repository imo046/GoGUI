package src

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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

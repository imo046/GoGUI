package main

import (
	"GoGUI/src"
	"gioui.org/app"
	"gioui.org/unit"
	"log"
	"os"
)

func main() {
	//must be called last inside the 'main' function
	defer app.Main()

	go func() {
		w := app.NewWindow(
			app.Title("Timer"),
			app.Size(unit.Dp(400), unit.Dp(600)), //device independent pixels (x,y)
		)

		if err := src.Draw(w); err != nil {
			log.Panic(err)
		}
		os.Exit(0)
	}()
}

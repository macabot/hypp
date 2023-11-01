package main

import (
	"github.com/macabot/hypp/examples/mouse-drag-and-drop/app"
	_ "github.com/macabot/hypp/jsd"
	"github.com/macabot/hypp/window"
)

func main() {
	el := window.Document().GetElementById("app")
	if el.IsNull() {
		panic("Could not find element with id 'app'")
	}
	app.Run(
		el,
	)

	select {} // run Go forever
}

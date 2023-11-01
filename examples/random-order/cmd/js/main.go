package main

import (
	"github.com/macabot/hypp/examples/random-order/app"
	_ "github.com/macabot/hypp/jsd"
	"github.com/macabot/hypp/window"
)

// TODO why does this behave different from hyperapp.html?
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

package main

import (
	"github.com/macabot/hypp/examples/todo-list/app"
	"github.com/macabot/hypp/js"
	_ "github.com/macabot/hypp/jsd"
	"github.com/macabot/hypp/window"
)

func main() {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'")
	}
	app.Run(
		window.Element{Value: el},
	)

	select {} // run Go forever
}

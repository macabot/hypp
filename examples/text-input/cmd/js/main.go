package main

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/examples/todo-list/app"
	"github.com/macabot/hypp/js"
	_ "github.com/macabot/hypp/jsd"
)

func main() {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'")
	}
	app.Run(
		hypp.Element{Value: el},
	)

	select {} // run Go forever
}

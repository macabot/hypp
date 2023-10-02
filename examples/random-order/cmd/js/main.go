package main

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/examples/todo-list/app"
	_ "github.com/macabot/hypp/jsd"
)

// TODO why does this behave different from hyperapp.html?
func main() {
	el := hypp.Document().GetElementById("app")
	if el.IsNull() {
		panic("Could not find element with id 'app'")
	}
	app.Run(
		el,
	)

	select {} // run Go forever
}

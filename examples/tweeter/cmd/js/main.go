package main

import (
	"syscall/js"

	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/examples/tweeter/app"
)

func main() {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'")
	}
	app.Run(
		jsd.Driver{},
		jsd.Node(el),
	)

	select {} // run Go forever
}

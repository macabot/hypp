package main

import (
	"syscall/js"

	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/examples/todo-list/app"
)

func main() {
	app.Run(jsd.Node(js.Global().Get("document").Call("getElementById", "app")))
}

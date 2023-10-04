package main

import (
	"github.com/macabot/hypp/examples/gif-search/app"
	_ "github.com/macabot/hypp/js"
	"github.com/macabot/hypp/window"
)

// Build with -ldflags="-X 'main.APIKey=<api_key>'"
var APIKey string

func main() {
	el := window.Document().GetElementById("app")
	if el.IsNull() {
		panic("Could not find element with id 'app'")
	}
	app.APIKey = APIKey
	app.Run(
		el,
	)

	select {} // run Go forever
}

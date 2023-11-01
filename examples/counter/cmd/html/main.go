package main

import (
	"fmt"

	"github.com/macabot/hypp/examples/counter/app"
	"github.com/macabot/hypp/tag"
)

func main() {
	fmt.Println(tag.RenderToString(app.View(&app.State{})))
}

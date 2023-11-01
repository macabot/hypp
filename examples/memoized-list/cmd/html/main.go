package main

import (
	"fmt"

	"github.com/macabot/hypp/examples/memoized-list/app"
	"github.com/macabot/hypp/tag"
)

func main() {
	state := &app.State{
		List: []string{"a", "b", "c"},
	}
	fmt.Println(tag.RenderToString(app.View(state)))
}

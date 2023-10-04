package main

import (
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/examples/svg-circles/app"
	"github.com/macabot/hypp/tag"
)

func main() {
	fmt.Println(tag.RenderToString(app.View(&hypp.EmptyState{})))
}

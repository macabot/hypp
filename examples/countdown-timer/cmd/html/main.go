package main

import (
	"fmt"
	"time"

	"github.com/macabot/hypp/examples/countdown-timer/app"
	"github.com/macabot/hypp/tag"
)

func main() {
	state := &app.State{
		Count:  10 * time.Second,
		Paused: true,
	}
	fmt.Println(tag.RenderToString(app.View(state)))
}

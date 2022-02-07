package main

import (
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/driver/html"
	"github.com/macabot/hypp/examples/memoized-list/app"
)

func main() {
	driver := html.Driver{}
	node := driver.CreateElement("main", hypp.Option[hypp.ElementCreationOptions]{})
	app.Run(driver, node)
	fmt.Println(node.(*html.Node).InnerHTML(
		&html.RenderOptions{Deterministic: true},
	))
}

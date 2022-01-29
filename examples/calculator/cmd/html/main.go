package main

import (
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/driver/html"
	"github.com/macabot/hypp/examples/calculator/app"
)

// FIXME number buttons should not have class attribute with empty value.
func main() {
	driver := html.Driver{}
    // TODO remove things like ElementCreationOptions that were added to mimick the browser implementation, but aren't used by hypp.
	node := driver.CreateElement("main", hypp.Option[hypp.ElementCreationOptions]{})
	app.Run(driver, node)
	fmt.Println(node.(*html.Node).InnerHTML())
}

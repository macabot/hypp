// Package app draws SVG circles.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/preYMW

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"github.com/macabot/hypp/tag/svg"
)

func main(children ...*hypp.VNode) *hypp.VNode {
	return html.Main(nil, children...)
}

func circle(id string, props hypp.HProps) *hypp.VNode {
	props.Set("id", id)
	return svg.Circle(props)
}

func use(href string, props hypp.HProps) *hypp.VNode {
	props.Set("href", href)
	return svg.Use(props)
}

func Run(driver hypp.Driver, node hypp.Node) {
	hypp.App(hypp.AppProps[*hypp.EmptyState]{
		Driver: driver,
		Init:   &hypp.EmptyState{},
		View: func(_ *hypp.EmptyState) *hypp.VNode {
			return main(
				html.Svg(
					hypp.HProps{"viewBox": "0 0 30 10"},
					circle("symbol", hypp.HProps{
						"cx":     5,
						"cy":     5,
						"r":      4,
						"stroke": "#0366d6",
					}),
					use("#symbol", hypp.HProps{
						"x":    10,
						"fill": "#0366d6",
					}),
					use("#symbol", hypp.HProps{
						"x":    20,
						"fill": "white",
					}),
				),
			)
		},
		Node: node,
	})
}

// This file is based on https://codepen.io/jorgebucaran/pen/wjvEBj
package app

import (
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type MyState struct {
	hypp.EmptyState
	todos []TodoItem
	value string
}

func (m MyState) clone() *MyState {
	return &m
}

func container(children ...*hypp.VNode) *hypp.VNode {
	return html.Section(nil, children...)
}

func main(children ...*hypp.VNode) *hypp.VNode {
	return html.Main(nil, children...)
}

func button(onclick hypp.Dispatchable, text string) *hypp.VNode {
	return html.Button(hypp.HProps{"onclick": onclick}, hypp.Text(text))
}

func separator() *hypp.VNode {
	return html.Hr(nil)
}

type focusProps struct {
	id string
	preventScroll interface{} // TODO
}

func justFocus(_ *MyState, payload hypp.Payload) {

}

func focus(id string, preventScroll) hypp.Effect {
	return hypp.Effect{
		Effecter: justFocus,
		Payload: focusProps{id: id, preventScroll: preventScroll},
	}
}

func domID(n int) string {
	return fmt.Sprintf("input-%d", n)
}

func Run(driver hypp.Driver, node hypp.Node) {
	hypp.App(hypp.AppProps[*MyState]{
		Driver: driver,
		Init: &MyState{},
		View: func(state *MyState) *hypp.VNode {
			// TODO
		},
		Node: node,
	})
}

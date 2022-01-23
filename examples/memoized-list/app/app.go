// Package app lets you write 140 characters in a textarea.
// Click "Tweet" to empty the textarea.
package app

// This file is based on https://github.com/jorgebucaran/hyperapp/blob/7a5c5c8e1e92387ab7295daf5bf2448490d23eb6/docs/api/memo.md#example

import (
	"math"
	"math/rand"
	"strings"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type List []string

func (l List) Hash() string {
	return strings.Join(l, "")
}

type MyState struct {
	hypp.EmptyState
	list    List
	counter int
}

func (m MyState) clone() *MyState {
	return &m
}

var hex = "0123456789ABCDEF"

func randomHex() string {
	i := int(math.Floor(rand.Float64() * 16))
	return hex[i : i+1]
}

func randomColor() string {
	color := "#"
	for i := 0; i < 6; i++ {
		color += randomHex()
	}
	return color
}

func listView(data hypp.MemoData) *hypp.VNode {
	list := data.(List)
	return html.P(hypp.HProps{
		"style": map[string]string{
			"backgroundColor": randomColor(),
			"color":           randomColor(),
		},
	}, hypp.Text(strings.Join(list, ", ")))
}

func moreItems(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.list = append(newState.list, randomHex())
	return newState
}

func increment(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.counter++
	return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
	hypp.App(hypp.AppProps[*MyState]{
		Driver: driver,
		Init: &MyState{
			list: []string{"a", "b", "c"},
		},
		View: func(state *MyState) *hypp.VNode {
			return html.Main(
				nil,
				html.Button(
					hypp.HProps{"onclick": hypp.Action[*MyState](moreItems)},
					hypp.Text("Grow list"),
				),
				html.Button(
					hypp.HProps{"onclick": hypp.Action[*MyState](increment)},
					hypp.Text("+1 to counter"),
				),
				html.P(nil, hypp.Textf("Counter: %d", state.counter)),
				html.P(nil, hypp.Textf("Regular view showing list:")),
				listView(state.list),
				html.P(nil, hypp.Text("Memoized view showing list:")),
				hypp.Memo(listView, state.list),
			)
		},
		Node: node,
	})
}

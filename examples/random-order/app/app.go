package app

import (
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type State struct {
	hypp.EmptyState
	lists [][3]int
	index int
}

func (s State) clone() *State {
	return &s
}

func incrementIndex(state *State, payload hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.index = (newState.index + 1) % len(newState.lists)
	return newState
}

func button() *hypp.VNode {
	return html.Button(
		hypp.HProps{
			"onclick": hypp.Action[*State](incrementIndex),
		},
		hypp.Text("Next"),
	)
}

func list(list [3]int) *hypp.VNode {
	children := make([]*hypp.VNode, len(list))
	for i, item := range list {
		children[i] = html.Div(
			hypp.HProps{
				"class": []string{"item", fmt.Sprintf("position-%d", i)},
				"key":   item,
			},
			hypp.Textf("%d", item),
		)
	}
	return html.Div(
		hypp.HProps{
			"class": "container",
		},
		children...,
	)
}

func Run(node hypp.Element) {
	hypp.App(hypp.AppProps[*State]{
		Init: &State{
			lists: [][3]int{
				{1, 2, 3},
				{2, 3, 1},
			},
			index: 0,
		},
		View: func(state *State) *hypp.VNode {
			return html.Main(
				nil,
				button(),
				list(state.lists[state.index]),
			)
		},
		Node: node,
	})
}

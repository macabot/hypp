// Package app implements a counter.
// Click the buttons to increase or decrease the count.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/MrBgMy

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type State struct {
	hypp.EmptyState
	count int
}

func (m State) clone() *State {
	return &m
}

func subtract(state *State, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.count--
	return newState
}

func add(state *State, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.count++
	return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
	hypp.App(hypp.AppProps[*State]{
		Driver: driver,
		Init:   &State{},
		View: func(state *State) *hypp.VNode {
			return html.Main(
				nil,
				html.H1(nil, hypp.Textf("%d", state.count)),
				html.Button(
					hypp.HProps{"onclick": hypp.Action[*State](subtract)},
					hypp.Text("ー"),
				),
				html.Button(
					hypp.HProps{"onclick": hypp.Action[*State](add)},
					hypp.Text("＋"),
				),
			)
		},
		Node: node,
	})
}

// Package app implements a counter.
// Click the buttons to increase or decrease the count.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/MrBgMy

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type MyState struct {
	hypp.EmptyState
	count int
}

func (m MyState) clone() *MyState {
	return &m
}

func subtract(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.count--
	return newState
}

func add(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.count++
	return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
	hypp.App(hypp.AppProps[*MyState]{
		Driver: driver,
		Init:   &MyState{},
		View: func(state *MyState) *hypp.VNode {
			return html.Main(
				nil,
				html.H1(nil, hypp.Textf("%d", state.count)),
				html.Button(
					hypp.HProps{"onclick": hypp.Action[*MyState](subtract)},
					hypp.Text("ー"),
				),
				html.Button(
					hypp.HProps{"onclick": hypp.Action[*MyState](add)},
					hypp.Text("＋"),
				),
			)
		},
		Node: node,
	})
}

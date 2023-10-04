// Package app renders a page that says "👋 Hi.".
package app

// This file is based on https://codepen.io/jorgebucaran/pen/Qdwpxy

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/window"
)

type State struct {
	hypp.EmptyState
	message string
}

func Run(node window.Element) {
	hypp.App(hypp.AppProps[*State]{
		Init: &State{message: "👋 Hi."},
		View: func(state *State) *hypp.VNode {
			return hypp.H("h1", nil, hypp.Text(state.message))
		},
		Node: node,
	})
}

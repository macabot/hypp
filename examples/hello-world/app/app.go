// This file is based on https://codepen.io/jorgebucaran/pen/Qdwpxy
package app

import (
    "github.com/macabot/hypp"
)

type MyState struct {
    hypp.EmptyState
    message string
}

func Run(driver hypp.Driver, node hypp.Node) {
    hypp.App(hypp.AppProps[*MyState]{
        Driver: driver,
        Init: &MyState{message: "👋 Hi."},
        View: func(state *MyState) *hypp.VNode {
            return hypp.H("h1", nil, hypp.Text(state.message))
        },
        Node: node,
    })
}

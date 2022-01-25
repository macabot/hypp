// Package app synchronizes the h1 title with the value of the input field.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/qRMEGX

import (
	"strings"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type State struct {
	hypp.EmptyState
	message string
}

func (m State) clone() *State {
	return &m
}

func withPayload[S hypp.State](filter func(payload hypp.Payload) hypp.Dispatchable) hypp.Action[S] {
	return func(_ S, payload hypp.Payload) hypp.Dispatchable {
		return filter(payload)
	}
}

func input[S hypp.State](oninput hypp.Action[S], props hypp.HProps) *hypp.VNode {
	props.Set("oninput", oninput)
	return html.Input(props)
}

func title(text string) *hypp.VNode {
	return html.H1(nil, hypp.Text(text))
}

func setText(state *State, payload hypp.Payload) hypp.Dispatchable {
	message := payload.(string)
	newState := state.clone()
	newState.message = message
	return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
	hypp.App(hypp.AppProps[*State]{
		Driver: driver,
		Init:   &State{},
		View: func(state *State) *hypp.VNode {
			t := state.message
			if strings.TrimSpace(state.message) == "" {
				t = "🤷"
			}
			return html.Main(
				nil,
				title(t),
				input(
					withPayload[*State](func(payload hypp.Payload) hypp.Dispatchable {
						event := payload.(hypp.Event)
						return hypp.ActionAndPayload[*State]{
							Action:  setText,
							Payload: event.Target().Value(),
						}
					}),
					hypp.HProps{
						"placeholder": "Type in something...",
						"value":       state.message,
						"type":        "text",
					},
				),
			)
		},
		Node: node,
	})
}

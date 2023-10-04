// Package app lets you write 140 characters in a textarea.
// Click "Tweet" to empty the textarea.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/bgWBdV

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"github.com/macabot/hypp/window"
)

type State struct {
	hypp.EmptyState
	content string
	count   int
}

func (m State) clone() *State {
	return &m
}

func button(onclick hypp.Dispatchable, text string, props hypp.HProps) *hypp.VNode {
	props.Set("onclick", onclick)
	return html.Button(props, hypp.Text(text))
}

func title(text string) *hypp.VNode {
	return html.H1(nil, hypp.Text(text))
}

var maxLength = 140

type contentAndLength struct {
	content string
	length  int
}

func setText(state *State, payload hypp.Payload) hypp.Dispatchable {
	cl := payload.(contentAndLength)
	newState := state.clone()
	if cl.length > maxLength {
		newState.count = 0
	} else {
		newState.content = cl.content
		newState.count = state.count + len(state.content) - len(cl.content)
	}
	return newState
}

func Run(node window.Element) {
	hypp.App(hypp.AppProps[*State]{
		Init: &State{count: maxLength},
		View: func(state *State) *hypp.VNode {
			var oninput hypp.Action[*State] = func(_ *State, payload hypp.Payload) hypp.Dispatchable {
				event := payload.(window.Event)
				content := event.Target().Value()
				return hypp.ActionAndPayload[*State]{
					Action: setText,
					Payload: contentAndLength{
						content: content,
						length:  len(content),
					},
				}
			}
			return html.Main(
				nil,
				title("Tweeter 🦤"),
				html.Textarea(hypp.HProps{
					"placeholder": "What's on your mind?",
					"oninput":     oninput,
					"value":       state.content,
					"rows":        4 + len(state.content)/100,
				}),
				html.Section(
					nil,
					hypp.Textf("%d", state.count),
					button(&State{}, "Tweet", hypp.HProps{
						"disabled": state.count >= maxLength,
					}),
				),
			)
		},
		Node: node,
	})
}

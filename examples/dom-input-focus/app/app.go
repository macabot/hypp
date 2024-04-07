// Package app shows how to focus an input element.
// Click on "Add new" to add a new input element and focus on it.
// Click one of the "Focus" buttons to focus on the input element next to it.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/wjvEBj

import (
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"github.com/macabot/hypp/window"
)

type State struct {
	hypp.EmptyState
	count int
}

func (m State) clone() *State {
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
	id            string
	preventScroll bool
}

func justFocus(_ hypp.Dispatch, payload hypp.Payload) {
	props := payload.(focusProps)
	window.RequestAnimationFrame(func() {
		window.Document().
			GetElementById(props.id).
			Call("focus", map[string]any{
				"preventScroll": props.preventScroll,
			})
	})
}

func focus(id string, preventScroll bool) hypp.Effect {
	return hypp.Effect{
		Effecter: justFocus,
		Payload:  focusProps{id: id, preventScroll: preventScroll},
	}
}

func domID(n int) string {
	return fmt.Sprintf("input-%d", n)
}

func addAndFocus(state *State, _ hypp.Payload) hypp.Dispatchable {
	id := domID(state.count)
	newState := state.clone()
	newState.count++
	return hypp.StateAndEffects[*State]{
		State: newState,
		Effects: []hypp.Effect{
			focus(id, false),
		},
	}
}

func focusStateEffect(state *State, payload hypp.Payload) hypp.Dispatchable {
	return hypp.StateAndEffects[*State]{
		State: state,
		Effects: []hypp.Effect{
			focus(payload.(string), false),
		},
	}
}

func Run(node window.Element) {
	hypp.App(hypp.AppProps[*State]{
		Init: &State{count: 1},
		View: func(state *State) *hypp.VNode {
			children := make([]*hypp.VNode, state.count)
			for i := range children {
				children[i] = container(
					html.Input(hypp.HProps{
						"type": "text",
						"id":   domID(i),
					}),
					button(hypp.ActionAndPayload[*State]{
						Action:  focusStateEffect,
						Payload: domID(i),
					}, "Focus"),
				)
			}
			children = append(
				children,
				separator(),
				button(hypp.Action[*State](addAndFocus), "Add new"),
			)
			return main(children...)
		},
		Node: node,
	})
}

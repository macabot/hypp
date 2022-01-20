// This file is based on https://codepen.io/jorgebucaran/pen/wjvEBj
package app

import (
	"fmt"

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

var window hypp.Window

func justFocus(_ hypp.Dispatch, payload hypp.Payload) {
	props := payload.(focusProps)
	window.RequestAnimationFrame(func() {
		window.EscapeToValue().
			Get("document").
			Call("getElementById", props.id).
			Call("focus", map[string]interface{}{
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

func addAndFocus(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	id := domID(state.count)
	newState := state.clone()
	newState.count++
	return hypp.StateAndEffects[*MyState]{
		State: newState,
		Effects: []hypp.Effect{
			focus(id, false),
		},
	}
}

func focusStateEffect(state *MyState, payload hypp.Payload) hypp.Dispatchable {
	return hypp.StateAndEffects[*MyState]{
		State: state,
		Effects: []hypp.Effect{
			focus(payload.(string), false),
		},
	}
}

func Run(driver hypp.Driver, node hypp.Node) {
	window = driver.Window()
	hypp.App(hypp.AppProps[*MyState]{
		Driver: driver,
		Init:   &MyState{count: 1},
		View: func(state *MyState) *hypp.VNode {
			children := make([]*hypp.VNode, state.count)
			for i := range children {
				children[i] = container(
					html.Input(hypp.HProps{
						"type": "text",
						"id":   domID(i),
					}),
					button(hypp.ActionAndPayload[*MyState]{
						Action:  focusStateEffect,
						Payload: domID(i),
					}, "Focus"),
				)
			}
			children = append(
				children,
				separator(),
				button(hypp.Action[*MyState](addAndFocus), "Add new"),
			)
			return main(children...)
		},
		Node: node,
	})
}

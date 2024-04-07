// Package app tracks the position of your mouse.
// Click anywhere to stop/start tracking your mouse.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/Bpyraw

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"github.com/macabot/hypp/window"
)

type State struct {
	hypp.EmptyState
	x          int
	y          int
	isTracking bool
}

func (m State) clone() *State {
	return &m
}

type mouseProps struct {
	name         string
	dispatchable hypp.Dispatchable
}

func on(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(mouseProps)
	listener := func(event window.Event) {
		dispatch(props.dispatchable, event)
	}
	id := window.AddEventListener(props.name, listener)
	return func() {
		window.RemoveEventListener(props.name, id)
	}
}

func onMouseMove(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: on,
		Payload:    mouseProps{name: "mousemove", dispatchable: dispatchable},
	}
}

func onClick(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: on,
		Payload:    mouseProps{name: "click", dispatchable: dispatchable},
	}
}

func draggable(content string, props hypp.HProps) *hypp.VNode {
	props.Set("class", "draggable")
	return html.Span(props, hypp.Text(content))
}

func title(text string) *hypp.VNode {
	return html.H1(nil, hypp.Text(text))
}

func titlef(format string, args ...any) *hypp.VNode {
	return html.H1(nil, hypp.Textf(format, args...))
}

func strong(text string) *hypp.VNode {
	return html.Strong(nil, hypp.Text(text))
}

func move(state *State, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(window.Event).Value
	newState := state.clone()
	newState.x = event.Get("x").Int()
	newState.y = event.Get("y").Int()
	return newState
}

func toggle(state *State, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.isTracking = !state.isTracking
	return newState
}

func Run(node window.Element) {
	hypp.App(hypp.AppProps[*State]{
		Init: &State{x: -1, y: -1, isTracking: true},
		View: func(state *State) *hypp.VNode {
			var t string
			if state.isTracking {
				t = " stop "
			} else {
				t = " start "
			}
			var v *hypp.VNode
			if state.x-state.y == 0 {
				v = title("Move your 🐭")
			} else {
				v = titlef("%d, %d", state.x, state.y)
			}
			return html.Main(
				nil,
				hypp.Text("Click anywhere to"),
				strong(t),
				hypp.Text("tracking the mouse."),
				v,
			)
		},
		Subscriptions: func(state *State) []hypp.Subscription {
			move := onMouseMove(hypp.Action[*State](move))
			move.Disabled = !state.isTracking
			return []hypp.Subscription{
				move,
				onClick(hypp.Action[*State](toggle)),
			}
		},
		Node: node,
	})
}

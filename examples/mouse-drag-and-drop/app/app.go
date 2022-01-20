// This file is based on https://codepen.io/jorgebucaran/pen/apzYvo
package app

import (
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type MyState struct {
	hypp.EmptyState
	dragging bool
	offsetX  int
	offsetY  int
	x        int
	y        int
}

func (m MyState) clone() *MyState {
	return &m
}

var window hypp.Window

type mouseProps struct {
	name         string
	dispatchable hypp.Dispatchable
}

func on(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(mouseProps)
	listener := func(event hypp.Event) {
		dispatch(props.dispatchable, event)
	}
	id := window.AddEventListener(props.name, listener)
	return func() {
		window.RemoveEventListener(props.name, id)
	}
}

func onMouseUp(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: on,
		Payload:    mouseProps{name: "mouseup", dispatchable: dispatchable},
	}
}

func onMouseMove(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: on,
		Payload:    mouseProps{name: "mousemove", dispatchable: dispatchable},
	}
}

func draggable(content string, props hypp.HProps) *hypp.VNode {
	props.Set("class", "draggable")
	return html.Span(props, hypp.Text(content))
}

func title(text string) *hypp.VNode {
	return html.H1(nil, hypp.Text(text))
}

func drop(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.dragging = false
	return newState
}

func drag(state *MyState, payload hypp.Payload) hypp.Dispatchable {
	props := payload.(hypp.Event).EscapeToValue()
	newState := state.clone()
	newState.dragging = true
	newState.offsetX = props.Get("offsetX").Int()
	newState.offsetY = props.Get("offsetY").Int()
	newState.x = props.Get("pageX").Int()
	newState.y = props.Get("pageY").Int()
	return newState
}

func move(state *MyState, payload hypp.Payload) hypp.Dispatchable {
	if !state.dragging {
		return state
	}
	event := payload.(hypp.Event).EscapeToValue()
	newState := state.clone()
	newState.x = event.Get("pageX").Int()
	newState.y = event.Get("pageY").Int()
	return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
	window = driver.Window()
	hypp.App(hypp.AppProps[*MyState]{
		Driver: driver,
		Init:   &MyState{x: 5, y: 20},
		View: func(state *MyState) *hypp.VNode {
			draggableContent := "🛸"
			titleText := "Drag the UFO!"
			if state.dragging {
				draggableContent = "👽"
				titleText = "Good job!"
			}
			return html.Main(
				nil,
				draggable(draggableContent, hypp.HProps{
					"onmousedown": hypp.Action[*MyState](drag),
					"style": map[string]string{
						"cursor":     "move",
						"left":       fmt.Sprintf("%dpx", state.x-state.offsetX),
						"top":        fmt.Sprintf("%dpx", state.y-state.offsetY),
						"position":   "absolute",
						"userSelect": "none",
					},
				}),
				title(titleText),
			)
		},
		Subscriptions: func(_ *MyState) []hypp.Subscription {
			return []hypp.Subscription{
				onMouseUp(hypp.Action[*MyState](drop)),
				onMouseMove(hypp.Action[*MyState](move)),
			}
		},
		Node: node,
	})
}

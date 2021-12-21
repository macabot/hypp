// This file is based on https://codepen.io/jorgebucaran/pen/zNxRLy
package app

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type MyState struct {
	hypp.EmptyState
	todos []TodoItem
	value string
}

type TodoItem struct {
	isEditing bool
	lastValue string
	value     string
}

func preventDefault(action hypp.Action) hypp.Action {
	return func(state MyState, payload hypp.Payload) hypp.Dispatchable {
		event := payload.(hypp.Event)
		event.PreventDefault()
		return action(state, payload)
	}
}

func withPayload(filter func(e hypp.Event) hypp.ActionAndPayload) hypp.Action {
	return func(_ MyState, payload hypp.Payload) hypp.ActionAndPayload {
		return filter(payload)
	}
}

func targetValue(action hypp.Action) hypp.ActionAndPayload {
	return withPayload(func(e hypp.Event) hypp.ActionAndPayload {
		return hypp.ActionAndPayload{
			Action:  action,
			Payload: e.Target().Value(),
		}
	})
}

func form(onsubmit hypp.ActionLike, props hypp.HProps, children []hypp.VNode) hypp.VNode {
	props["onsubmit"] = preventDefault(onsubmit)
	return html.Form(
		props,
		children...,
	)
}

func checkbox() hypp.VNode {
	return html.Input(
		hypp.HProps{
			"type": "checkbox",
		},
	)
}

func submit(value string) hypp.VNode {
	return html.Input(
		hypp.HProps{
			"type":  "submit",
			"value": value,
		},
	)
}

func input(oninput hypp.Action, props hypp.HProps) hypp.VNode {
	props["type"] = "text"
	props["oninput"] = targetValue(oninput)
	return html.Input(props)
}

func todosView(value string) VNode {
	return html.Label(
		nil,
		checkbox(),
		html.Span(
			nil,
			hypp.Text(value),
		),
	)
}

func newValue(state MyState, value hypp.Payload) hypp.Dispatchable {
	state.value = value
	return value
}

func newTodo(state MyState, _ hypp.Payload) hypp.Dispatchable {
	if len(state.todos) == 0 {
		return state
	}
	state.todos = append(state.todos, TodoItem{value: state.value})
	return state
}

func Run() {
	hypp.App[MyState](hypp.AppProps{
		Init: MyState{
			todos: []TodoItem{{value: "Learn Hypp"}},
			value: "",
		},
		View: func(state MyState) VNode {
			children := todosView("To-do list ✏️", state.todos)
			children = append(children, form(
				newTodo,
				nil,
				input(
					newValue,
					hypp.HProps{
						"value": state.value,
					},
				),
				submit("Add new"),
			))
			return html.Main(
				nil,
				children...,
			)
		},
		Node: js.Global().Get("document").Call("getElementById", "app"),
	})
}

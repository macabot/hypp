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

func (m MyState) clone() *MyState {
	return &m
}

type TodoItem struct {
	isEditing bool
	lastValue string
	value     string
}

func preventDefault[S hypp.State](action hypp.Action[S]) hypp.Action[S] {
	return func(state S, payload hypp.Payload) hypp.Dispatchable {
		event := payload.(hypp.Event)
		event.PreventDefault()
		return action(state, payload)
	}
}

func withPayload[S hypp.State](filter func(e hypp.Event) hypp.ActionAndPayload[S]) hypp.Action[S] {
	return func(_ S, payload hypp.Payload) hypp.Dispatchable {
		return filter(payload.(hypp.Event))
	}
}

func targetValue[S hypp.State](action hypp.Action[S]) hypp.Action[S] {
	return withPayload(func(e hypp.Event) hypp.ActionAndPayload[S] {
		return hypp.ActionAndPayload[S]{
			Action:  action,
			Payload: e.Target().Value(),
		}
	})
}

func form[S hypp.State](onsubmit hypp.Action[S], props hypp.HProps, children ...hypp.VNode) hypp.VNode {
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

func input[S hypp.State](oninput hypp.Action[S], props hypp.HProps) hypp.VNode {
	props["type"] = "text"
	props["oninput"] = targetValue(oninput)
	return html.Input(props)
}

func todoItemView(value string) hypp.VNode {
	return html.Label(
		nil,
		checkbox(),
		html.Span(
			nil,
			hypp.Text(value),
		),
	)
}

func todosView(title string, todos []TodoItem) []hypp.VNode {
	ulChildren := make([]hypp.VNode, len(todos))
	for i, todo := range todos {
		ulChildren[i] = html.Li(nil, todoItemView(todo.value))
	}
	return []hypp.VNode{
		html.H1(nil, hypp.Text(title)),
		html.Ul(nil, ulChildren...),
	}
}

func newValue(state *MyState, value hypp.Payload) hypp.Dispatchable {
	state = state.clone()
	state.value = value.(string)
	return state
}

func newTodo(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	if len(state.todos) == 0 {
		return state
	}
	state = state.clone()
	state.todos = append(state.todos, TodoItem{value: state.value})
	return state
}

func Run(node hypp.Node) {
	hypp.App[*MyState](hypp.AppProps[*MyState]{
		Init: MyState{
			todos: []TodoItem{{value: "Learn Hypp"}},
			value: "",
		},
		View: func(state *MyState) hypp.VNode {
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
		Node: node,
	})
}

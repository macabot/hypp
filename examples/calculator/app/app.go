// This file is based on https://codepen.io/jorgebucaran/pen/PmjRov
package app

import (
    "github.com/macabot/hypp"
    "github.com/macabot/hypp/tag/html"
)

type MyState struct {
    hypp.EmptyState
    fn string
    carry float64
    value float64
    hasCarry bool
}

func (m MyState) clone() *MyState {
    return &m
}

var computer = map[string]func(a, b float64) float64 {
    "+": func(a, b float64) float64 { return a + b },
    "-": func(a, b float64) float64 { return a - b },
    "×": func(a, b float64) float64 { return a * b },
    "÷": func(a, b float64) float64 { return a / b },
}

var initialState = &MyState{}

func clear() *MyState {
    return initialState
}

func newDigit(state *MyState, payload hypp.Payload) hypp.Dispatchable {
    number := payload.(float64)
    state.clone()
    state.hasCarry = false
    v := 0.0
    if !state.hasCarry {
        v = state.value
    }
    state.value = v * 10 + number
    return state
}

func newFn(state *MyState, payload hypp.Payload) hypp.Dispatchable {
    fn := payload.(string)
    newState := state.clone()
    newState.fn = fn
    newState.hasCarry = true
    newState.carry = state.value
    if state.hasCarry || state.fn == "" {
        newState.value = state.value
    } else {
        newState.value = computer[state.fn](state.carry, state.value)
    }
    return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
    hypp.App[*MyState](hypp.AppProps[*MyState]{
        Driver: driver,
        Init: &MyState{}, // TODO
        View: func(state *MyState) *hypp.VNode {
            // TODO
        },
        Node: node,
    })
}

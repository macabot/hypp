// This file is based on https://codepen.io/jorgebucaran/pen/PmjRov
package app

import (
    "strconv"

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
var computerKeys = []string{"+", "-", "×", "÷"}

func clear(_ *MyState, payload hypp.Payload) hypp.Dispatchable {
    return &MyState{}
}

func newDigit(state *MyState, payload hypp.Payload) hypp.Dispatchable {
    number := payload.(float64)
    newState := state.clone()
    newState.hasCarry = false
    v := 0.0
    if !state.hasCarry {
        v = state.value
    }
    newState.value = v * 10 + number
    return newState
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

func equal(state *MyState, _ hypp.Payload) hypp.Dispatchable {
    newState := state.clone()
    newState.hasCarry = true
    if !state.hasCarry {
        newState.carry = state.value
    }
    if state.fn != "" {
        a := state.carry
        b := state.value
        if state.hasCarry {
            a = state.value
            b = state.carry
        }
        newState.value = computer[state.fn](a, b)
    }
    return newState
}

func displayView(value float64) *hypp.VNode {
    return html.Div(hypp.HProps{"class": "display"}, hypp.Text(strconv.FormatFloat(value, 'f', -1, 64)))
}

func keysView(children ...*hypp.VNode) *hypp.VNode {
    return html.Div(hypp.HProps{"class": "keys"}, children...)
}

func fnView(keys []string) []*hypp.VNode {
    out := make([]*hypp.VNode, len(keys))
    for i, fn := range keys {
        out[i] = html.Button(hypp.HProps{
            "class": "function",
            "onclick": hypp.ActionAndPayload[*MyState]{Action: newFn, Payload: fn},
        }, hypp.Text(fn))
    }
    return out
}

func digitsView(digits []float64) []*hypp.VNode {
    out := make([]*hypp.VNode, len(digits))
    for i, digit := range digits {
        out[i] = html.Button(hypp.HProps{
            "class": map[string]bool{
                "zero": digit == 0,
            },
            "onclick": hypp.ActionAndPayload[*MyState]{Action: newDigit, Payload: digit},
        }, hypp.Textf("%.0f", digit))
    }
    return out
}

func acView() *hypp.VNode {
    return html.Button(hypp.HProps{"onclick": hypp.Action[*MyState](clear)}, hypp.Text("AC"))
}

func eqView() *hypp.VNode {
    return html.Button(hypp.HProps{"onclick": hypp.Action[*MyState](equal), "class": "equal"}, hypp.Text("="))
}

func Run(driver hypp.Driver, node hypp.Node) {
    hypp.App(hypp.AppProps[*MyState]{
        Driver: driver,
        Init: &MyState{},
        View: func(state *MyState) *hypp.VNode {
            keys := fnView(computerKeys)
            keys = append(keys, digitsView([]float64{7, 8, 9, 4, 5, 6, 1, 2, 3, 0})...)
            keys = append(keys, acView(), eqView())
            return html.Main(
                nil,
                displayView(state.value),
                keysView(keys...),
            )
        },
        Node: node,
    })
}

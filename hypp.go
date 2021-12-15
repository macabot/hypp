// +build go1.18
// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.d.ts
package hypp

import (
    "fmt"
)

func App[S any](props AppProps[S]) Dispatch {
    return func(dispatchable Dispatchable, payload Payload) {}
}

type HProps map[string]interface{}

func (h HProps) key() Option[string] {
    if h == nil {
        return Option[string]{}
    }
    if key, ok := h["key"]; ok {
        return Option[string]{V: fmt.Sprintf("%v", key), OK: ok}
    }
    return Option[string]{}
}

func H(tag string, props HProps, children...VNode) VNode {
    return VNode{}
}

func Memo[D any](view func(data D) VNode, data D) VNode {
    return VNode{}
}

func Text(value string) VNode {
    return VNode{}
}

type Payload interface{}

type Action[S any] func(state S, payload Payload) Dispatchable

func (_ Action[S]) IAmDispatchable() {}

type Node struct{
    parentNode Option[Node]
}

type Subscriptions[S any] func(state S) []Subscription[S]

type Render func()

type AppProps[S any] struct {
    Init Dispatchable
    Subscriptions Subscriptions[S]
    DispatchInitializer func(dispatch Dispatch) Dispatch
    View func(state S) VNode
    Node Option[Node]

    vdom VNode
    dispatch Dispatch
    subs []Subscription[S]
    render Render
    busy bool
    state S
}

type Dispatch func(dispatchable Dispatchable, payload Payload)

type Dispatchable interface {
    IAmDispatchable()
}

type StateAndEffects[S any] struct {
    State S
    Effect []Effect[S]
}

func (_ StateAndEffects[S]) IAmDispatchable() {}

type ActionAndPayload[S any] struct {
    Action Action[S]
    Payload Payload
}

func (_ ActionAndPayload[S]) IAmDispatchable() {}

type Effect[S any] struct {
    Effecter func(dispatch Dispatch, payload Payload)
    Payload Payload
}

type Subscription[S any] struct {
    Subscriber func(dispatch Dispatch, payload Payload) Unsubscribe
    Payload Payload
}

type Unsubscribe func()

type Option[T any] struct {
    V T
    OK bool
}

type VNode struct {
    props HProps
    children []VNode
    node Option[Node]
    events map[string]interface{} // Action[S] | ActionAndPayload[S, P]
    key Option[string]
    tag string
    memoView func(data interface{}) VNode
    memo interface{} // Indexable
    kind int
}

// +build go1.18
// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.d.ts
package hypp

func App[S any](props AppProps[S]) Dispatch {
    return func(dispatchable Dispatchable, payload Payload) {}
}

type HProps map[string]interface{}

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

type Node struct{}

type AppProps[S any] struct {
    Init Dispatchable
    Subscriptions func(state S) []Subscription[S]
    Dispatch func(dispatch Dispatch) Dispatch
    view func(state S) VNode
    node Node
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

type VNode struct {
    props HProps
    children []VNode
    node Node
    events map[string]interface{} // Action[S] | ActionAndPayload[S, P]
    key *string
    tag interface{} // string | func(data Indexable) VNode
    memo interface{} // Indexable
    kind int
}

// +build go1.18
// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.d.ts
package hypp

import (
    "fmt"
)

type State interface {
    comparable
    Dispatchable
}

func App[S State](props AppProps[S]) Dispatch {
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

func (h HProps) get(key string) interface{} {
    if h == nil {
        return nil
    }
    return h[key]
}

func (h HProps) has(key string) bool {
    if h == nil {
        return false
    }
    _, ok := h[key]
    return ok
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

type actionLike interface {
    Dispatchable
    iAmActionLike()
}

type Action[S State] func(state S, payload Payload) Dispatchable

func (_ Action[S]) IAmDispatchable() {}

func (_ Action[S]) iAmActionLike() {}

type EventListener func(VNode, Event)

type Node interface {
    parentNode() Option[Node]
    nodeType() int
    nodeValue() string
    setNodeValue(nodeValue string)
    nodeName() string
    childNodes() []Node
    insertBefore(newNode, referenceNode Node) Node
    removeChild(child Node)
    get(name string) interface{}
    has(name string) bool
    set(name string, value interface{})
    appendChild(child Node) Node
    removeEventListener(kind string, listener EventListener)
    addEventListener(kind string, listener EventListener)
    removeAttribute(name string)
    setAttribute(name string, value interface{})
    setEvent(name string, value interface{})
    getEvent(name string) interface{}
}

type Subscriptions[S State] func(state S) []Subscription[S]

type Render func()

type AppProps[S State] struct {
    Init Dispatchable
    Subscriptions Subscriptions[S]
    DispatchInitializer func(dispatch Dispatch) Dispatch
    View func(state S) VNode
    Node Node

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

type StateAndEffects[S State] struct {
    State S
    Effects []Effect[S]
}

func (_ StateAndEffects[S]) IAmDispatchable() {}

type ActionAndPayload[S State] struct {
    Action Action[S]
    Payload Payload
}

func (_ ActionAndPayload[S]) IAmDispatchable() {}

func (_ ActionAndPayload[S]) iAmActionLike() {}

type Effect[S State] struct {
    Effecter func(dispatch Dispatch, payload Payload)
    Payload Payload
}

type Subscription[S State] struct {
    Subscriber func(dispatch Dispatch, payload Payload) Unsubscribe
    Payload Payload
    unsubscribe Unsubscribe
    Disabled bool
}

type Unsubscribe func()

type Option[T any] struct {
    V T
    OK bool
}

type VNode struct {
    props HProps
    children []VNode
    node Node // Can be nil
    events map[string]actionLike // Action[S] | ActionAndPayload[S]
    key Option[string]
    tag string
    memoView func(data interface{}) VNode
    memo interface{} // Indexable
    kind int
    isNil bool
}

var NilVNode = VNode{isNil: true}

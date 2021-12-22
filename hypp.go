//go:build go1.18
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
	return app[S](props)
}

type HProps map[string]interface{}

func (h HProps) key() Option[string] {
	if h == nil {
		return Option[string]{}
	}
	if key, ok := h["key"]; ok {
		return Option[string]{V: fmt.Sprint(key), OK: ok}
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

func H(tag string, props HProps, children ...VNode) VNode {
	return h(tag, props, children)
}

func Memo(view func(data interface{}) VNode, data interface{}) VNode {
	return memo(view, data)
}

func Text(value string) VNode {
	return text(value, nil)
}

type Payload interface{}

type ActionLike interface { // TODO should be exported?
	Dispatchable
	iAmActionLike()
}

type Action[S State] func(state S, payload Payload) Dispatchable

func (_ Action[S]) IAmDispatchable() {}

func (_ Action[S]) iAmActionLike() {}

type Event interface {
	Dispatchable
    Type() string
	PreventDefault()
	Target() EventTarget
}

type EventTarget interface {
	Value() interface{}
}

type EventListener func(Event)

type Node interface {
	ParentNode() Option[Node]
	NodeType() int
	NodeValue() string
	SetNodeValue(nodeValue string)
	NodeName() string
	ChildNodes() []Node
	InsertBefore(newNode, referenceNode Node) Node
	RemoveChild(child Node)
	Get(name string) Value
	Has(name string) bool
	Set(name string, value interface{})
	AppendChild(child Node) Node
	RemoveEventListener(kind string, listener EventListener)
	AddEventListener(kind string, listener EventListener)
	RemoveAttribute(name string)
	SetAttribute(name string, value interface{})
	Events() Events
	Style() Style
}

type Value interface {
	Int() int
	String() string
	Bool() bool
}

type Events interface {
	Set(name string, event Event)
	Get(name string) Event
}

type Style interface {
	SetProperty(propertyName, value string)
	Set(name, value string)
	Get(name string) string
}

type Subscriptions[S State] func(state S) []Subscription[S]

type Render func()

type AppProps[S State] struct {
	Init                Dispatchable
	Subscriptions       Subscriptions[S]
	DispatchInitializer func(dispatch Dispatch) Dispatch
	View                func(state S) VNode
	Node                Node

	vdom     VNode
	dispatch Dispatch
	subs     []Subscription[S]
	render   Render
	busy     bool
	state    S
}

type Dispatch func(dispatchable Dispatchable, payload Payload)

type Dispatchable interface {
	IAmDispatchable()
}

type StateAndEffects[S State] struct {
	State   S
	Effects []Effect[S]
}

func (_ StateAndEffects[S]) IAmDispatchable() {}

type ActionAndPayload[S State] struct {
	Action  Action[S]
	Payload Payload
}

func (_ ActionAndPayload[S]) IAmDispatchable() {}

func (_ ActionAndPayload[S]) iAmActionLike() {}

type Effect[S State] struct {
	Effecter func(dispatch Dispatch, payload Payload)
	Payload  Payload
}

type Subscription[S State] struct {
	Subscriber  func(dispatch Dispatch, payload Payload) Unsubscribe
	Payload     Payload
	unsubscribe Unsubscribe
	Disabled    bool
}

type Unsubscribe func()

type Option[T any] struct {
	V  T
	OK bool
}

type VNode struct {
	props    HProps
	children []VNode
	node     Node                  // Can be nil
	events   map[string]ActionLike // Action[S] | ActionAndPayload[S]
	key      Option[string]
	tag      string
	memoView func(data interface{}) VNode
	memo     interface{} // Indexable
	kind     int
	isNil    bool
}

var NilVNode = VNode{isNil: true}

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

func (h HProps) Key() Option[string] {
	if key := h.Get("key"); key.OK {
		return Option[string]{V: fmt.Sprint(key.V), OK: true}
	}
	return Option[string]{}
}

func (h HProps) Get(key string) Option[interface{}] {
	if h == nil {
		return Option[interface{}]{}
	}
	v, ok := h[key]
	return Option[interface{}]{V: v, OK: ok}
}

func (h HProps) Has(key string) bool {
	if h == nil {
		return false
	}
	_, ok := h[key]
	return ok
}

func (h *HProps) Set(key string, value interface{}) {
	if *h == nil {
		*h = HProps{}
	}
	m := *h
	m[key] = value
}

func H(tag string, props HProps, children ...*VNode) *VNode {
	return h(tag, props, children)
}

func Memo(view func(data interface{}) *VNode, data interface{}) *VNode {
	return memo(view, data)
}

func Text(value string) *VNode {
	return text(value, nil)
}

func Textf(format string, a ...interface{}) *VNode {
	return Text(fmt.Sprintf(format, a...))
}

type Payload interface{}

type ActionLike interface {
	Dispatchable
	iAmActionLike()
}

type Action[S State] func(state S, payload Payload) Dispatchable

func (_ Action[S]) IAmDispatchable() {}

func (_ Action[S]) iAmActionLike() {}

type Event interface {
	EscapeToValuer
    Type() string
	PreventDefault()
	Target() EventTargetValuer
}

type EventTargetValuer interface {
	Value() string
}

type Value interface {
	Bool() bool
	Call(m string, args ...interface{}) Value
	Delete(p string)
	Equal(w Value) bool
	Float() float64
	Get(p string) Value
	Index(i int) Value
	InstanceOf(t Value) bool
	Int() int
	Invoke(args ...interface{}) Value
	IsNaN() bool
	IsNull() bool
	IsUndefined() bool
	Length() int
	New(args ...interface{}) Value
	Set(p string, x interface{})
	SetIndex(i int, x interface{})
	String() string
	Truthy() bool
	Type() Type
}

type Type int

const (
	TypeUndefined Type = iota
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeSymbol
	TypeObject
	TypeFunction
)

func (t Type) String() string {
	switch t {
	case TypeUndefined:
		return "undefined"
	case TypeNull:
		return "null"
	case TypeBoolean:
		return "boolean"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeSymbol:
		return "symbol"
	case TypeObject:
		return "object"
	case TypeFunction:
		return "function"
	default:
		panic("bad type")
	}
}


type EventListener func(Event)

type EventListenerID interface {
	IAmAnEventListenerID()
}

type EventTarget interface {
	RemoveEventListener(kind string, listenerID EventListenerID)
	AddEventListener(kind string, listener EventListener) EventListenerID
}

type Node interface {
	EventTarget
	ParentNode() Node
	NodeType() int
	NodeValue() string
	SetNodeValue(nodeValue string)
	NodeName() string
	ChildNodes() []Node
	InsertBefore(newNode, referenceNode Node) Node
	RemoveChild(child Node)
	Get(name string) Option[interface{}]
	// In implements the in-operator
	// See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/in
	In(name string) bool
	Set(name string, value interface{})
	AppendChild(child Node) Node
	RemoveAttribute(name string)
	SetAttribute(name string, value interface{})
	Events() Events
	Style() Style
	EventListenerID(kind string) EventListenerID
	SetEventListenerID(kind string, eventListenerID EventListenerID)
}

type EscapeToValuer interface {
	EscapeToValue() Value
}

type Window interface {
	EscapeToValuer
	EventTarget
	RequestAnimationFrame(f func()) int
}

type Events interface {
	Set(name string, event Dispatchable)
	Get(name string) Dispatchable
	Del(name string)
}

type Style interface {
	SetProperty(propertyName, value string)
	Set(name, value string)
	Get(name string) string
}

type Subscriptions[S State] func(state S) []Subscription

type Render func()

type ElementCreationOptions struct {
	Is string
}

type Driver interface {
	CreateTextNode(data string) Node
	CreateElementNS(namespaceURI, qualifiedName string, options Option[ElementCreationOptions]) Node
	CreateElement(tagName string, options Option[ElementCreationOptions]) Node
	Enqueue(render func())
	Window() Window
}

type AppProps[S State] struct {
	Driver Driver
	Init                Dispatchable
	Subscriptions       Subscriptions[S]
	DispatchInitializer func(dispatch Dispatch) Dispatch
	View                func(state S) *VNode
	Node                Node

	vdom     *VNode
	dispatch Dispatch
	subs     []Subscription
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
	Effects []Effect
}

func (_ StateAndEffects[S]) IAmDispatchable() {}

type ActionAndPayload[S State] struct {
	Action  Action[S]
	Payload Payload
}

func (_ ActionAndPayload[S]) IAmDispatchable() {}

func (_ ActionAndPayload[S]) iAmActionLike() {}

type Effect struct {
	Effecter func(dispatch Dispatch, payload Payload)
	Payload  Payload
}

type Subscription struct {
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
	children vKids
	node     Node                  // Can be nil
	events   map[string]ActionLike // Action[S] | ActionAndPayload[S]
	key      Option[string]
	tag      string
	memoView func(data interface{}) *VNode
	memo     interface{} // Indexable
	kind     int
}

type vKids []*VNode

func (v vKids) getKey(i int) Option[string] {
	if i < len(v) {
		return v[i].key
	}
	return Option[string]{}
}

func (v vKids) get(i int) *VNode {
	if i < len(v) {
		return v[i]
	}
	return nil
}

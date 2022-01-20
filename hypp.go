//go:build go1.18

// Package hypp creates reactive web applications.
package hypp

// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.d.ts

import (
	"fmt"
)

// State constrains the state that is used in the hypp application.
// It must be comparable and Dispatchable.
//
// Most often you will embed the EmptyState:
//	package example
//
//	type MyState struct {
//		hypp.EmptyState
//	}
//
// Alternatively, you can explicitly make your state Dispatchable:
//	package example
//
//	type MyState string
//
//	func(_ MyState) IAmDispatchable() {}
type State interface {
	comparable
	Dispatchable
}

// EmptyState implements the State constraint.
// Embed the EmptyState in your state to implement the State constraint:
//	package example
//
//	type MyState struct {
//		hypp.EmptyState
//		Foo string
//		Bar int
//	}
type EmptyState struct{}

// IAmDispatchable makes the EmptyState Dispatchable.
func (_ EmptyState) IAmDispatchable() {}

// App creates a new application.
func App[S State](props AppProps[S]) Dispatch {
	return app[S](props)
}

// HProps are the properties to create a *VNode.
//
// The allowed value type depends on the key:
//	| Key               | Value type                                            |
//	| ----------------- | ----------------------------------------------------- |
//	| Starts with "on"  | Dispatchable                                          |
//	| "class"           | bool, int, float64, string, []string, map[string]bool |
//	| "style"           | bool, int, float64, string, map[string]string         |
//	| Other             | bool, int, float64, string                            |
type HProps map[string]interface{}

// Key returns the "key" property, if available.
// The value is always converted into a string.
func (h HProps) Key() Option[string] {
	if key := h.Get("key"); key.OK {
		return Option[string]{V: fmt.Sprint(key.V), OK: true}
	}
	return Option[string]{}
}

// Get returns the requested key, if available.
func (h HProps) Get(key string) Option[interface{}] {
	if h == nil {
		return Option[interface{}]{}
	}
	v, ok := h[key]
	return Option[interface{}]{V: v, OK: ok}
}

// Has returns true if the requested key is found.
func (h HProps) Has(key string) bool {
	if h == nil {
		return false
	}
	_, ok := h[key]
	return ok
}

// Set sets the given key value pair.
// It is safe to call this method on a nil value.
func (h *HProps) Set(key string, value interface{}) {
	if *h == nil {
		*h = HProps{}
	}
	m := *h
	m[key] = value
}

// H creates a new *VNode specified by tag.
//
// See the tag package for functions that create specific tags:
//	package main
//
//	import (
//		"github.com/macabot/hypp"
//		"github.com/macabot/hypp/tag/html"
//	)
//
//	func main() {
//		hypp.H("main", hypp.HProps{"class": "main"})
//		// Is equivalent to
//		html.Main(hypp.HProps{"class": "main"})
//	}
func H(tag string, props HProps, children ...*VNode) *VNode {
	return h(tag, props, children)
}

func Memo(view func(data interface{}) *VNode, data interface{}) *VNode {
	return memo(view, data)
}

// Text creates a text *VNode.
func Text(value string) *VNode {
	return text(value, nil)
}

// Textf creates a text *VNode by interpolating the format with the arguments.
func Textf(format string, a ...interface{}) *VNode {
	return Text(fmt.Sprintf(format, a...))
}

// Payload is the value that is dispatched.
type Payload interface{}

type Action[S State] func(state S, payload Payload) Dispatchable

// IAmDispatchable makes Action Dispatchable.
func (_ Action[S]) IAmDispatchable() {}

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
	Driver              Driver
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

func (a *AppProps[S]) init() {
	if a.DispatchInitializer == nil {
		a.DispatchInitializer = dispatchInitializerID
	}
	if a.Init == nil {
		a.Init = EmptyState{}
	}
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
	node     Node // Can be nil
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

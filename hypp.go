//go:build go1.18

// Package hypp creates reactive web applications.
package hypp

// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.d.ts

import (
	"errors"
	"fmt"

	"github.com/macabot/hypp/js"
	"github.com/macabot/hypp/window"
)

// State constrains the state that is used in the hypp application.
// It must be comparable and [Dispatchable].
//
// Most often you will embed the [EmptyState]:
//
//	package example
//
//	type State struct {
//		hypp.EmptyState
//	}
//
// Alternatively, you can explicitly make your state Dispatchable:
//
//	package example
//
//	type State string
//
//	func(_ State) IAmDispatchable() {}
type State interface {
	comparable
	Dispatchable
}

// EmptyState implements the [State] constraint.
// Embed the EmptyState in your state to implement the State constraint:
//
//	package example
//
//	type State struct {
//		hypp.EmptyState
//		Foo string
//		Bar int
//	}
type EmptyState struct{}

// IAmDispatchable makes the EmptyState [Dispatchable].
func (EmptyState) IAmDispatchable() {}

// App creates a new application.
//
// It panics if the [js.GetDriver] returns nil.
// It also panics if [AppProps.Validate] returns an error for the given props.
func App[S State](props AppProps[S]) Dispatch {
	return app(props)
}

// HProps are the properties to create a [VNode].
//
// The allowed value type depends on the key:
//
//	| Key               | Value type                                            |
//	| ----------------- | ----------------------------------------------------- |
//	| Starts with "on"  | Dispatchable                                          |
//	| "class"           | bool, int, float64, string, []string, map[string]bool |
//	| "style"           | map[string]string                                     |
//	| Other             | bool, int, float64, string                            |
type HProps map[string]any

// key returns the "key" property, if available.
// The value is always converted into a string.
func (h HProps) key() option[string] {
	if key := h.get("key"); key.OK {
		return option[string]{V: fmt.Sprint(key.V), OK: true}
	}
	return option[string]{}
}

// clone returns a shallow clone of the HProps.
func (h HProps) clone() HProps {
	clone := make(HProps, len(h))
	for k, v := range h {
		clone[k] = v
	}
	return clone
}

// get returns the requested key, if available.
func (h HProps) get(key string) option[any] {
	if h == nil {
		return option[any]{}
	}
	v, ok := h[key]
	return option[any]{V: v, OK: ok}
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
func (h *HProps) Set(key string, value any) {
	if *h == nil {
		*h = HProps{}
	}
	m := *h
	m[key] = value
}

// H creates a new [VNode] specified by tag.
//
// See the tag package for functions that create specific tags:
//
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

type MemoData interface {
	Hash() string
}

func Memo(view func(data MemoData) *VNode, data MemoData) *VNode {
	return memo(view, data)
}

// Text creates a text [VNode].
func Text(value string) *VNode {
	return text(value, window.Element{})
}

// Textf creates a text [VNode] by interpolating the format with the arguments.
func Textf(format string, a ...any) *VNode {
	return Text(fmt.Sprintf(format, a...))
}

// Payload is the value that is dispatched.
type Payload any

// Action is a function which describes a transition between the current state and the next state.
// It must not perform any side-effects, but it may return side-effects using [StateAndEffects].
//
// An action is dispatched by either a DOM event, the effecter of an [Effect], or the subscriber of a [Subscription].
// When dispatched, an action always receives the current [State] as its first argument and an optional [Payload] as its second argument.
// An action that is dispatched by a DOM event will receive a [window.Event] as payload.
// An action that is dispatched by an [ActionAndPayload] will receive the 'Payload' field as payload.
type Action[S State] func(state S, payload Payload) Dispatchable

// IAmDispatchable makes Action [Dispatchable].
func (_ Action[S]) IAmDispatchable() {}

type Subscriptions[S State] func(state S) []Subscription

// AppProps is passed as an argument when creating an [App].
type AppProps[S State] struct {
	// Init is the dispatchable that initializes the app.
	// If Init is nil, it is replaced with the EmptyState.
	Init Dispatchable
	// Optional slice of subscriptions.
	Subscriptions Subscriptions[S]
	// Optional function that wraps the Dispatch function.
	DispatchWrapper func(dispatch Dispatch) Dispatch
	// View renders the app given the state.
	// It cannot be nil.
	View func(state S) *VNode
	// Node must have a parentNode that is not null.
	Node window.Element

	vdom     *VNode
	dispatch Dispatch
	subs     []Subscription
	render   func()
	busy     bool
	state    S
}

// Validate returns an error if one of the following is true:
//   - View is nil.
//   - Node is falsy.
//   - Node has a null parentNode.
func (a AppProps[S]) Validate() error {
	if a.View == nil {
		return errors.New("hypp: AppProps.View cannot be nil")
	} else if !a.Node.Truthy() {
		return errors.New("hypp: AppProps.Node cannot be falsy")
	} else if a.Node.ParentNode().IsNull() {
		return errors.New("hypp: AppProps.Node must have a parent node")
	}
	return nil
}

func validateDriver() error {
	if js.GetDriver() == nil {
		return errors.New("hypp: Driver in hypp/js cannot be nil")
	}
	return nil
}

func (a *AppProps[S]) init() {
	if a.DispatchWrapper == nil {
		a.DispatchWrapper = dispatchWrapperID
	}
	if a.Init == nil {
		a.Init = EmptyState{}
	}
	if err := validateDriver(); err != nil {
		panic(err)
	}
	if err := a.Validate(); err != nil {
		panic(err)
	}
}

type Dispatch func(dispatchable Dispatchable, payload Payload)

// Dispatchable is implemented by types that, when dispatched, change the state.
// There are four dispatchable types:
//   - Types that implement the [State] constraint.
//     For example, types that embed the [EmptyState].
//   - [StateAndEffects]
//   - [Action]
//   - [ActionAndPayload]
type Dispatchable interface {
	IAmDispatchable()
}

type StateAndEffects[S State] struct {
	State   S
	Effects []Effect
}

// IAmDispatchable makes StateAndEffects [Dispatchable].
func (_ StateAndEffects[S]) IAmDispatchable() {}

// ActionAndPayload contains an [Action] and [Payload].
// When the action is dispatched, it receives the current state as its first argument and the payload as its second argument.
type ActionAndPayload[S State] struct {
	Action  Action[S]
	Payload Payload
}

// IAmDispatchable makes ActionAndPayload [Dispatchable].
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

// VNodeKind indicates the type of [VNode].
type VNodeKind int

// Each constant corresponds to an element's [nodeType].
// Use [H] to create an ElementNode VNode.
// Use [Text] or [Textf] to create a TextNode VNode.
//
// [nodeType]: https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeType
const (
	// ElementNode indicates a VNode that renders an element node.
	ElementNode VNodeKind = 1
	// TextNode indicates a VNode that renders text inside an element node.
	TextNode VNodeKind = 3
)

type VNode struct {
	props    HProps
	children vKids
	node     window.Element // Can be empty
	tag      string
	memoView func(data MemoData) *VNode
	memoData MemoData
	kind     VNodeKind
}

// Props returns the VNode's properties.
func (n VNode) Props() HProps {
	return n.props
}

// Children returns the VNode's children.
func (n VNode) Children() []*VNode {
	return n.children
}

// Tag returns the VNode's tag.
func (n VNode) Tag() string {
	return n.tag
}

// Kind returns the VNode's [VNodeKind].
func (n VNode) Kind() VNodeKind {
	return n.kind
}

func (n VNode) key() option[string] {
	return n.props.key()
}

type vKids []*VNode

func (v vKids) getKey(i int) option[string] {
	if i < len(v) {
		return v[i].key()
	}
	return option[string]{}
}

func (v vKids) get(i int) *VNode {
	if i < len(v) {
		return v[i]
	}
	return nil
}

package window

import (
	"github.com/macabot/hypp/js"
)

type EventTarget struct {
	// js.Value cannot be embedded because the name conflicts with the method [EventTarget.Value].
	V js.Value
}

// Value returns the value of the EventTarget.
func (t EventTarget) Value() string {
	return t.V.Get("value").String()
}

type Event struct {
	js.Value
}

// Type returns a string containing the event's type.
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/type
func (e Event) Type() string {
	return e.Value.Get("type").String()
}

// PreventDefault tells the user agent that if the event does not get explicitly handled, its default action should not be taken as it normally would be.
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/preventDefault
func (e Event) PreventDefault() {
	e.Value.Call("preventDefault")
}

// StopImmediatePropagation prevents other listeners of the same event from being called.
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/stopImmediatePropagation
func (e Event) StopImmediatePropagation() {
	e.Value.Call("stopImmediatePropagation")
}

// StopPropagation prevents further propagation of the current event in the capturing and bubbling phases.
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/stopPropagation
func (e Event) StopPropagation() {
	e.Value.Call("stopPropagation")
}

// Target returns a reference to the object onto which the event was dispatched.
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/target
func (e Event) Target() EventTarget {
	return EventTarget{e.Value.Get("target")}
}

type EventListenerID struct {
	js.Func
}

type EventListener func(Event)

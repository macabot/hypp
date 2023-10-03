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

func (e Event) Type() string {
	return e.Value.Get("type").String()
}

func (e Event) PreventDefault() {
	e.Value.Call("preventDefault")
}

func (e Event) StopImmediatePropagation() {
	e.Value.Call("stopImmediatePropagation")
}

func (e Event) StopPropagation() {
	e.Value.Call("stopPropagation")
}

func (e Event) Target() EventTarget {
	return EventTarget{e.Value.Get("target")}
}

type EventListenerID struct {
	js.Value
}

type EventListener func(Event)

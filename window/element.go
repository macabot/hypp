package window

import (
	"fmt"

	"github.com/macabot/hypp/js"
	"github.com/macabot/hypp/util"
)

// Element represents an HTML element.
// See https://developer.mozilla.org/en-US/docs/Web/API/Element
type Element struct {
	js.Value
}

// RemoveEventListener removes an [EventListener] previously registered with [Node.AddEventListener].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener
func (e Element) RemoveEventListener(kind string, listenerID EventListenerID) {
	e.Value.Call("removeEventListener", kind, listenerID.Value)
}

// AddEventListener sets up a function that will be called whenever the specified event is delivered to the [Node].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener
func (e Element) AddEventListener(kind string, listener EventListener) EventListenerID {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		listener(Event{args[0]})
		return nil
	})
	e.Value.Call("addEventListener", kind, f)
	return EventListenerID{js.ValueOf(f)}
}

// ParentNode returns the parent [Node].
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/parentNode
func (e Element) ParentNode() Element {
	return Element{e.Value.Get("parentNode")}
}

// NodeType returns an integer that identifies what the node is.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeType
func (e Element) NodeType() int {
	return e.Value.Get("nodeType").Int()
}

func (e Element) NodeValue() string {
	return e.Value.Get("nodeValue").String()
}

func (e Element) SetNodeValue(nodeValue string) {
	e.Value.Set("nodeValue", nodeValue)
}

func (e Element) NodeName() string {
	return e.Value.Get("nodeName").String()
}

func (e Element) ChildNodes() []Element {
	children := e.Value.Get("childNodes")
	l := children.Length()
	out := make([]Element, l)
	for i := 0; i < l; i++ {
		out[i] = Element{children.Index(i)}
	}
	return out
}

func (e Element) InsertBefore(newNode, referenceNode Element) Element {
	return Element{e.Value.Call(
		"insertBefore",
		newNode.Value,
		referenceNode.Value,
	)}
}

func (e Element) RemoveChild(child Element) {
	e.Value.Call("removeChild", child.Value)
}

func (e Element) Get(name string) util.Option[any] {
	if !e.In(name) {
		return util.Option[any]{}
	}
	v := e.Value.Get(name)
	kind := v.Type()
	switch kind {
	case js.TypeUndefined, js.TypeNull:
		return util.Option[any]{OK: true}
	case js.TypeBoolean:
		return util.Option[any]{V: v.Bool(), OK: true}
	case js.TypeNumber:
		if js.Global().Get("Number").Call("isInteger", v).Bool() {
			return util.Option[any]{V: v.Int(), OK: true}
		} else {
			return util.Option[any]{V: v.Float(), OK: true}
		}
	case js.TypeString:
		return util.Option[any]{V: v.String(), OK: true}
	default:
		panic(fmt.Errorf("hypp: cannot get node property of type '%s'", kind))
	}
}

func (e Element) In(name string) bool {
	getPrototypeOf := js.Global().Get("Object").Get("getPrototypeOf")
	for v := e.Value; !v.IsNull(); v = getPrototypeOf.Invoke(v) {
		if v.Call("hasOwnProperty", name).Bool() {
			return true
		}
	}
	return false
}

func validateValue(value any) {
	switch value.(type) {
	case nil, bool, int, float64, string:
		// Do nothing
	default:
		panic(fmt.Errorf("hypp: expected nil, bool, int, float64 or string. Got %+v of type %T\n", value, value))
	}
}

func (e Element) Set(name string, value any) {
	validateValue(value)
	e.Value.Set(name, value)
}

func (e Element) AppendChild(child Element) Element {
	return Element{e.Value.Call("appendChild", child.Value)}
}

func (e Element) RemoveAttribute(name string) {
	e.Value.Call("removeAttribute", name)
}

func (e Element) SetAttribute(name string, value any) {
	e.Value.Call("setAttribute", name, value)
}

// TODO move from package 'window' to 'hypp'
// func (e Element) Eventsx() Events {
// 	v := e.Value
// 	if v.Get("events").IsUndefined() {
// 		e.Value.Set("events", map[string]any{})
// 	}
// 	return Events{v.Get("events")}
// }

func (e Element) SetStyleProperty(propertyName, value string) {
	e.Value.Get("style").Call("setProperty", propertyName, value)
}

func (e Element) SetStyle(name, value string) {
	e.Value.Get("style").Set(name, value)
}

func (e Element) EventListenerID(kind string) EventListenerID {
	listeners := e.Value.Get("eventListeners")
	if listeners.IsUndefined() {
		return EventListenerID{}
	}
	listener := listeners.Get(kind)
	if listener.IsUndefined() {
		return EventListenerID{}
	}
	return EventListenerID{listener}
}

func (e Element) SetEventListenerID(kind string, eventListenerID EventListenerID) {
	v := e.Value
	id := eventListenerID.Value
	listeners := v.Get("eventListeners")
	if listeners.IsUndefined() {
		v.Set("eventListeners", map[string]any{kind: id})
	} else {
		listeners.Set(kind, id)
	}
}

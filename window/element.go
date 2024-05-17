package window

import (
	"fmt"

	"github.com/macabot/hypp/js"
)

// Element represents an HTML element.
// See https://developer.mozilla.org/en-US/docs/Web/API/Element
type Element struct {
	js.Value
}

// RemoveEventListener removes an [EventListener] previously registered with [Node.AddEventListener].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener
//
// It frees up resources allocated for the event listener.
func (e Element) RemoveEventListener(kind string, listenerID EventListenerID) {
	e.Value.Call("removeEventListener", kind, listenerID.Value)
	listenerID.Release()
}

// AddEventListener sets up a function that will be called whenever the specified event is delivered to the [Node].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener
//
// It returns an EventListenerID that contains the js.Func corresponding to the given EventListener.
// Use this EventListenerID when calling [Element.RemoveEventListener].
// This ensures the resources allocated by the js.Func are released.
func (e Element) AddEventListener(kind string, listener EventListener) EventListenerID {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		listener(Event{args[0]})
		return nil
	})
	e.Value.Call("addEventListener", kind, f)
	return EventListenerID{f}
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

// NodeValue returns the value of the node.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeValue
func (e Element) NodeValue() string {
	return e.Value.Get("nodeValue").String()
}

// SetNodeValue sets the value of the node.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeValue
func (e Element) SetNodeValue(nodeValue string) {
	e.Value.Set("nodeValue", nodeValue)
}

// NodeName returns the name of the node.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeName
func (e Element) NodeName() string {
	return e.Value.Get("nodeName").String()
}

// ChildNodes returns the child nodes.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/childNodes
func (e Element) ChildNodes() []Element {
	children := e.Value.Get("childNodes")
	l := children.Length()
	out := make([]Element, l)
	for i := 0; i < l; i++ {
		out[i] = Element{children.Index(i)}
	}
	return out
}

// InsertBefore inserts the newNode before the referenceNode as a child.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/insertBefore
func (e Element) InsertBefore(newNode, referenceNode Element) Element {
	return Element{e.Value.Call(
		"insertBefore",
		newNode.Value,
		referenceNode.Value,
	)}
}

// RemoveChild removes the child node.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/removeChild
func (e Element) RemoveChild(child Element) {
	e.Value.Call("removeChild", child.Value)
}

func validateValue(value any) {
	switch value.(type) {
	case nil, bool, int, float64, string:
		// Do nothing
	default:
		panic(fmt.Errorf("hypp/window: expected nil, bool, int, float64 or string. Got %+v of type %T", value, value))
	}
}

func (e Element) Set(name string, value any) {
	validateValue(value)
	e.Value.Set(name, value)
}

// AppendChild adds a node to the end of the list of children.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node/appendChild
func (e Element) AppendChild(child Element) Element {
	return Element{e.Value.Call("appendChild", child.Value)}
}

// RemoveAttribute removes the attribute with the specified name from the element.
// See https://developer.mozilla.org/en-US/docs/Web/API/Element/removeAttribute
func (e Element) RemoveAttribute(name string) {
	e.Value.Call("removeAttribute", name)
}

// SetAttribute sets the value of an attribute on the specified element. If the attribute already exists, the value is updated; otherwise a new attribute is added with the specified name and value.
// See https://developer.mozilla.org/en-US/docs/Web/API/Element/setAttribute
func (e Element) SetAttribute(name string, value any) {
	e.Value.Call("setAttribute", name, value)
}

func (e Element) SetStyleProperty(propertyName, value string) {
	e.Value.Get("style").Call("setProperty", propertyName, value)
}

func (e Element) SetStyle(name, value string) {
	e.Value.Get("style").Set(name, value)
}

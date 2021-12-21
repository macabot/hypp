package js

import (
	"syscall/js"

	"github.com/macabot/hypp"
)

func document() js.Value {
	return js.Global().Get("document")
}

type Driver struct{}

func (d Driver) CreateTextNode(data string) Node {
	return document().Call("createTextNode", data)
}

func elementCreationOptionsToValue(options hypp.Option[hypp.ElementCreationOptions]) js.Value {
	out := js.Undefined()
	if options.OK {
		out = map[string]interface{}{
			"is": options.Is,
		}
	}
	return out
}

func (d Driver) CreateElementNS(namespaceURI, qualifiedName string, options hypp.Option[hypp.ElementCreationOptions]) Node {
	return document().Call("createElementNS", namespaceURI, qualifiedName, elementCreationOptionsToValue(options))
}

func (d Driver) CreateElement(tagName string, options Option[hypp.ElementCreationOptions]) Node {
	return document().Call("createElement", tagName, elementCreationOptionsToValue(options))
}

func (d Driver) Enqueue(render func()) {
	js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		render()
	}))
}

type Node js.Value

func (n Node) JSValue() js.Value {
	return js.Value(n)
}

func (n Node) ParentNode() hypp.Option[hypp.Node] {
	parent := js.Value(n)
	return hypp.Option[hypp.Node]{
		V:  Node(parent),
		OK: !parent.IsNull(),
	}
}

func (n Node) NodeType() int {
	return js.Value(n).Get("nodeType").Int()
}

func (n Node) NodeValue() string {
	return js.Value(n).Get("nodeValue").String()
}

func (n Node) SetNodeValue(nodeValue string) {
	js.Value(n).Set("nodeValue", nodeValue)
}

func (n Node) NodeName() string {
	return js.Value(n).Get("nodeName").String()
}

func (n Node) ChildNodes() []hypp.Node {
	children := js.Value(n).Get("childNodes")
	l := children.Length()
	out := make([]hypp.Node, l)
	for i := 0; i < l; i++ {
		out[i] = Node(children.Index(i))
	}
	return out
}

func (n Node) InsertBefore(newNode, referenceNode hypp.Node) hypp.Node {
	return Node(js.Value(n).Call("insertBefore", newNode, referenceNode))
}

func (n Node) RemoveChild(child hypp.Node) {
	js.Value(n).Call("removeChild", child)
}

func (n Node) Get(name string) hypp.Value {
	return js.Value(n).Get(name)
}

func (n Node) Has(name string) bool {
	return js.Value(n).Has(name)
}

func (n Node) Set(name string, value interface{}) {
	js.Value(n).Set(name, value)
}

func (n Node) AppendChild(child hypp.Node) hypp.Node {
	return Node(js.Value(n).Call("appendChild", child))
}

func (n Node) RemoveEventListener(kind string, listener hypp.EventListener) {
	js.Value(n).Call("removeEventListener", kind, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener(Event(args[0]))
	}))
}

func (n Node) AddEventListener(kind string, listener hypp.EventListener) {
	js.Value(n).Call("addEventListener", kind, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener[Event(args[0])]
	}))
}

func (n Node) RemoveAttribute(name string) {
	js.Value(n).Call("removeAttribute", name)
}

func (n Node) SetAttribute(name string, value interface{}) {
	js.Value(n).Call("setAttribute", name, value)
}

func (n Node) Style() hypp.Style {
	return Style(js.Value(n).Get("style"))
}

type Events js.Value

func (e Events) Set(name string, value interface{}) {
	js.Value(e).Set(name, value)
}

func (e Events) Get(name string) hypp.Value {
	return js.Value(e).Get(name)
}

type Style js.Value

func (s Style) SetProperty(propertyName, value string) {
	js.Value(s).Call("setProperty", propertyName, value)
}

func (s Style) Set(name, value string) {
	js.Value(s).Set(name, value)
}

func (s Style) Get(name string) string {
	js.Value(s).Get(name).String()
}

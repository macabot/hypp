package js

import (
	"fmt"
	"sync"
	"syscall/js"

	"github.com/macabot/hypp"
)

func document() js.Value {
	return js.Global().Get("document")
}

var _ hypp.Driver = Driver{}

type Driver struct{}

func (d Driver) CreateTextNode(data string) hypp.Node {
	return Node(document().Call("createTextNode", data))
}

func elementCreationOptionsToValue(options hypp.Option[hypp.ElementCreationOptions]) js.Value {
	out := js.Undefined()
	if options.OK {
		out = js.ValueOf(map[string]interface{}{
			"is": options.V.Is,
		})
	}
	return out
}

func (d Driver) CreateElementNS(namespaceURI, qualifiedName string, options hypp.Option[hypp.ElementCreationOptions]) hypp.Node {
	return Node(document().Call("createElementNS", namespaceURI, qualifiedName, elementCreationOptionsToValue(options)))
}

func (d Driver) CreateElement(tagName string, options hypp.Option[hypp.ElementCreationOptions]) hypp.Node {
	return Node(document().Call("createElement", tagName, elementCreationOptionsToValue(options)))
}

func (d Driver) Enqueue(render func()) {
	js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		render()
		return nil
	}))
}

func hyppNodeToValue(node hypp.Node) js.Value {
	if node == nil {
		return js.Null()
	}
	return node.(Node).JSValue()
}

var _ hypp.Node = Node{}

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
	return Node(js.Value(n).Call(
		"insertBefore",
		hyppNodeToValue(newNode),
		hyppNodeToValue(referenceNode),
	))
}

func (n Node) RemoveChild(child hypp.Node) {
	js.Value(n).Call("removeChild", hyppNodeToValue(child))
}

func (n Node) Get(name string) hypp.Option[interface{}] {
	if !n.Has(name) {
		return hypp.Option[interface{}]{}
	}
	v := js.Value(n).Get(name)
	kind := v.Type()
	switch kind {
	case js.TypeUndefined, js.TypeNull:
		return hypp.Option[interface{}]{OK: true}
	case js.TypeBoolean:
		return hypp.Option[interface{}]{V: v.Bool(), OK: true}
	case js.TypeNumber:
		if js.Global().Get("Number").Call("isInteger", v).Bool() {
			return hypp.Option[interface{}]{V: v.Int(), OK: true}
		} else {
			return hypp.Option[interface{}]{V: v.Float(), OK: true}
		}
	case js.TypeString:
		return hypp.Option[interface{}]{V: v.String(), OK: true}
	default:
		panic(fmt.Errorf("js: cannot get Node property of type '%s'", kind))
	}
}

func (n Node) Has(name string) bool {
	return js.Value(n).Call("hasOwnProperty", name).Bool()
}

func (n Node) Set(name string, value interface{}) {
	js.Value(n).Set(name, value)
}

func (n Node) AppendChild(child hypp.Node) hypp.Node {
	return Node(js.Value(n).Call("appendChild", hyppNodeToValue(child)))
}

func (n Node) RemoveEventListener(kind string, listener hypp.EventListener) {
	js.Value(n).Call("removeEventListener", kind, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener(Event(args[0]))
		return nil
	}))
}

func (n Node) AddEventListener(kind string, listener hypp.EventListener) {
	js.Value(n).Call("addEventListener", kind, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener(Event(args[0]))
		return nil
	}))
}

func (n Node) RemoveAttribute(name string) {
	js.Value(n).Call("removeAttribute", name)
}

func (n Node) SetAttribute(name string, value interface{}) {
	js.Value(n).Call("setAttribute", name, value)
}

func (n Node) Events() hypp.Events {
	v := js.Value(n)
	if v.Get("events").IsUndefined() {
		js.Value(n).Set("events", map[string]interface{}{})
	}
	return Events(v.Get("events"))
}

func (n Node) Style() hypp.Style {
	return Style(js.Value(n).Get("style"))
}

var _ hypp.Events = Events{}

type Events js.Value

func (e Events) JSValue() js.Value {
	return js.Value(e)
}

type GlobalDispatchables struct {
	mu sync.Mutex
	i int
	v map[int]hypp.Dispatchable
}

func (g *GlobalDispatchables) Set(value hypp.Dispatchable) int {
	g.mu.Lock()
	i := g.i
	g.v[i] = value
	g.i++
	g.mu.Unlock()
	return i
}

func (g GlobalDispatchables) Get(i int) hypp.Dispatchable {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.v[i]
}

func (g *GlobalDispatchables) Del(i int) {
	g.mu.Lock()
	delete(g.v, i)
	g.mu.Unlock()
}

var globalNodeEvents = GlobalDispatchables{v: map[int]hypp.Dispatchable{}}

func (e Events) Set(name string, value hypp.Dispatchable) {
	i := globalNodeEvents.Set(value)
	js.Value(e).Set(name, i)
}

func (e Events) Get(name string) hypp.Dispatchable {
	return globalNodeEvents.Get(js.Value(e).Get(name).Int())
}

func (e Events) Del(name string) {
	globalNodeEvents.Del(js.Value(e).Get(name).Int())
}

var _ hypp.Event = Event{}

type Event js.Value

func (e Event) JSValue() js.Value {
	return js.Value(e)
}

// func (e Event) IAmDispatchable() {}

func (e Event) Type() string {
	return js.Value(e).Get("type").String()
}

func (e Event) PreventDefault() {
	js.Value(e).Call("preventDefault")
}

func (e Event) Target() hypp.EventTarget {
	return EventTarget(js.Value(e).Get("target"))
}

var _ hypp.EventTarget = EventTarget{}

type EventTarget js.Value

func (e EventTarget) JSValue() js.Value {
	return js.Value(e)
}

func (e EventTarget) Value() interface{} {
	return js.Value(e).Get("value")
}

var _ hypp.Style = Style{}

type Style js.Value

func (s Style) JSValue() js.Value {
	return js.Value(s)
}

func (s Style) SetProperty(propertyName, value string) {
	js.Value(s).Call("setProperty", propertyName, value)
}

func (s Style) Set(name, value string) {
	js.Value(s).Set(name, value)
}

func (s Style) Get(name string) string {
	return js.Value(s).Get(name).String()
}

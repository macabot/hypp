package js

import (
	"errors"
	"fmt"
	"sync"
	"syscall/js"

	"github.com/macabot/hypp"
)

func document() js.Value {
	return js.Global().Get("document")
}

type Driver struct{}

var _ hypp.Driver = Driver{}

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

func (d Driver) Window() hypp.Window {
	return Window(js.Global())
}

func (d Driver) ValidateAppPropsNode(node hypp.Node) error {
	if node == nil {
		return errors.New("hypp/driver/js: AppProps.Node cannot be nil")
	} else if js.Value(node.ParentNode().(Node)).IsNull() {
		return errors.New("hypp/driver/js: AppProps.Node must have a parent")
	}
	return nil
}

type Window js.Value

var _ hypp.Window = Window{}

func (w Window) EscapeToValue() hypp.Value {
	return EscapeToValuer(w).EscapeToValue()
}

func (w Window) RemoveEventListener(kind string, listenerID hypp.EventListenerID) {
	EventTarget(w).RemoveEventListener(kind, listenerID)
}

func (w Window) AddEventListener(kind string, listener hypp.EventListener) hypp.EventListenerID {
	return EventTarget(w).AddEventListener(kind, listener)
}

func (w Window) RequestAnimationFrame(f func()) int {
	return js.Value(w).Call(
		"requestAnimationFrame",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			f()
			return nil
		}),
	).Int()
}

func hyppNodeToValue(node hypp.Node) js.Value {
	if node == nil {
		return js.Null()
	}
	return js.Value(node.(Node))
}

type EventTarget js.Value

var _ hypp.EventTarget = EventTarget{}

func (e EventTarget) RemoveEventListener(kind string, listenerID hypp.EventListenerID) {
	js.Value(e).Call("removeEventListener", kind, js.Value(listenerID.(EventListenerID)))
}

func (e EventTarget) AddEventListener(kind string, listener hypp.EventListener) hypp.EventListenerID {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener(Event(args[0]))
		return nil
	})
	js.Value(e).Call("addEventListener", kind, f)
	return EventListenerID(js.ValueOf(f))
}

type EventListenerID js.Value

var _ hypp.EventListenerID = EventListenerID{}

func (e EventListenerID) IAmAnEventListenerID() {}

type Node js.Value

var _ hypp.Node = Node{}

func (n Node) ParentNode() hypp.Node {
	return Node(js.Value(n).Get("parentNode"))
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
	child.Events().(Events).deleteAll()
	js.Value(n).Call("removeChild", hyppNodeToValue(child))
}

func (n Node) Get(name string) hypp.Option[interface{}] {
	if !n.In(name) {
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

func (n Node) In(name string) bool {
	getPrototypeOf := js.Global().Get("Object").Get("getPrototypeOf")
	for v := js.Value(n); !v.IsNull(); v = getPrototypeOf.Invoke(v) {
		if v.Call("hasOwnProperty", name).Bool() {
			return true
		}
	}
	return false
}

func validateValue(value interface{}) {
	switch value.(type) {
	case nil, bool, int, float64, string:
		// Do nothing
	default:
		panic(fmt.Errorf("hypp: expected nil, bool, int, float64 or string. Got %+v of type %T\n", value, value))
	}
}

func (n Node) Set(name string, value interface{}) {
	validateValue(value)
	js.Value(n).Set(name, value)
}

func (n Node) AppendChild(child hypp.Node) hypp.Node {
	return Node(js.Value(n).Call("appendChild", hyppNodeToValue(child)))
}

func (n Node) RemoveEventListener(kind string, listenerID hypp.EventListenerID) {
	EventTarget(n).RemoveEventListener(kind, listenerID)
}

func (n Node) AddEventListener(kind string, listener hypp.EventListener) hypp.EventListenerID {
	return EventTarget(n).AddEventListener(kind, listener)
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

func (n Node) SetStyleProperty(propertyName, value string) {
	js.Value(n).Get("style").Call("setProperty", propertyName, value)
}

func (n Node) SetStyle(name, value string) {
	js.Value(n).Get("style").Set(name, value)
}

func (n Node) EventListenerID(kind string) hypp.EventListenerID {
	listeners := js.Value(n).Get("eventListeners")
	if listeners.IsUndefined() {
		return nil
	}
	listener := listeners.Get(kind)
	if listener.IsUndefined() {
		return nil
	}
	return EventListenerID(listener)
}

func (n Node) SetEventListenerID(kind string, eventListenerID hypp.EventListenerID) {
	v := js.Value(n)
	id := js.Value(eventListenerID.(EventListenerID))
	listeners := v.Get("eventListeners")
	if listeners.IsUndefined() {
		v.Set("eventListeners", map[string]interface{}{kind: id})
	} else {
		listeners.Set(kind, id)
	}
}

type Events js.Value

var _ hypp.Events = Events{}

type dispatchablesRepo struct {
	mu sync.Mutex
	i  int
	v  map[int]hypp.Dispatchable
}

func (g *dispatchablesRepo) Add(value hypp.Dispatchable) int {
	g.mu.Lock()
	i := g.i
	g.v[i] = value
	g.i++
	g.mu.Unlock()
	return i
}

func (g *dispatchablesRepo) Set(i int, value hypp.Dispatchable) {
	g.mu.Lock()
	g.v[i] = value
	g.mu.Unlock()
}

func (g *dispatchablesRepo) Del(i int) {
	g.mu.Lock()
	delete(g.v, i)
	g.mu.Unlock()
}

func (g *dispatchablesRepo) Get(i int) hypp.Dispatchable {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.v[i]
}

var globalNodeEvents = dispatchablesRepo{v: map[int]hypp.Dispatchable{}}

func (e Events) Set(name string, value hypp.Dispatchable) {
	v := js.Value(e)
	if p := v.Get(name); p.Type() == js.TypeNumber {
		globalNodeEvents.Set(p.Int(), value)
	} else {
		i := globalNodeEvents.Add(value)
		v.Set(name, i)
	}
}

func (e Events) Get(name string) hypp.Dispatchable {
	v := js.Value(e)
	if p := v.Get(name); p.Type() == js.TypeNumber {
		return globalNodeEvents.Get(p.Int())
	}
	return nil
}

func (e Events) Del(name string) {
	v := js.Value(e)
	if p := v.Get(name); p.Type() == js.TypeNumber {
		globalNodeEvents.Del(p.Int())
		v.Delete(name)
	}
}

func (e Events) deleteAll() {
	names := js.Global().Get("Object").Call("keys", js.Value(e))
	l := names.Length()
	for i := 0; i < l; i++ {
		name := names.Index(i).String()
		e.Del(name)
	}
}

type EscapeToValuer js.Value

var _ hypp.EscapeToValuer = EscapeToValuer{}

func (e EscapeToValuer) EscapeToValue() hypp.Value {
	return Value{js.Value(e)}
}

type Event js.Value

var _ hypp.Event = Event{}

func (e Event) EscapeToValue() hypp.Value {
	return EscapeToValuer(e).EscapeToValue()
}

func (e Event) Type() string {
	return js.Value(e).Get("type").String()
}

func (e Event) PreventDefault() {
	js.Value(e).Call("preventDefault")
}

func (e Event) StopImmediatePropagation() {
	js.Value(e).Call("stopImmediatePropagation")
}

func (e Event) StopPropagation() {
	js.Value(e).Call("stopPropagation")
}

func (e Event) Target() hypp.EventTargetValuer {
	return EventTargetValuer(js.Value(e).Get("target"))
}

type EventTargetValuer js.Value

var _ hypp.EventTargetValuer = EventTargetValuer{}

func (e EventTargetValuer) Value() string {
	return js.Value(e).Get("value").String()
}

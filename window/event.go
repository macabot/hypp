package window

import (
	"sync"

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

// TODO move all code below from package 'window' to 'hypp'.

type dispatchablesRepo struct {
	mu sync.Mutex
	i  int
	v  map[int]Dispatchable
}

func (g *dispatchablesRepo) Add(value Dispatchable) int {
	g.mu.Lock()
	i := g.i
	g.v[i] = value
	g.i++
	g.mu.Unlock()
	return i
}

func (g *dispatchablesRepo) Set(i int, value Dispatchable) {
	g.mu.Lock()
	g.v[i] = value
	g.mu.Unlock()
}

func (g *dispatchablesRepo) Del(i int) {
	g.mu.Lock()
	delete(g.v, i)
	g.mu.Unlock()
}

func (g *dispatchablesRepo) Get(i int) Dispatchable {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.v[i]
}

var globalNodeEvents = dispatchablesRepo{v: map[int]Dispatchable{}}

type Events struct {
	js.Value
}

func (e Events) Set(name string, value Dispatchable) {
	v := e.Value
	if p := v.Get(name); p.Type() == js.TypeNumber {
		globalNodeEvents.Set(p.Int(), value)
	} else {
		i := globalNodeEvents.Add(value)
		v.Set(name, i)
	}
}

func (e Events) Get(name string) Dispatchable {
	v := e.Value
	if p := v.Get(name); p.Type() == js.TypeNumber {
		return globalNodeEvents.Get(p.Int())
	}
	return nil
}

func (e Events) Del(name string) {
	v := e.Value
	if p := v.Get(name); p.Type() == js.TypeNumber {
		globalNodeEvents.Del(p.Int())
		v.Delete(name)
	}
}

func (e Events) deleteAll() {
	names := js.Global().Get("Object").Call("keys", e.Value)
	l := names.Length()
	for i := 0; i < l; i++ {
		name := names.Index(i).String()
		e.Del(name)
	}
}

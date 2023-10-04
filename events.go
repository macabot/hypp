package hypp

import (
	"sync"

	"github.com/macabot/hypp/js"
	"github.com/macabot/hypp/window"
)

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

type events struct {
	js.Value
}

func (e events) Set(name string, value Dispatchable) {
	v := e.Value
	if p := v.Get(name); p.Type() == js.TypeNumber {
		globalNodeEvents.Set(p.Int(), value)
	} else {
		i := globalNodeEvents.Add(value)
		v.Set(name, i)
	}
}

func (e events) Get(name string) Dispatchable {
	v := e.Value
	if p := v.Get(name); p.Type() == js.TypeNumber {
		return globalNodeEvents.Get(p.Int())
	}
	return nil
}

func (e events) Del(name string) {
	v := e.Value
	if p := v.Get(name); p.Type() == js.TypeNumber {
		globalNodeEvents.Del(p.Int())
		v.Delete(name)
	}
}

func (e events) deleteAll() {
	names := js.Global().Get("Object").Call("keys", e.Value)
	l := names.Length()
	for i := 0; i < l; i++ {
		name := names.Index(i).String()
		e.Del(name)
	}
}

func removeChild(parent, child window.Element) {
	events{child.Value.Get("events")}.deleteAll()
	parent.RemoveChild(child)
}

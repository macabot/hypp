package window

import "github.com/macabot/hypp/js"

type Win struct {
	js.Value
}

func Window() Win {
	return Win{js.Global()}
}

func (w Win) RequestAnimationFrame(f func()) int {
	return w.Call(
		"requestAnimationFrame",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			f()
			return nil
		}),
	).Int()
}

// RemoveEventListener removes an [EventListener] previously registered with [Node.AddEventListener].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener
func (w Win) RemoveEventListener(kind string, listenerID EventListenerID) {
	w.Value.Call("removeEventListener", kind, listenerID.Value)
}

// AddEventListener sets up a function that will be called whenever the specified event is delivered to the [Node].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener
func (w Win) AddEventListener(kind string, listener EventListener) EventListenerID {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		listener(Event{args[0]})
		return nil
	})
	w.Value.Call("addEventListener", kind, f)
	return EventListenerID{js.ValueOf(f)}
}

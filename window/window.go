package window

import "github.com/macabot/hypp/js"

func RequestAnimationFrame(f func()) int {
	return js.Global().Call(
		"requestAnimationFrame",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			f()
			return nil
		}),
	).Int()
}

// RemoveEventListener removes an [EventListener] previously registered with [Node.AddEventListener].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener
func RemoveEventListener(kind string, listenerID EventListenerID) {
	js.Global().Call("removeEventListener", kind, listenerID.Value)
}

// AddEventListener sets up a function that will be called whenever the specified event is delivered to the [Node].
// See https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener
func AddEventListener(kind string, listener EventListener) EventListenerID {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		listener(Event{args[0]})
		return nil
	})
	js.Global().Call("addEventListener", kind, f)
	return EventListenerID{js.ValueOf(f)}
}

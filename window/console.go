package window

import "github.com/macabot/hypp/js"

type Con struct {
	js.Value
}

// Console returns the console.
// See https://developer.mozilla.org/en-US/docs/Web/API/console
func Console() Con {
	return Con{js.Global().Get("console")}
}

// Log outputs a message to the web console.
// See https://developer.mozilla.org/en-US/docs/Web/API/console/log
func (c Con) Log(args ...any) {
	c.Value.Call("log", args...)
}

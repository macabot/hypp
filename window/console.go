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

// Debug outputs a debug message to the web console.
// See https://developer.mozilla.org/en-US/docs/Web/API/console/debug_static
func (c Con) Debug(args ...any) {
	c.Value.Call("debug", args...)
}

// Info outputs an informational message to the web console.
// See https://developer.mozilla.org/en-US/docs/Web/API/console/info_static
func (c Con) Info(args ...any) {
	c.Value.Call("info", args...)
}

// Log outputs a message to the web console.
// See https://developer.mozilla.org/en-US/docs/Web/API/console/log
func (c Con) Log(args ...any) {
	c.Value.Call("log", args...)
}

// Warn outputs a warning message to the web console.
// See https://developer.mozilla.org/en-US/docs/Web/API/console/warn_static
func (c Con) Warn(args ...any) {
	c.Value.Call("warn", args...)
}

// Error outputs an error message to the web console.
// See https://developer.mozilla.org/en-US/docs/Web/API/console/error_static
func (c Con) Error(args ...any) {
	c.Value.Call("error", args...)
}

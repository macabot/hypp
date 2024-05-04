package hypp_test

import (
	"testing"
)

// When creating AppProps, the state is initialized to its default value.
// The dispatch function only updates the DOM if the state has changed.
// This test ensures that the DOM is always updated with calling the dispatch function with AppProps.Init, even if AppProps.Init resolves to the state's default value.
func TestViewIsCalledWithInitWhenStateIsUnchanged(t *testing.T) {
	t.Skip("Not yet implemented")
	// node := window.Document().CreateElement("div", nil)
	// hypp.App[struct{}](hypp.AppProps[struct{}]{
	// 	View: func(state struct{}) *hypp.VNode {
	// 		return html.Main(nil)
	// 	},
	// 	Node: node,
	// })
	// TODO continue
}

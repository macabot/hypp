package window

import "github.com/macabot/hypp/js"

type Doc struct {
	js.Value
}

// Document returns the window's document.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document
func Document() Doc {
	return Doc{js.Global().Get("document")}
}

// CreateTextNode creates a new text node.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/createTextNode
func (d Doc) CreateTextNode(data string) Element {
	return Element{d.Call("createTextNode", data)}
}

// ElementCreationOptions represents the options when creating an element.
type ElementCreationOptions struct {
	// Is is the tag name of a custom element.
	Is string
}

// Value returns the js.Value representation.
func (o *ElementCreationOptions) Value() js.Value {
	if o == nil {
		return js.Undefined()
	}
	return js.ValueOf(map[string]any{
		"is": o.Is,
	})
}

// CreateElementNS creates an element with the specified namespace URI and qualified name.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/createElementNS
func (d Doc) CreateElementNS(namespaceURI, qualifiedName string, options *ElementCreationOptions) Element {
	return Element{d.Call("createElementNS", namespaceURI, qualifiedName, options.Value())}
}

// CreateElement creates the HTML element specified by tagName.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/createElement
func (d Doc) CreateElement(tagName string, options *ElementCreationOptions) Element {
	return Element{d.Call("createElement", tagName, options.Value())}
}

// GetElementById returns the [Element] whose id property matches the specified string.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementById
func (d Doc) GetElementById(id string) Element {
	return Element{d.Call("getElementById", id)}
}

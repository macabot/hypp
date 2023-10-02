package window

import "github.com/macabot/hypp/js"

type Doc struct {
	js.Value
}

func Document() Doc {
	return Doc{js.Global().Get("document")}
}

func (d Doc) CreateTextNode(data string) Element {
	return Element{d.Call("createTextNode", data)}
}

type ElementCreationOptions struct {
	Is string
}

func (o *ElementCreationOptions) Value() js.Value {
	if o == nil {
		return js.Undefined()
	}
	return js.ValueOf(map[string]any{
		"is": o.Is,
	})
}

func (d Doc) CreateElementNS(namespaceURI, qualifiedName string, options *ElementCreationOptions) Element {
	return Element{d.Call("createElementNS", namespaceURI, qualifiedName, options.Value())}
}

func (d Doc) CreateElement(tagName string, options *ElementCreationOptions) Element {
	return Element{d.Call("createElement", tagName, options.Value())}
}

func (d Doc) GetElementById(id string) Element {
	return Element{d.Call("getElementById", id)}
}

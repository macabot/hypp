package html

import (
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/macabot/hypp"
)

var ssrNode = 1
var textNode = 3

type Driver struct{}

var _ hypp.Driver = Driver{}

func (d Driver) CreateTextNode(data string) hypp.Node {
	return &Node{
		nodeType:  textNode,
		nodeValue: data,
		nodeName:  "#text",
	}
}

func (d Driver) CreateElementNS(namespaceURI, qualifiedName string, options hypp.Option[hypp.ElementCreationOptions]) hypp.Node {
	return &Node{
		nodeType:     ssrNode,
		nodeName:     qualifiedName,
		namespaceURI: namespaceURI,
	}
}

func (d Driver) CreateElement(tagName string, options hypp.Option[hypp.ElementCreationOptions]) hypp.Node {
	return d.CreateElementNS("http://www.w3.org/1999/xhtml", tagName, options)
}

func (d Driver) Window() hypp.Window {
	return Window{}
}

type Window struct {
	EventTarget
}

var _ hypp.Window = Window{}

func (w Window) EscapeToValue() hypp.Value {
	return nil
}

func (w Window) RequestAnimationFrame(f func()) int {
	f()
	return 1
}

type EventTarget struct{}

func (e EventTarget) RemoveEventListener(kind string, listenerID hypp.EventListenerID) {}

func (e EventTarget) AddEventListener(kind string, listenerID hypp.EventListener) hypp.EventListenerID {
	return EventListenerID{}
}

type EventListenerID struct{}

var _ hypp.EventListenerID = EventListenerID{}

func (e EventListenerID) IAmAnEventListenerID() {}

type Node struct {
	EventTarget
	parentNode   *Node
	nodeType     int
	nodeValue    string
	nodeName     string
	namespaceURI string
	childNodes   []hypp.Node
	attributes   hypp.Map[string, string]
}

var _ hypp.Node = &Node{}

// See https://w3c.github.io/html-reference/syntax.html#void-element
var voidElements = hypp.NewSet[string](
	"area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr",
)

type RenderOptions struct {
	// Deterministic, when true, will ensure the rendered HTML will always be the same.
	// It will sort the attributes by their keys.
	Deterministic bool
}

// OuterHTML renders the Node to an HTML string.
// The HTML will include the Node's tag.
// The given options can be nil.
func (n Node) OuterHTML(options *RenderOptions) string {
	if n.nodeType == textNode {
		return n.nodeValue
	}
	open := "<" + n.nodeName
	if options != nil && options.Deterministic {
		keys := make([]string, len(n.attributes))
		i := 0
		for k := range n.attributes {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		for _, k := range keys {
			open += fmt.Sprintf(" %s=%q", k, n.attributes[k])
		}
	} else {
		for k, v := range n.attributes {
			open += fmt.Sprintf(" %s=%q", k, v)
		}
	}
	open += ">"
	if len(n.childNodes) > 0 {
		s := open
		for _, child := range n.childNodes {
			s += child.(*Node).OuterHTML(options)
		}
		s += "</" + n.nodeName + ">"
		return s
	} else if voidElements.Has(n.nodeName) {
		return open
	} else {
		return open + "</" + n.nodeName + ">"
	}
}

// InnerHTML renders the child nodes to an HTML string.
// The HTML does not include the Node's tag.
// The given options can be nil.
func (n Node) InnerHTML(options *RenderOptions) string {
	s := ""
	for _, child := range n.childNodes {
		s += child.(*Node).OuterHTML(options)
	}
	return s
}

func (n Node) ParentNode() hypp.Node {
	return n.parentNode
}

func (n Node) NodeType() int {
	return n.nodeType
}

func (n Node) NodeValue() string {
	return n.nodeValue
}

func (n *Node) SetNodeValue(nodeValue string) {
	n.nodeValue = nodeValue
}

func (n Node) NodeName() string {
	return n.nodeName
}

func (n Node) ChildNodes() []hypp.Node {
	return n.childNodes
}

func (n *Node) InsertBefore(newNode, referenceNode hypp.Node) hypp.Node {
	parentNode := newNode.ParentNode().(*Node)
	if parentNode != nil {
		parentNode.RemoveChild(newNode)
	}
	if referenceNode == nil {
		newNode.(*Node).parentNode = n
		return n.AppendChild(newNode)
	} else {
		for i, child := range n.ChildNodes() {
			if child == referenceNode {
				newNode.(*Node).parentNode = n
				n.childNodes = append(n.childNodes[:i+1], n.childNodes[i:]...)
				n.childNodes[i] = newNode
				return newNode
			}
		}
		panic(errors.New("html: referenceNode is not a child of this Node"))
	}
}

func (n *Node) RemoveChild(child hypp.Node) {
	for i, c := range n.childNodes {
		if c == child {
			child.(*Node).parentNode = nil
			n.childNodes = append(n.childNodes[:i], n.childNodes[i+1:]...)
		}
	}
	panic(errors.New("html: cannot remove Node that is not a child of this Node"))
}

func (n Node) Get(name string) hypp.Option[interface{}] {
	if n.attributes == nil {
		return hypp.Option[interface{}]{}
	}
	name = camelToKebab(name)
	v, ok := n.attributes[name]
	return hypp.Option[interface{}]{V: v, OK: ok}
}

func (n Node) In(name string) bool {
	return false
}

func (n *Node) Set(name string, value interface{}) {
	name = camelToKebab(name)
	n.attributes.Set(name, fmt.Sprint(value))
}

func (n *Node) AppendChild(child hypp.Node) hypp.Node {
	n.childNodes = append(n.childNodes, child)
	return child
}

func (n *Node) RemoveAttribute(name string) {
	delete(n.attributes, name)
}

func (n *Node) SetAttribute(name string, value interface{}) {
	n.attributes.Set(name, fmt.Sprint(value))
}

func (n Node) Events() hypp.Events {
	return &Events{}
}

func (n Node) Style() hypp.Style {
	// TODO
	return Style{}
}

func (n Node) EventListenerID(kind string) hypp.EventListenerID {
	return EventListenerID{}
}

func (n *Node) SetEventListenerID(kind string, eventListenerID hypp.EventListenerID) {}

type Events struct{}

var _ hypp.Events = &Events{}

func (e *Events) Set(name string, value hypp.Dispatchable) {}

func (e Events) Get(name string) hypp.Dispatchable {
	return nil
}

func (e *Events) Del(name string) {}

type EscapeToValuer struct{}

var _ hypp.EscapeToValuer = EscapeToValuer{}

func (e EscapeToValuer) EscapeToValue() hypp.Value {
	return nil
}

type Event struct {
	EscapeToValuer
}

var _ hypp.Event = Event{}

func (e Event) Type() string {
	return ""
}

func (e Event) PreventDefault() {}

func (e Event) Target() hypp.EventTargetValuer {
	return EventTargetValuer{}
}

type EventTargetValuer struct{}

var _ hypp.EventTargetValuer = EventTargetValuer{}

func (e EventTargetValuer) Value() string {
	return ""
}

type Style map[string]string

var _ hypp.Style = Style{}

func (s Style) SetProperty(propertyName, value string) {
	s[propertyName] = value
}

var matchCamelCase = regexp.MustCompile("([a-z0-9])([A-Z])")

func camelToKebab(s string) string {
	return matchCamelCase.ReplaceAllString(s, "${1}-${2}")
}

func (s Style) Set(name, value string) {
	s[camelToKebab(name)] = value
}

func (s Style) Get(name string) string {
	return s[camelToKebab(name)]
}

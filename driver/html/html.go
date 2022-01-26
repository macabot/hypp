package html

import (
	"errors"
	"regexp"

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
	}
}

func (d Driver) CreateElementNS(namespaceURI, qualifiedName string, options hypp.Option[hypp.ElementCreationOptions]) hypp.Node {
	return &Node{
		nodeType: ssrNode,
		nodeName: qualifiedName,
	}
}

func (d Driver) CreateElement(tagName string, options hypp.Option[hypp.ElementCreationOptions]) hypp.Node {
	return &Node{
		nodeType: ssrNode,
		nodeName: tagName,
	}
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
	parentNode *Node
	nodeType   int
	nodeValue  string
	nodeName   string
	childNodes []hypp.Node
}

var _ hypp.Node = &Node{}

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
	if newNode.ParentNode() != nil {
		newNode.ParentNode().RemoveChild(newNode)
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
	// TODO
	return hypp.Option[interface{}]{}
}

func (n Node) In(name string) bool {
	// TODO
	return false
}

func (n *Node) Set(name string, value interface{}) {
	// TODO
}

func (n *Node) AppendChild(child hypp.Node) hypp.Node {
	n.childNodes = append(n.childNodes, child)
	return child
}

func (n *Node) RemoveAttribute(name string) {
	// TODO
}

func (n *Node) SetAttribute(name string, value interface{}) {
	// TODO
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
	return matchCamelCase.ReplaceAllString(s, "${1}_${2}")
}

func (s Style) Set(name, value string) {
	s[camelToKebab(name)] = value
}

func (s Style) Get(name string) string {
	return s[camelToKebab(name)]
}

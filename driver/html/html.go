package html

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/macabot/hypp"
	"golang.org/x/net/html"
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

func (d Driver) ValidateAppPropsNode(node hypp.Node) error {
	return nil
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
	style        hypp.Map[string, string]
}

var _ hypp.Node = &Node{}

// See https://w3c.github.io/html-reference/syntax.html#void-element
var voidElements = hypp.NewSet(
	"area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr",
)

type RenderOptions struct {
	// Deterministic, when true, will ensure the rendered HTML will always be the same.
	// It will sort map values by their keys, such as attributes and style properties.
	Deterministic bool
}

func (r *RenderOptions) isDeterministic() bool {
	return r != nil && r.Deterministic
}

// OuterHTML renders the Node to an HTML string.
// The HTML will include the Node's tag.
// The given options can be nil.
func (n Node) OuterHTML(options *RenderOptions) string {
	if n.nodeType == textNode {
		return n.nodeValue
	}
	attributes := n.attributes.Copy()
	if n.style != nil {
		if attributes == nil {
			attributes = hypp.Map[string, string]{}
		}
		attributes["style"] = renderStyle(n.style, options)
	}
	if attributes.Has("class") {
		attributes["class"] = renderClass(attributes["class"], options)
	}
	open := "<" + n.nodeName
	if options.isDeterministic() {
		keys := make([]string, len(attributes))
		i := 0
		for k := range attributes {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		for _, k := range keys {
			open += fmt.Sprintf(" %s=%q", k, attributes[k])
		}
	} else {
		for k, v := range attributes {
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
	if referenceNode == nil {
		return n.AppendChild(newNode)
	}
	parentNode := newNode.ParentNode().(*Node)
	if parentNode != nil {
		parentNode.RemoveChild(newNode)
	}
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

func (n *Node) RemoveChild(child hypp.Node) {
	for i, c := range n.childNodes {
		if c == child {
			child.(*Node).parentNode = nil
			n.childNodes = append(n.childNodes[:i], n.childNodes[i+1:]...)
			return
		}
	}
	panic(errors.New("html: cannot remove Node that is not a child of this Node"))
}

func (n Node) Get(name string) hypp.Option[interface{}] {
	o := n.attributes.GetOption(camelToKebab(name))
	return hypp.Option[interface{}]{V: o.V, OK: o.OK}
}

func (n Node) In(name string) bool {
	return false
}

func (n *Node) Set(name string, value interface{}) {
	n.SetAttribute(camelToKebab(name), value)
}

func (n *Node) AppendChild(child hypp.Node) hypp.Node {
	parentNode := child.ParentNode().(*Node)
	if parentNode != nil {
		parentNode.RemoveChild(child)
	}
	c := child.(*Node)
	c.parentNode = n
	n.childNodes = append(n.childNodes, c)
	return child
}

func (n *Node) RemoveAttribute(name string) {
	delete(n.attributes, name)
}

var matchUnescapedDoubleQuote = regexp.MustCompile(`\\([\s\S])|(")`)

func (n *Node) SetAttribute(name string, value interface{}) {
	if name == "style" {
		if m, ok := value.(map[string]string); ok {
			n.style = hypp.Map[string, string]{}
			for k, v := range m {
				n.SetStyle(k, v)
			}
			return
		}
	}
	if n.attributes == nil {
		n.attributes = hypp.Map[string, string]{}
	}
	s := fmt.Sprint(value)
	s = matchUnescapedDoubleQuote.ReplaceAllString(s, `$1$2`)
	n.attributes[name] = s
}

func (n Node) Events() hypp.Events {
	return &Events{}
}

func (n *Node) SetStyleProperty(propertyName, value string) {
	if n.style == nil {
		n.style = hypp.Map[string, string]{}
	}
	if value == "" {
		delete(n.style, propertyName)
	} else {
		n.style[propertyName] = html.EscapeString(value)
	}
}

func (n *Node) SetStyle(name, value string) {
	n.SetStyleProperty(camelToKebab(name), value)
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

func renderStyle(style hypp.Map[string, string], options *RenderOptions) string {
	parts := make([]string, len(style))
	i := 0
	for key, value := range style {
		parts[i] = key + ": " + value + ";"
		i++
	}
	if options.isDeterministic() {
		sort.Strings(parts)
	}
	return strings.Join(parts, " ")
}

var matchSpaces = regexp.MustCompile(`\s+`)

func renderClass(class string, options *RenderOptions) string {
	if !options.isDeterministic() {
		return class
	}
	parts := matchSpaces.Split(class, -1)
	sort.Strings(parts)
	return strings.Join(parts, " ")
}

var matchCamelCase = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func camelToKebab(s string) string {
	return strings.ToLower(matchCamelCase.ReplaceAllString(s, "${1}-${2}"))
}

func nodeToVNode(node *html.Node) *hypp.VNode {
	switch node.Type {
	case html.TextNode:
		return hypp.Text(node.Data)
	case html.ElementNode:
		hProps := hypp.HProps{}
		for _, attribute := range node.Attr {
			hProps[attribute.Key] = attribute.Val
		}
		var children []*hypp.VNode
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			children = append(children, nodeToVNode(child))
		}
		return hypp.H(node.Data, hProps, children...)
	default:
		return nil
	}
}

func Parse(r io.Reader) (*hypp.VNode, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return nodeToVNode(node), nil
}

func ParseFragment(r io.Reader, context *html.Node) ([]*hypp.VNode, error) {
	nodes, err := html.ParseFragment(r, context)
	if err != nil {
		return nil, err
	}
	vNodes := make([]*hypp.VNode, len(nodes))
	for i, node := range nodes {
		vNodes[i] = nodeToVNode(node)
	}
	return vNodes, nil
}

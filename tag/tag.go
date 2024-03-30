package tag

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/macabot/hypp"
	"golang.org/x/exp/maps"
	"golang.org/x/net/html"
)

var (
	ErrNilWriter = errors.New("writer cannot be nil")
	ErrNilReader = errors.New("reader cannot be nil")
)

// Render renders the node to the given writer.
// It is based on https://pkg.go.dev/golang.org/x/net/html#Render
func Render(w io.Writer, node *hypp.VNode) error {
	if w == nil {
		return ErrNilWriter
	}
	if node == nil {
		return nil
	}
	element, err := vNodeToNode(node, nil)
	if err != nil {
		return err
	}
	return html.Render(w, element)
}

// RenderToString renders the node to a string.
func RenderToString(node *hypp.VNode) (string, error) {
	if node == nil {
		return "", nil
	}
	w := &bytes.Buffer{}
	if err := Render(w, node); err != nil {
		return "", err
	}
	return w.String(), nil
}

// RenderFragment renders a fragment of nodes to the given writer.
func RenderFragment(w io.Writer, fragment []*hypp.VNode) error {
	if w == nil {
		return ErrNilWriter
	}
	for _, node := range fragment {
		if err := Render(w, node); err != nil {
			return err
		}
	}
	return nil
}

// RenderFragmentToString renders a fragment of nodes to a string.
func RenderFragmentToString(fragment []*hypp.VNode) (string, error) {
	w := &bytes.Buffer{}
	if err := RenderFragment(w, fragment); err != nil {
		return "", err
	}
	return w.String(), nil
}

// Parse returns the VNode tree for the HTML from the given Reader.
// It is based on https://pkg.go.dev/golang.org/x/net/html#Parse
func Parse(r io.Reader) (*hypp.VNode, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return nodeToVNode(node)
}

// ParseFromString is the same as [Parse], but takes as argument a string.
func ParseFromString(s string) (*hypp.VNode, error) {
	r := strings.NewReader(s)
	return Parse(r)
}

// ParseFragment parses a fragment of HTML and returns the VNodes that were found.
// If the fragment is the InnerHTML for an existing element, pass that element as context.
// It is based on https://pkg.go.dev/golang.org/x/net/html#ParseFragment
func ParseFragment(r io.Reader, context *html.Node) ([]*hypp.VNode, error) {
	if r == nil {
		return nil, ErrNilReader
	}
	nodes, err := html.ParseFragment(r, context)
	if err != nil {
		return nil, err
	}
	vNodes := make([]*hypp.VNode, len(nodes))
	for i, node := range nodes {
		childNode, err := nodeToVNode(node)
		if err != nil {
			return nil, err
		}
		vNodes[i] = childNode
	}
	return vNodes, nil
}

// ParseFragmentFromString is the same as [ParseFragment], but takes as argument a string.
func ParseFragmentFromString(s string, context *html.Node) ([]*hypp.VNode, error) {
	r := strings.NewReader(s)
	return ParseFragment(r, context)
}

// nodeToVNode converts a *html.Node to a *hypp.VNode.
func nodeToVNode(node *html.Node) (*hypp.VNode, error) {
	switch node.Type {
	case html.TextNode:
		return hypp.Text(node.Data), nil
	case html.ElementNode:
		hProps := hypp.HProps{}
		for _, attribute := range node.Attr {
			hProps[attribute.Key] = attribute.Val
		}
		var children []*hypp.VNode
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			childNode, err := nodeToVNode(child)
			if err != nil {
				return nil, err
			}
			children = append(children, childNode)
		}
		return hypp.H(node.Data, hProps, children...), nil
	default:
		return nil, fmt.Errorf("hypp/tag: cannot parse NodeType %v", node.Type)
	}
}

func renderStyle(style map[string]string) string {
	keys := maps.Keys(style)
	sort.Strings(keys)
	parts := make([]string, len(style))
	for i, key := range keys {
		parts[i] = key + ": " + style[key] + ";"
	}
	return strings.Join(parts, " ")
}

func propsToAttributes(props hypp.HProps) []html.Attribute {
	keys := maps.Keys(props)
	sort.Strings(keys)
	var attributes []html.Attribute
	for _, key := range keys {
		if key == "key" {
			continue
		}
		if len(key) >= 2 && key[0] == 'o' && key[1] == 'n' {
			continue
		}
		value := props[key]
		if key == "style" {
			value = renderStyle(value.(map[string]string))
		}
		attributes = append(attributes, html.Attribute{
			Key: key,
			Val: fmt.Sprint(value),
		})
	}
	return attributes
}

func vNodeToNode(node *hypp.VNode, parent *html.Node) (*html.Node, error) {
	switch node.Kind() {
	case hypp.TextNode:
		return &html.Node{
			Type: html.TextNode,
			Data: node.Tag(),
		}, nil
	case hypp.SSRNode:
		element := &html.Node{
			Parent: parent,
			Type:   html.ElementNode,
			Data:   node.Tag(),
			Attr:   propsToAttributes(node.Props()),
		}
		var firstChild *html.Node
		var lastChild *html.Node
		children := node.Children()
		var childElements []*html.Node
		for _, child := range children {
			if child == nil {
				continue
			}
			childElement, err := vNodeToNode(child, parent)
			if err != nil {
				return nil, err
			}
			childElements = append(childElements, childElement)
			if firstChild == nil {
				firstChild = childElement
			}
			lastChild = childElement
		}
		for i := range childElements {
			var prevSibling *html.Node
			if i > 0 {
				prevSibling = childElements[i-1]
			}
			var nextSibling *html.Node
			if i < len(childElements)-1 {
				nextSibling = childElements[i+1]
			}
			childElements[i].PrevSibling = prevSibling
			childElements[i].NextSibling = nextSibling
		}
		element.FirstChild = firstChild
		element.LastChild = lastChild
		return element, nil
	default:
		return nil, fmt.Errorf("hypp/tag: cannot render node kind %v", node.Kind())
	}
}

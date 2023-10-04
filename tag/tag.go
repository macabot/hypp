package tag

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/util"
	"golang.org/x/exp/maps"
	"golang.org/x/net/html"
)

// See https://w3c.github.io/html-reference/syntax.html#void-element
var voidElements = util.NewSet(
	"area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr",
)

// TODO what about SVG void elements?

func renderStyle(style map[string]string) string {
	keys := maps.Keys(style)
	sort.Strings(keys)
	parts := make([]string, len(style))
	for i, key := range keys {
		// TODO should html.EscapeString be called here or be called for all property values?
		parts[i] = key + ": " + html.EscapeString(style[key]) + ";"
	}
	return strings.Join(parts, " ")
}

// TODO should the Render function use https://pkg.go.dev/golang.org/x/net/html#Render ?

func Render(w io.Writer, node *hypp.VNode) error {
	if node.Kind() == hypp.TextNode {
		// TODO should the value be passed through html.EscapeString?
		_, err := w.Write([]byte(node.Tag()))
		return err
	}

	if err := hypp.ValidateHProps(node.Props(), node.Tag()); err != nil {
		return err
	}

	if _, err := w.Write([]byte("<" + node.Tag())); err != nil {
		return err
	}

	props := node.Props()
	keys := maps.Keys(props)
	sort.Strings(keys)
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
		if _, err := fmt.Fprintf(w, " %s=%q", key, props[key]); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(">")); err != nil {
		return err
	}

	if len(node.Children()) == 0 && voidElements.Has(node.Tag()) {
		return nil
	}

	for _, child := range node.Children() {
		Render(w, child)
	}

	_, err := w.Write([]byte("</" + node.Tag() + ">"))
	return err
}

func RenderToString(node *hypp.VNode) (string, error) {
	w := &bytes.Buffer{}
	if err := Render(w, node); err != nil {
		return "", err
	}
	return w.String(), nil
}

func RenderFragment(w io.Writer, fragment []*hypp.VNode) error {
	for _, node := range fragment {
		if err := Render(w, node); err != nil {
			return err
		}
	}
	return nil
}

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
		return nil, fmt.Errorf("hypp/tag: could not parse NodeType %v", node.Type)
	}
}

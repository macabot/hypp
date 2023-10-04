package hypp

import (
	"fmt"
	"io"
	"sort"

	"github.com/macabot/hypp/util"
)

// See https://w3c.github.io/html-reference/syntax.html#void-element
var voidElements = util.NewSet(
	"area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr",
)

func Render(w io.Writer, nodes ...*VNode) error {
	for _, node := range nodes {
		if _, err := renderVNode(w, node); err != nil {
			return err
		}
	}
	return nil
}

func renderVNode(w io.Writer, node *VNode) error {
	if node.kind == TextNode {
		_, err := w.Write([]byte(node.tag))
		return err
	}

	if _, err := w.Write([]byte("<" + node.tag)); err != nil {
		return err
	}

	keys := make([]string, len(node.props))
	i := 0
	for key := range node.props {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		if _, err := fmt.Fprintf(w, " %s=%q", key, node.props[key]); err != nil {
			return err
		}
	}

	if len(node.children) == 0 && voidElements.Has(node.tag) {
		if _, err := w.Write([]byte("/")); err != nil {
			return err
		}
	}

	// TODO continue

	// if len(node.children) > 0 {
	// 	for _, child := range node.children {
	// 		renderVNode(w, child)
	// 	}
	// }

	return nil
}

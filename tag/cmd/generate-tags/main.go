package main

import (
	"fmt"
	"os"
	"strings"
)

var fileTemplate = `package %s

// DO NOT EDIT
// This file was generated using github.com/macabot/hypp/tag/cmd/generate-tags

import "github.com/macabot/hypp"

%s`

var funcTemplate = `func %s(props hypp.HProps, children ...*hypp.VNode) *hypp.VNode {
	return hypp.H("%s", props, children...)
}
`

// main generates the tag functions.
// The first argument is the package name.
// All following arguments are the tag names. Each tag name must match `[a-z]+`.
func main() {
	packageName := os.Args[1]
	funcs := make([]string, len(os.Args)-2)
	for i, arg := range os.Args[2:] {
		funcs[i] = fmt.Sprintf(funcTemplate, title(arg), arg)
	}

	fmt.Printf(
		fileTemplate,
		packageName,
		strings.Join(funcs, "\n"),
	)
}

func title(s string) string {
	const diff = 'a' - 'A'
	first := s[0]
	if first < 'a' || first > 'z' {
		panic("invalid tag name: " + s)
	}
	return string(first-diff) + string(s[1:])
}

package main

import (
	"fmt"
	"os"
	"strings"
)

var fileTemplate = `package %s

// DO NOT EDIT
// This file was generated using cmd/generate-html-tags

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag"
)

var Tags = tag.NewSet(
%s)

%s`

var funcTemplate = `func %s(props hypp.HProps, children ...*hypp.VNode) *hypp.VNode {
	return hypp.H("%s", props, children...)
}
`

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	packageName := os.Args[1]
	tags := make([]string, len(os.Args)-2)
	funcs := make([]string, len(os.Args)-2)
	for i, arg := range os.Args[2:] {
		tags[i] = "\t\"" + arg + "\",\n"
		funcs[i] = fmt.Sprintf(funcTemplate, strings.Title(arg), arg)
	}

	fmt.Printf(
		fileTemplate,
		packageName,
		strings.Join(tags, ""),
		strings.Join(funcs, "\n"),
	)
}

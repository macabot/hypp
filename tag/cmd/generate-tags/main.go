package main

import (
	"fmt"
	"os"
	"strings"
)

var fileTemplate = `// DO NOT EDIT
// This file was generated using cmd/generate-html-tags
package %s

import "github.com/macabot/hypp"

%s
`

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
	funcs := make([]string, len(os.Args)-2)
	for i, arg := range os.Args[2:] {
		funcs[i] = fmt.Sprintf(funcTemplate, strings.Title(arg), arg)
	}

	fmt.Printf(fileTemplate, packageName, strings.Join(funcs, "\n"))
}

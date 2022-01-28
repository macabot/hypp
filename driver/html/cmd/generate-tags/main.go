package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var fileTemplate = `package %s

// DO NOT EDIT
// This file was generated using github.com/macabot/hypp/driver/html/cmd/generate-tags

import "github.com/macabot/hypp"

var hierarchy = map[string]string{
%s}

var ownProps = map[string]hypp.Set[string]{
%s}`

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Data struct {
	Hierarchy map[string]string   `json:"hierarchy"`
	OwnProps  map[string][]string `json:"ownProps"`
}

func main() {
	packageName := os.Args[1]
	filename := os.Args[2]

	b, err := ioutil.ReadFile(filename)
	panicIf(err)
	data := Data{}
	panicIf(json.Unmarshal(b, &data))

	h := make([]string, len(data.Hierarchy))
	i := 0
	for key, value := range data.Hierarchy {
		h[i] = "\t\"" + key + "\": \"" + value + "\",\n"
		i++
	}

	o := make([]string, len(data.OwnProps))
	i = 0
	for key, value := range data.OwnProps {
		o[i] = "\t\"" + key + "\": hypp.NewSet(\"" + strings.Join(value, "\", \"") + "\"),\n"
		i++
	}

	fmt.Printf(
		fileTemplate,
		packageName,
		strings.Join(h, ""),
		strings.Join(o, ""),
	)
}

// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/tests/index.test.js
package hypp

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateVirtualNode(t *testing.T) {
    assert.Equal(t, VNode{
        props: HProps{"foo": true},
        tag: "zord",
        kind: 1,
    }, H("zord", HProps{"foo": true}))
}

func TestCreateTextNode(t *testing.T) {
    assert.Equal(t, VNode{
        tag: "hyper",
        kind: 3,
    }, Text("hyper"))
}

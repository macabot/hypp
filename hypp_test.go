package hypp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// This test is based on https://github.com/jorgebucaran/hyperapp/blob/8e6a4908954186ed8466ae9ef15ec172ddbea2d0/tests/index.test.js#L7
func TestCreateVirtualNode(t *testing.T) {
	assert.Equal(t, &VNode{
		props: HProps{"foo": true},
		tag:   "zord",
		kind:  1,
	}, H("zord", HProps{"foo": true}))
}

// This test is based on https://github.com/jorgebucaran/hyperapp/blob/8e6a4908954186ed8466ae9ef15ec172ddbea2d0/tests/index.test.js#L20
func TestCreateTextNode(t *testing.T) {
	assert.Equal(t, &VNode{
		tag:  "hyper",
		kind: 3,
	}, Text("hyper"))
}

func TestSetValueOnNilHProps(t *testing.T) {
	var props HProps
	props.Set("foo", "bar")
	assert.Equal(t, HProps{"foo": "bar"}, props)
}

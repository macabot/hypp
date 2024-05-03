package hypp

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

// TestKeepKeyOnVNodeShallowClone should ensure that the key of a VNode is kept when performing a shallow clone.
//
// The following hyperapp pull request contains "a small improvement regarding how keys are handled": https://github.com/jorgebucaran/hyperapp/pull/1090
// This change removes "key" from the HProps when creating a VNode.
// Similar changes were made to hypp: https://github.com/macabot/hypp/pull/26
// The problem, however, was that a shallow clone of a VNode no longer contained the key.
func TestKeepKeyOnVNodeShallowClone(t *testing.T) {
	span := H("span", HProps{"key": "foo"}, Text("test"))

	shallowClone := func(n *VNode) *VNode {
		return H(n.Tag(), n.Props(), n.Children()...)
	}

	spanClone := shallowClone(span)
	assert.Equal(
		t,
		option[string]{V: "foo", OK: true},
		spanClone.Props().key(),
	)
}

func TestDontPanicOnShortPropertyKey(t *testing.T) {
	assert.NoError(t, ValidateHProps(HProps{"": "y"}))
	assert.NoError(t, ValidateHProps(HProps{"o": "y"}))
}

func TestValidateHPropsReturnsErrorIfOtherHasInvalidType(t *testing.T) {
	assert.Error(t, ValidateHProps(HProps{"x": []string{"foo"}}))
}

package tag_test

import (
	"testing"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FIXME tests are not passing

func TestRender(t *testing.T) {
	node := hypp.H(
		"a",
		hypp.HProps{
			"attr-a1": 33,
			"attr-a2": "a2",
		},
		hypp.H(
			"p",
			hypp.HProps{
				"attr-b1": "hello",
				"attr-b2": `"world"`,
			},
			hypp.Text("This is a test"),
			hypp.H("br", nil),
		),
	)
	output, err := tag.RenderToString(node)
	require.NoError(t, err)

	assert.Equal(
		t,
		`<a attr-a1="33" attr-a2="a2"><p attr-b1="hello" attr-b2="&#34;world&#34;">This is a test<br/></p></a>`,
		output,
	)
}

func TestRenderWithStyle(t *testing.T) {
	node := hypp.H(
		"div",
		hypp.HProps{
			"style": map[string]string{
				"background-color": "red",
				"background-image": `url("hypp.png")`,
			},
		},
	)
	output, err := tag.RenderToString(node)
	require.NoError(t, err)

	assert.Equal(
		t,
		`<div style="background-color: red; background-image: url(&#34;hypp.png&#34;);"></div>`,
		output,
	)
}

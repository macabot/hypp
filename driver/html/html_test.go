package html

import (
	"testing"

	"github.com/macabot/hypp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var eco = hypp.Option[hypp.ElementCreationOptions]{}

func TestInnerHTML(t *testing.T) {
	driver := Driver{}

	wrapper := Node{}

	a := driver.CreateElement("a", eco)
	a.Set("attrA1", 33)
	a.Set("attrA2", "a2")

	b := driver.CreateElement("b", eco)
	b.Set("attrB1", "hello")
	b.Set("attrB2", `"world"`)

	c := driver.CreateTextNode("This is a test")

	d := driver.CreateElement("br", eco)

	wrapper.AppendChild(a)
	a.AppendChild(b)
	b.AppendChild(c)
	b.AppendChild(d)

	assert.Equal(
		t,
		`<a attr-a1="33" attr-a2="a2"><b attr-b1="hello" attr-b2="\"world\"">This is a test<br></b></a>`,
		wrapper.InnerHTML(&RenderOptions{Deterministic: true}),
	)
}

func TestRenderStyle(t *testing.T) {
	driver := Driver{}
	div := driver.CreateElement("div", eco).(*Node)
	div.SetAttribute("style", map[string]string{
		"backgroundColor": "red",
	})
	div.SetStyle("backgroundImage", `url("hypp.png")`)

	assert.Equal(
		t,
		`<div style="background-color: red; background-image: url(&#34;hypp.png&#34;);"></div>`,
		div.OuterHTML(&RenderOptions{Deterministic: true}),
	)
}

func TestRenderClassSlice(t *testing.T) {
	driver := Driver{}
	div := driver.CreateElement("div", eco).(*Node)
	div.SetAttribute("class", "b c a")

	assert.Equal(
		t,
		`<div class="a b c"></div>`,
		div.OuterHTML(&RenderOptions{Deterministic: true}),
	)
}

func TestAppendChild(t *testing.T) {
	driver := Driver{}
	a := driver.CreateElement("a", eco)
	b := driver.CreateElement("b", eco)
	c := driver.CreateElement("c", eco)
	d := driver.CreateElement("d", eco)
	out := a.AppendChild(b)
	assert.Equal(t, b, out)

	d.AppendChild(c)
	assert.Equal(t, d, c.ParentNode())

	a.AppendChild(c)

	assert.Equal(t, a, b.ParentNode())
	assert.Equal(t, a, c.ParentNode())
	require.Len(t, a.ChildNodes(), 2)
	assert.Equal(t, b, a.ChildNodes()[0])
	assert.Equal(t, c, a.ChildNodes()[1])

	a.AppendChild(b)

	require.Len(t, a.ChildNodes(), 2)
	assert.Equal(t, c, a.ChildNodes()[0])
	assert.Equal(t, b, a.ChildNodes()[1])
}

func TestInsertBeforeChangesParent(t *testing.T) {
	driver := Driver{}
	a := driver.CreateElement("a", eco)
	b := driver.CreateElement("b", eco)
	a.AppendChild(b)
	assert.Equal(t, a, b.ParentNode())
	assert.Len(t, a.ChildNodes(), 1)

	c := driver.CreateElement("c", eco)
	out := c.InsertBefore(b, nil)
	assert.Equal(t, b, out)
	assert.Equal(t, c, b.ParentNode())
	assert.Len(t, c.ChildNodes(), 1)
	assert.Len(t, a.ChildNodes(), 0)
}

func TestInsertBeforeNode(t *testing.T) {
	driver := Driver{}
	a := driver.CreateElement("a", eco)
	b := driver.CreateElement("b", eco)
	c := driver.CreateElement("c", eco)

	a.AppendChild(b)
	require.Len(t, a.ChildNodes(), 1)
	assert.Equal(t, b, a.ChildNodes()[0])

	a.InsertBefore(c, b)
	require.Len(t, a.ChildNodes(), 2)
	assert.Equal(t, c, a.ChildNodes()[0])
	assert.Equal(t, b, a.ChildNodes()[1])
}

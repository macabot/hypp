package html

import (
    "testing"

    "github.com/macabot/hypp"
    "github.com/stretchr/testify/assert"
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
    b.Set("attrB2", "world")

    c := driver.CreateTextNode("This is a test")

    d := driver.CreateElement("br", eco)

    wrapper.AppendChild(a)
    a.AppendChild(b)
    b.AppendChild(c)
    b.AppendChild(d)

    assert.Equal(
        t,
        `<a attr-a1="33" attr-a2="a2"><b attr-b1="hello" attr-b2="world">This is a test<br></b></a>`,
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

// This file is based on https://codepen.io/jorgebucaran/pen/bgWBdV
package app

import (
    "github.com/macabot/hypp"
    "github.com/macabot/hypp/tag/html"
)

type MyState struct {
    hypp.EmptyState
    content string
    count int
}

func (m MyState) clone() *MyState {
    return &m
}

func button[S hypp.State](onclick hypp.Action[S], text string, props hypp.HProps) *hypp.VNode {
    props.Set("onclick", onclick)
    return html.Button(props, hypp.Text(text))
}

func title(text string) *hypp.VNode {
    return html.H1(nil, hypp.Text(text))
}

var maxLength = 140

type contentAndLength struct {
    content string
    length int
}

func setText(state *MyState, payload hypp.Payload) hypp.Dispatchable {
    cl := payload.(contentAndLength)
    newState := state.clone()
    if cl.length > maxLength {
        newState.count = 0
    } else {
        newState.content = cl.content
        newState.count = state.count + len(state.content) - len(cl.content)
    }
    return newState
}

func tweet(_ *MyState, _ hypp.Payload) hypp.Dispatchable {
    return &MyState{count: maxLength}
}

func Run(driver hypp.Driver, node hypp.Node) {
    hypp.App[*MyState](hypp.AppProps[*MyState]{
        Driver: driver,
        Init: &MyState{count: maxLength},
        View: func(state *MyState) *hypp.VNode {
            var oninput hypp.Action[*MyState] = func(_ *MyState, payload hypp.Payload) hypp.Dispatchable {
                event := payload.(hypp.Event)
                content := event.Target().Value()
                return hypp.ActionAndPayload[*MyState]{
                    Action: setText,
                    Payload: contentAndLength{
                        content: content,
                        length: len(content),
                    },
                }
            }
            return html.Main(
                nil,
                title("Tweeter 🦤"),
                html.Textarea(hypp.HProps{
                    "placeholder": "What's on your mind?",
                    "oninput": oninput,
                    "value": state.content,
                    "rows": 4 + len(state.content) / 100,
                }),
                html.Section(
                    nil,
                    hypp.Textf("%d", state.count),
                    button(tweet, "Tweet", hypp.HProps{
                        "disabled": state.count >= maxLength,
                    }),
                ),
            )
        },
        Node: node,
    })
}

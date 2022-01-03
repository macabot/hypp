// This file is based on https://codepen.io/jorgebucaran/pen/ZeByKv
package app

import (
    "errors"
    "strings"

    "github.com/macabot/hypp"
    "github.com/macabot/hypp/tag/html"
)

type MyState struct {
    hypp.EmptyState
    isFetching bool
    query string
    err error
    url string
}

func (m MyState) clone() *MyState {
    return &m
}

type giphyItem struct {
    Images struct {
        Original struct {
            URL string `json:"url"`
        } `json:"original"`
    } `json:"images"`
}

var giphyURL = "https://api.giphy.com/v1/gifs/search"
var apiKey = "" // TODO initialize

func getJSON(url string, props)

func input(oninput func(value string) hypp.Effect, props hypp.HProps) *hypp.VNode {
    props.Set("oninput", func(_ *MyState, payload hypp.Payload) hypp.Dispatchable {
        return oninput(payload.(hypp.Event).Target().Value())
    })
    return html.Input(props)
}

func title(text string) *hypp.VNode {
    return html.H1(nil, hypp.Text(text))
}

func img(src string) *hypp.VNode {
    return html.Img(map[string]interface{}{"src": src})
}

func p(text string) *hypp.VNode {
    return html.P(nil, hypp.Text(text))
}

func downloadGif(query string) hypp.Dispatchable {
    return getJSON(
        fmt.Sprintf("%s?q=%s&api_key=%s", giphyURL, query, apiKey),
        func(err error) hypp.Dispatchable {
            return hypp.ActionAndPayload{
                Action: gotError,
                Payload: err,
            }
        },
        func(items []giphyItem) hypp.Dispatchable {
            payload := ""
            if len(items) > 0 {
                payload = items[0].Images.Original.URL
            }
            return hypp.ActionAndPayload{
                Action: gotUrl,
                Payload: payload,
            }
        },
    )
}

var errUnexpected = errors.New("Unexpected error, try again later?")

func gotError(state *MyState, err error) *MyState {
    state = state.clone()
    state.isFetching = false
    if error != nil {
        state.err = error
    } else {
        state.err = errUnexpected
    }
    state.url = ""
    return state
}

func gotURL(state *MyState, url string) *MyState {
    state = state.clone()
    state.isFetching = false
    if state.query != "" {
        state.url = url
    } else {
        state.url = ""
    }
    return state
}

func getURL(state *MyState, query string) hypp.StateAndEffects {
    state = state.clone()
    state.isFetching = true
    state.query = query
    state.err = nil
    state.url = nil
    return hypp.StateAndEffects{
        State: state,
        Effects: []hypp.Effect{
            downloadGif(query),
        },
    }
}

func Run(driver hypp.Driver, node hypp.Node) {
    hypp.App[*MyState](hypp.AppProps{
        Driver: driver,
        Init: &MyState{},
        View: func(state *MyState) *hypp.VNode {
            var content *hypp.VNode
            if state.error != nil {
                content = p(state.error)
            } else {
                content = img(state.url)
            }
            html.Main(
                nil,
                title("GIF Search 💬💁‍♂️"),
                input(func(value string) hypp.Effect {
                    return hypp.Effect{
                        Effecter: getURL,
                        Payload: strings.TrimSpace(value),
                    }
                }, map[string]interface{}{
                    "placeholder": "Search GIFs...",
                    "type": "text",
                }),
                content,
            )
        },
        Node: node,
    })
}

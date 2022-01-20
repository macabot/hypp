// This file is based on https://codepen.io/jorgebucaran/pen/ZeByKv
package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type MyState struct {
	hypp.EmptyState
	isFetching bool
	query      string
	err        error
	url        string
}

func (m MyState) clone() *MyState {
	return &m
}

type giphyBody struct {
	Data []struct {
		Images struct {
			Original struct {
				URL string `json:"url"`
			} `json:"original"`
		} `json:"images"`
	} `json:"data"`
}

var giphyURL = "https://api.giphy.com/v1/gifs/search"
var APIKey = ""

type requestProps struct {
	url   string
	onErr func(err error) hypp.Dispatchable
	onOK  func(items giphyBody) hypp.Dispatchable
}

func request(dispatch hypp.Dispatch, payload hypp.Payload) {
	props := payload.(requestProps)
	go func() {
		res, err := http.Get(props.url)
		if err != nil {
			dispatch(props.onErr(err), nil)
			return
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			dispatch(props.onErr(fmt.Errorf("unexpected status code %d", res.StatusCode)), nil)
			return
		}
		var body giphyBody
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			dispatch(props.onErr(err), nil)
			return
		}
		dispatch(props.onOK(body), nil)
	}()
}

func getJSON(url string, props requestProps) hypp.Effect {
	props.url = url
	return hypp.Effect{
		Effecter: request,
		Payload:  props,
	}
}

func input[S hypp.State](oninput func(value string) hypp.ActionAndPayload[S], props hypp.HProps) *hypp.VNode {
	props.Set("oninput", hypp.Action[*MyState](func(_ *MyState, payload hypp.Payload) hypp.Dispatchable {
		return oninput(payload.(hypp.Event).Target().Value())
	}))
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

func downloadGif(query string) hypp.Effect {
	return getJSON(
		fmt.Sprintf("%s?q=%s&api_key=%s", giphyURL, query, APIKey),
		requestProps{
			onErr: func(err error) hypp.Dispatchable {
				return hypp.ActionAndPayload[*MyState]{
					Action:  gotError,
					Payload: err,
				}
			},
			onOK: func(body giphyBody) hypp.Dispatchable {
				payload := ""
				if len(body.Data) > 0 {
					payload = body.Data[0].Images.Original.URL
				}
				return hypp.ActionAndPayload[*MyState]{
					Action:  gotURL,
					Payload: payload,
				}
			},
		},
	)
}

var errUnexpected = errors.New("Unexpected error, try again later?")

func gotError(state *MyState, payload hypp.Payload) hypp.Dispatchable {
	err := payload.(error)
	state = state.clone()
	state.isFetching = false
	if err != nil {
		state.err = err
	} else {
		state.err = errUnexpected
	}
	state.url = ""
	return state
}

func gotURL(state *MyState, payload hypp.Payload) hypp.Dispatchable {
	url := payload.(string)
	state = state.clone()
	state.isFetching = false
	if state.query != "" {
		state.url = url
	} else {
		state.url = ""
	}
	return state
}

func getURL(state *MyState, payload hypp.Payload) hypp.Dispatchable {
	query := payload.(string)
	state = state.clone()
	state.isFetching = true
	state.query = query
	state.err = nil
	state.url = ""
	return hypp.StateAndEffects[*MyState]{
		State: state,
		Effects: []hypp.Effect{
			downloadGif(query),
		},
	}
}

func Run(driver hypp.Driver, node hypp.Node) {
	hypp.App(hypp.AppProps[*MyState]{
		Driver: driver,
		Init:   &MyState{},
		View: func(state *MyState) *hypp.VNode {
			var content *hypp.VNode
			if state.err != nil {
				content = p(state.err.Error())
			} else {
				content = img(state.url)
			}
			return html.Main(
				nil,
				title("GIF Search 💬💁‍♂️"),
				input(func(value string) hypp.ActionAndPayload[*MyState] {
					return hypp.ActionAndPayload[*MyState]{
						Action:  getURL,
						Payload: strings.TrimSpace(value),
					}
				}, map[string]interface{}{
					"placeholder": "Search GIFs...",
					"type":        "text",
				}),
				content,
			)
		},
		Node: node,
	})
}

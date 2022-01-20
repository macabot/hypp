// This file is based on https://codepen.io/jorgebucaran/pen/evOZLv
package app

import (
	"fmt"
	"time"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type MyState struct {
	hypp.EmptyState
	count  time.Duration
	paused bool
}

func (m MyState) clone() *MyState {
	return &m
}

type timeProps struct {
	delay        time.Duration
	dispatchable hypp.Dispatchable
}

func humanizeTime(t time.Duration) string {
	return fmt.Sprintf("%02.f:%02.f", t.Minutes(), t.Seconds())
}

func interval(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(timeProps)
	ticker := time.NewTicker(props.delay)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				dispatch(props.dispatchable, t)
			}
		}
	}()
	return func() {
		ticker.Stop()
	}
}

func every(delay time.Duration, dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: interval,
		Payload:    timeProps{delay: delay, dispatchable: dispatchable},
	}
}

func titlef(format string, args ...interface{}) *hypp.VNode {
	return html.H1(nil, hypp.Textf(format, args...))
}

func button(onclick hypp.Action[*MyState], text string) *hypp.VNode {
	return html.Button(hypp.HProps{"onclick": onclick}, hypp.Text(text))
}

var resetState = MyState{count: 10 * time.Second}

func tick(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	if state.count == 0 {
		newState := state.clone()
		newState.count = resetState.count
		newState.paused = !state.paused
		return newState
	} else if !state.paused {
		newState := state.clone()
		newState.count -= time.Second
		return newState
	} else {
		return state
	}
}

func reset(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.count = resetState.count
	return newState
}

func toggle(state *MyState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.paused = !state.paused
	return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
	state := resetState.clone()
	state.paused = true
	hypp.App(hypp.AppProps[*MyState]{
		Driver: driver,
		Init:   state,
		View: func(state *MyState) *hypp.VNode {
			var startStop string
			if state.paused {
				startStop = "▶️ Start"
			} else {
				startStop = "Pause ✋"
			}
			return html.Main(
				nil,
				titlef("⏱ %s", humanizeTime(state.count)),
				button(toggle, startStop),
				button(reset, "Reset"),
			)
		},
		Subscriptions: func(state *MyState) []hypp.Subscription {
			return []hypp.Subscription{
				every(time.Second, hypp.Action[*MyState](tick)),
			}
		},
		Node: node,
	})
}

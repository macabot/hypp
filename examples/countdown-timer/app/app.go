// Package app implements a countdown timer.
// Click Start/Pause to start or pause the countdown.
// Click Reset to reset the countdown.
// If the timer is running when clicking Reset, then the timer will continue the countdown after the timer is reset.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/evOZLv

import (
	"fmt"
	"time"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

type State struct {
	hypp.EmptyState
	count  time.Duration
	paused bool
}

func (m State) clone() *State {
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

func button(onclick hypp.Action[*State], text string) *hypp.VNode {
	return html.Button(hypp.HProps{"onclick": onclick}, hypp.Text(text))
}

var resetState = State{count: 10 * time.Second}

func tick(state *State, _ hypp.Payload) hypp.Dispatchable {
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

func reset(state *State, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.count = resetState.count
	return newState
}

func toggle(state *State, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.paused = !state.paused
	return newState
}

func Run(driver hypp.Driver, node hypp.Node) {
	state := resetState.clone()
	state.paused = true
	hypp.App(hypp.AppProps[*State]{
		Driver: driver,
		Init:   state,
		View: func(state *State) *hypp.VNode {
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
		Subscriptions: func(state *State) []hypp.Subscription {
			return []hypp.Subscription{
				every(time.Second, hypp.Action[*State](tick)),
			}
		},
		Node: node,
	})
}

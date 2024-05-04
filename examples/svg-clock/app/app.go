// Package app renders a clock with a single hand that marks the passage of every second.
package app

// This file is based on https://codepen.io/jorgebucaran/pen/PWMBLp

import (
	"math"
	"time"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"github.com/macabot/hypp/tag/svg"
	"github.com/macabot/hypp/window"
)

type State struct {
	time time.Time
}

func (m State) clone() *State {
	return &m
}

func angle(t time.Time) float64 {
	return (2 * math.Pi * float64(t.Unix())) / 60
}

type timeProps struct {
	delay        time.Duration
	dispatchable hypp.Dispatchable
}

func interval(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(timeProps)
	ticker := time.NewTicker(props.delay)
	go func() {
		for {
			t := <-ticker.C
			dispatch(props.dispatchable, t)
		}
	}()
	return func() {
		ticker.Stop()
	}
}

func getTime(dispatch hypp.Dispatch, payload hypp.Payload) {
	props := payload.(timeProps)
	dispatch(props.dispatchable, time.Now())
}

func every(delay time.Duration, dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: interval,
		Payload:    timeProps{delay: delay, dispatchable: dispatchable},
	}
}

func now(dispatchable hypp.Dispatchable) hypp.Effect {
	return hypp.Effect{
		Effecter: getTime,
		Payload:  timeProps{dispatchable: dispatchable},
	}
}

func tick(state *State, payload hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.time = payload.(time.Time)
	return newState
}

func Run(node window.Element) {
	hypp.App(hypp.AppProps[*State]{
		Init: hypp.StateAndEffects[*State]{
			State: &State{},
			Effects: []hypp.Effect{
				now(tick),
			},
		},
		View: func(state *State) *hypp.VNode {
			return html.Svg(
				hypp.HProps{
					"viewBox":      "0 0 100 100",
					"width":        "40%",
					"stroke-width": 2,
				},
				svg.Circle(hypp.HProps{
					"cx":     50,
					"cy":     50,
					"r":      45,
					"stroke": "#0366d6",
					"fill":   "white",
				}),
				svg.Line(hypp.HProps{
					"x1":           50,
					"y1":           50,
					"x2":           50 + 40*math.Cos(angle(state.time)),
					"y2":           50 + 40*math.Sin(angle(state.time)),
					"stroke":       "#0366d6",
					"stroke-width": 3,
				}),
			)
		},
		Subscriptions: func(state *State) []hypp.Subscription {
			return []hypp.Subscription{
				every(time.Second, tick),
			}
		},
		Node: node,
	})
}

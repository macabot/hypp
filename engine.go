// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.js
package hypp

import (
    "strings"
)

var ssrNode = 1
var textNode = 3
var svgNS = "http://www.w3.org/2000/svg"

func h(tag string, props HProps, children []VNode) VNode {
    return VNode{
        tag: tag,
        props: props,
        key: props.key(),
        children: children,
        kind: ssrNode,
    }
}

func memo(view func(data interface{}) VNode, data interface{}) VNode {
    return VNode{
        memoView: view,
        memo: data,
    }
}

func text(value string, node Node) VNode {
    return VNode{
        tag: value,
        kind: textNode,
        node: node,
    }
}

func dispatchInitializerID(dispatch Dispatch) Dispatch {
    return dispatch
}

func dispatchID(dispatchable Dispatchable, payload Payload) {}

func subscriptionsID[S State](state S) []Subscription[S] {
    return nil
}

func renderID() {}

func recycleNode(node Node) VNode {
    if node.nodeType() == textNode {
        return text(node.nodeValue(), node)
    } else {
        childNodes := node.childNodes()
        children := make([]VNode, len(childNodes))
        for i, childNode := range childNodes {
            children[i] = recycleNode(childNode)
        }
        return VNode{
            tag: strings.ToLower(node.nodeName()),
            children: children,
            kind: ssrNode,
            node: node,
        }
    }
}

func shouldRestart(a, b Payload) bool {
    return true // TODO implement
}

func patchSubs[S State](oldSubs, newSubs []Subscription[S], dispatch Dispatch) []Subscription[S] {
    var subs []Subscription[S]
    for i := 0; i < len(oldSubs) || i < len(newSubs); i++ {
        oldSub := oldSubs[i]
        newSub := newSubs[i]
        var sub Subscription[S]
        if !newSub.Disabled {
            if oldSub.Disabled || &newSub.Subscriber != &oldSub.Subscriber || shouldRestart(newSub.Payload, oldSub.Payload) {
                if !oldSub.Disabled {
                    oldSub.unsubscribe()
                }
                sub = Subscription[S]{
                    Subscriber: newSub.Subscriber,
                    Payload: newSub.Payload,
                    unsubscribe: newSub.Subscriber(dispatch, newSub.Payload),
                }
            } else {
                sub = oldSub
            }
        } else {
            if !oldSub.Disabled {
                oldSub.unsubscribe()
            }
            sub = Subscription[S]{
                Disabled: true,
            }
        }
        subs = append(subs, sub)
    }
    return subs
}

func enqueue(render func(), busy bool) {
    // TODO implement
}

func patch(
    parentNode Node,
    node Node,
    vdomOld VNode,
    vdom VNode,
    listener func(VNode, Event),
    busy bool,
) Option[Node] {
    return Option[Node]{
        V: node, // TODO implement
        OK: true,
    }
}

type Event struct {
    kind string // TODO implement
}

type EmptyState struct{}

func (_ EmptyState) IAmDispatchable() {}

func (a *AppProps[S]) init() {
    if a.DispatchInitializer == nil {
        a.DispatchInitializer = dispatchInitializerID
    }
    if a.Init == nil {
        a.Init = EmptyState{}
    }
}

func update[S State](appProps *AppProps[S], newState S) {
    if (appProps.state != newState) {
        appProps.state = newState
        // if appProps.state == EmptyState{} { // FIXME
        //     appProps.dispatch = dispatchID
        //     appProps.subscriptions = subscriptionsID[S]
        //     appProps.render = renderID
        // }
        if appProps.Subscriptions != nil {
            appProps.subs = patchSubs(
                appProps.subs,
                appProps.Subscriptions(appProps.state),
                appProps.dispatch,
            )
        }
        if appProps.View != nil && !appProps.busy {
            appProps.busy = true
            enqueue(appProps.render, appProps.busy)
        }
    }
}

func app[S State](appProps AppProps[S]) Dispatch {
    appProps.init()
    var dispatch Dispatch

    // appProps.vdom = appProps.Node.V // FIXME
    // if appProps.Node.OK {
    //     appProps.vdom = recycleNode(appProps.Node.V)
    // }

    listener := func(this VNode, event Event) {
        dispatch(this.events[event.kind], event)
    }

    appProps.render = func() {
        vdomOld := appProps.vdom
        appProps.vdom = appProps.View(appProps.state)
        appProps.busy = false
        appProps.Node = patch(
            appProps.Node.V.parentNode().V,
            appProps.Node.V,
            vdomOld,
            appProps.vdom,
            listener,
            appProps.busy,
        )
    }

    dispatch = func(dispatchable Dispatchable, props Payload) {
		switch v := dispatchable.(type) {
		case StateAndEffects[S]:
			update[S](&appProps, v.State)
			for _, effect := range v.Effects {
				effect.Effecter(dispatch, effect.Payload)
			}
		case Action[S]:
			dispatch(v(appProps.state, props), nil)
		case ActionAndPayload[S]:
			dispatch(v.Action, v.Payload)
		case S: // State
			update[S](&appProps, v)
		}
	}
	dispatch = appProps.DispatchInitializer(dispatch)
	dispatch(appProps.Init, nil)

	return dispatch
}

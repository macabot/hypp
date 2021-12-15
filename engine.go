// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.js
package hypp

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

func text(value string, node Option[Node]) VNode {
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

func subscriptionsID[S any](state S) []Subscription[S] {
    return nil
}

func renderID() {}

func recycleNode(node VNode) VNode {
    return node // TODO implement
}

func patchSubs[S any](oldSubs, newSubs []Subscription[S], dispatch Dispatch) []Subscription[S] {
    return newSubs // TODO implement
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

func update[S comparable](props *AppProps[S], newState S) {
    if (props.state != newState) {
        props.state = newState
        // if _, ok := props.state.(EmptyState); ok { // FIXME
        //     props.dispatch = dispatchID
        //     props.subscriptions = subscriptionsID[S]
        //     props.render = renderID
        // }
        if props.Subscriptions != nil {
            props.subs = patchSubs(props.subs, props.Subscriptions(props.state), props.dispatch)
        }
        if props.View != nil && !props.busy {
            props.busy = true
            enqueue(props.render, props.busy)
        }
    }
}

func app[S any](props AppProps[S]) Dispatch {
    props.init()
    var dispatch Dispatch

    // props.vdom = props.Node.V // FIXME
    // if props.Node.OK {
    //     props.vdom = recycleNode(props.Node.V)
    // }

    listener := func(this VNode, event Event) {
        // dispatch(this.events[event.kind], event) // FIXME
    }

    props.render = func() {
        vdomOld := props.vdom
        props.vdom = props.View(props.state)
        props.busy = false
        props.Node = patch(
            props.Node.V.parentNode.V,
            props.Node.V,
            vdomOld,
            props.vdom,
            listener,
            props.busy,
        )
    }

    dispatch = func(dispatchable Dispatchable, props Payload) {
		switch v := dispatchable.(type) {
		case StateAndEffects[S]:
			update(&props, v.State)
			for _, effect := range v.Effects {
				effect.Effecter(dispatch, effect.Payload)
			}
		case Action:
			dispatch(v(props.state, props), nil)
		case ActionAndPayload:
			dispatch(v.Action, v.Payload)
		default: // State
			update(&props, v)
		}
	}
	dispatch = dispatchInit(dispatch)
	// dispatch(init, nil) // FIXME

	return dispatch
}

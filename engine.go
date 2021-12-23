// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.js
package hypp

import (
	"fmt"
	"strings"
)

var ssrNode = 1
var textNode = 3
var svgNS = "http://www.w3.org/2000/svg"

func h(tag string, props HProps, children []VNode) VNode {
	return VNode{
		tag:      tag,
		props:    props,
		key:      props.Key(),
		children: children,
		kind:     ssrNode,
	}
}

func memo(view func(data interface{}) VNode, data interface{}) VNode {
	return VNode{
		memoView: view,
		memo:     data,
	}
}

func text(value string, node Node) VNode {
	return VNode{
		tag:  value,
		kind: textNode,
		node: node,
	}
}

func dispatchInitializerID(dispatch Dispatch) Dispatch {
	return dispatch
}

func recycleNode(node Node) VNode {
	if node.NodeType() == textNode {
		return text(node.NodeValue(), node)
	} else {
		childNodes := node.ChildNodes()
		children := make([]VNode, len(childNodes))
		for i, childNode := range childNodes {
			children[i] = recycleNode(childNode)
		}
		return VNode{
			tag:      strings.ToLower(node.NodeName()),
			children: children,
			kind:     ssrNode,
			node:     node,
		}
	}
}

func shouldRestart(a, b Payload) bool {
	return a != b // TODO implement
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
					Subscriber:  newSub.Subscriber,
					Payload:     newSub.Payload,
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

func hPropsKeys(oldProps, newProps HProps) []string {
	seen := map[string]struct{}{}
	keys := make([]string, len(oldProps))
	i := 0
	for key := range oldProps {
		seen[key] = struct{}{}
		keys[i] = key
		i++
	}
	for key := range newProps {
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}
	return keys
}

type vNodeMap map[string]VNode

func (s vNodeMap) Has(key string) bool {
	if s == nil {
		return false
	}
	_, ok := s[key]
	return ok
}

func (s vNodeMap) HasOption(key Option[string]) bool {
	if !key.OK {
		return false
	}
	return s.Has(key.V)
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Set(v T) {
	s[v] = struct{}{}
}

func stylePairKeys(a, b map[string]string) []string {
	seen := map[string]struct{}{}
	keys := make([]string, len(a))
	i := 0
	for key := range a {
		seen[key] = struct{}{}
		keys[i] = key
		i++
	}
	for key := range b {
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}
	return keys
}

type ElementCreationOptions struct {
	Is string
}

type Driver interface {
	CreateTextNode(data string) Node
	CreateElementNS(namespaceURI, qualifiedName string, options Option[ElementCreationOptions]) Node
	CreateElement(tagName string, options Option[ElementCreationOptions]) Node
	Enqueue(render func())
}

func createClass(obj interface{}) string {
	var parts []string
	switch v := obj.(type) {
	case string:
		return v
	case []string:
		parts = v
	case map[string]bool:
		for k, ok := range v {
			if ok {
				parts = append(parts, k)
			}
		}
	default:
		return fmt.Sprint(obj)
	}
	return strings.Join(parts, " ")
}

func patchProperty(node Node, key string, oldValue, newValue interface{}, listener EventListenerGenerator, isSvg bool) {
	if key == "key" {
		// Do nothing
	} else if key == "style" {
		oldStyle, ok1 := oldValue.(map[string]string)
		newStyle, ok2 := newValue.(map[string]string)
		if ok1 && ok2 {
			for _, k := range stylePairKeys(oldStyle, newStyle) {
				v := newStyle[k]
				if k[0] == '-' {
					node.Style().SetProperty(k, v)
				} else {
					node.Style().Set(k, v)
				}
			}
		} else {
			node.Set(key, newValue)
		}
	} else if key[0] == 'o' && key[1] == 'n' {
		key := key[2:]
		node.Events().Set(key, newValue.(Event))
		if newValue == nil {
			node.RemoveEventListener(key, listener(node))
		} else if oldValue == nil {
			node.AddEventListener(key, listener(node))
		}
	} else if !isSvg && key != "list" && key != "from" && node.Has(key) {
		if newValue == nil {
			node.Set(key, "")
		} else {
			node.Set(key, newValue)
		}
	} else {
		if newValue != nil && newValue != false && key == "class" {
			newValue = createClass(newValue)
		}
		if newValue == nil || newValue == false {
			node.RemoveAttribute(key)
		} else {
			node.SetAttribute(key, newValue)
		}
	}
}

func createNode(driver Driver, vdom VNode, listener EventListenerGenerator, isSvg bool) Node {
	props := vdom.props
	var node Node
	if vdom.kind == textNode {
		node = driver.CreateTextNode(vdom.tag)
	} else {
		isSvg = isSvg || vdom.tag == "svg"
		options := Option[ElementCreationOptions]{}
		if props.Has("is") {
			options.V.Is = fmt.Sprint(props.Get("is").V)
			options.OK = true
		}
		if isSvg {
			node = driver.CreateElementNS(svgNS, vdom.tag, options)
		} else {
			node = driver.CreateElement(vdom.tag, options)
		}
	}

	for k := range props {
		patchProperty(node, k, nil, props[k], listener, isSvg)
	}

	for i := 0; i < len(vdom.children); i++ {
		vdom.children[i] = maybeVNode(vdom.children[i], VNode{isNil: true})
		node.AppendChild(
			createNode(
				driver,
				vdom.children[i],
				listener,
				isSvg,
			),
		)
	}
	vdom.node = node
	return node
}

func patch(
	driver Driver,
	parent Node,
	node Node,
	oldVNode VNode,
	newVNode VNode,
	listener EventListenerGenerator,
	isSvg bool,
) Node {
	if &oldVNode == &newVNode {
		// Do nothing
	} else if !oldVNode.isNil && oldVNode.kind == textNode && newVNode.kind == textNode {
		if oldVNode.tag != newVNode.tag {
			node.SetNodeValue(newVNode.tag)
		}
	} else if oldVNode.isNil || oldVNode.tag != newVNode.tag {
		newVNode = maybeVNode(newVNode, VNode{isNil: true})
		node = parent.InsertBefore(
			createNode(driver, newVNode, listener, isSvg),
			node,
		)
		if !oldVNode.isNil {
			parent.RemoveChild(oldVNode.node)
		}
	} else {
		var tmpVKid VNode
		var oldVKid VNode

		var oldKey Option[string]
		var newKey Option[string]

		oldProps := oldVNode.props
		newProps := newVNode.props

		oldVKids := oldVNode.children
		newVKids := newVNode.children

		oldHead := 0
		newHead := 0
		oldTail := len(oldVKids) - 1
		newTail := len(newVKids) - 1

		isSvg := isSvg || newVNode.tag == "svg"

		allKeys := hPropsKeys(oldProps, newProps)
		for _, i := range allKeys {
			var cmpVal Option[interface{}]
			if i == "value" || i == "selected" || i == "checked" {
				cmpVal = node.Get(i)
			} else {
				cmpVal = oldProps.Get(i)
			}
			if cmpVal != newProps.Get(i) {
				patchProperty(node, i, oldProps.Get(i), newProps.Get(i), listener, isSvg)
			}
		}

		for newHead <= newTail && oldHead <= oldTail {
			oldKey := oldVKids[oldHead].key
			if !oldKey.OK || oldKey != newVKids[newHead].key {
				break
			}

			newVKids[newHead] = maybeVNode(
				newVKids[newHead],
				oldVKids[oldHead],
			)
			newHead++
			oldHead++
			patch(
				driver,
				node,
				oldVKids[oldHead].node,
				oldVKids[oldHead],
				newVKids[newHead],
				listener,
				isSvg,
			)
		}

		for newHead <= newTail && oldHead <= oldTail {
			oldKey = oldVKids[oldTail].key
			if !oldKey.OK || oldKey != newVKids[newTail].key {
				break
			}
			newVKids[newTail] = maybeVNode(
				newVKids[newTail],
				oldVKids[oldTail],
			)
			newTail--
			oldTail--
			patch(
				driver,
				node,
				oldVKids[oldTail].node,
				oldVKids[oldTail],
				newVKids[newTail],
				listener,
				isSvg,
			)
		}

		if oldHead > oldTail {
			for newHead <= newTail {
				newVKids[newHead] = maybeVNode(newVKids[newHead], VNode{isNil: true})
				newHead++
				oldVKid = oldVKids[oldHead]
				node.InsertBefore(
					createNode(
						driver,
						newVKids[newHead],
						listener,
						isSvg,
					),
					oldVKid.node,
				)
			}
		} else if newHead > newTail {
			for oldHead <= oldTail {
				node.RemoveChild(oldVKids[oldHead].node)
				oldHead++
			}
		} else {
			keyed := vNodeMap{}
			newKeyed := Set[string]{}
			for i := oldHead; i <= oldTail; i++ {
				oldKey = oldVKids[i].key
				if oldKey.OK {
					keyed[oldKey.V] = oldVKids[i]
				}
			}

			for newHead <= newTail {
				oldVKid = oldVKids[oldHead]
				oldKey = oldVKid.key
				newVKids[newHead] = maybeVNode(newVKids[newHead], oldVKid)
				newKey = newVKids[newHead].key

				if (oldKey.OK && newKeyed.Has(oldKey.V)) || newKey.OK && newKey == oldVKids[oldHead+1].key {
					if !oldKey.OK {
						node.RemoveChild(oldVKid.node)
					}
					oldHead++
					continue
				}

				if !newKey.OK || oldVNode.kind == ssrNode {
					if !oldKey.OK {
						patch(
							driver,
							node,
							oldVKid.node,
							oldVKid,
							newVKids[newHead],
							listener,
							isSvg,
						)
						newHead++
					}
					oldHead++
				} else {
					if oldKey == newKey {
						patch(
							driver,
							node,
							oldVKid.node,
							oldVKid,
							newVKids[newHead],
							listener,
							isSvg,
						)
						newKeyed.Set(newKey.V)
						oldHead++
					} else {
						tmpVKid = keyed[newKey.V]
						if !tmpVKid.isNil {
							patch(
								driver,
								node,
								node.InsertBefore(tmpVKid.node, oldVKid.node),
								tmpVKid,
								newVKids[newHead],
								listener,
								isSvg,
							)
							newKeyed.Set(newKey.V)
						} else {
							patch(
								driver,
								node,
								oldVKid.node,
								VNode{isNil: true},
								newVKids[newHead],
								listener,
								isSvg,
							)
						}
					}
					newHead++
				}
			}

			for oldHead <= oldTail {
				oldVKid = oldVKids[oldHead]
				oldHead++
				if !oldVKid.key.OK {
					node.RemoveChild(oldVKid.node)
				}
			}

			for i := range keyed {
				if !newKeyed.Has(i) {
					node.RemoveChild(keyed[i].node)
				}
			}
		}
	}

	newVNode.node = node
	return node
}

func propsChanged(a, b interface{}) bool {
	return true // TODO implement
}

func maybeVNode(newVNode, oldVNode VNode) VNode {
	if !newVNode.isNil {
		if newVNode.memoView != nil {
			if oldVNode.isNil || oldVNode.memo == nil || propsChanged(oldVNode.memo, newVNode.memo) {
				oldVNode = newVNode.memoView(newVNode.memo)
				oldVNode.memo = newVNode.memo
			}
			return oldVNode
		} else {
			return newVNode
		}
	} else {
		return text("", nil)
	}
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
	if appProps.state != newState {
		appProps.state = newState
		if appProps.Subscriptions != nil {
			appProps.subs = patchSubs(
				appProps.subs,
				appProps.Subscriptions(appProps.state),
				appProps.dispatch,
			)
		}
		if appProps.View != nil && !appProps.busy {
			appProps.busy = true
			appProps.Driver.Enqueue(appProps.render)
		}
	}
}

type EventListenerGenerator func(this Node) EventListener

func app[S State](appProps AppProps[S]) Dispatch {
	appProps.init()
	var dispatch Dispatch

	appProps.vdom = recycleNode(appProps.Node)

    listener := func(this Node) EventListener {
        return func(event Event) {
            dispatch(this.Events().Get(event.Type()), event)
        }
    }

	appProps.render = func() {
		vdomOld := appProps.vdom
		appProps.vdom = appProps.View(appProps.state)
		appProps.busy = false
		appProps.Node = patch(
			appProps.Driver,
			appProps.Node.ParentNode().V,
			appProps.Node,
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
		default:
			panic(fmt.Errorf("hypp: dispatchable has unexpected type '%[1]T'. Expected type 'StateAndEffects[%[2]T]', 'Action[%[2]T]', 'ActionAndPayload[%[2]T]' or '%[2]T'", dispatchable, appProps.state))
		}
	}
	dispatch = appProps.DispatchInitializer(dispatch)
	dispatch(appProps.Init, nil)

	return dispatch
}

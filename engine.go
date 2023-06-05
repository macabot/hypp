package hypp

// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.js

import (
	"fmt"
	"sort"
	"strings"
)

const svgNS = "http://www.w3.org/2000/svg"

func validateHProps(props HProps, tag string) {
	for key, value := range props {
		if key[0] == 'o' && key[1] == 'n' {
			if _, ok := value.(Dispatchable); !ok {
				panic(fmt.Errorf("hypp: expected '%s.%s' to have Dispatchable value. Got %+v of type %T\n", tag, key, value, value))
			}
		} else if key == "style" {
			if _, ok := value.(map[string]string); !ok {
				panic(fmt.Errorf("hypp: expected '%s.%s' to have value of type map[string]string. Got %+v of type %T\n", tag, key, value, value))
			}
		} else {
			switch value.(type) {
			case bool, int, float64, string:
				continue
			}
			if key == "class" {
				switch v := value.(type) {
				case []string, map[string]bool:
					// Do nothing
				default:
					panic(fmt.Errorf("hypp: expected '%s.%s' to have value of type bool, int, float64, string, []string or map[string]bool. Got %+v of type %T\n", tag, key, v, v))
				}
			}
		}
	}
}

func h(tag string, props HProps, children vKids) *VNode {
	validateHProps(props, tag)
	key := props.Key()
	props = props.clone()
	delete(props, "key")
	if classOption := props.Get("class"); classOption.OK {
		props.Set("class", createClass(classOption.V))
	}
	return &VNode{
		tag:      tag,
		props:    props,
		key:      key,
		children: children,
		kind:     SSRNode,
	}
}

func memo(view func(data MemoData) *VNode, data MemoData) *VNode {
	return &VNode{
		memoView: view,
		memoData: data,
	}
}

func text(value string, node Node) *VNode {
	return &VNode{
		tag:  value,
		kind: TextNode,
		node: node,
	}
}

func dispatchWrapperID(dispatch Dispatch) Dispatch {
	return dispatch
}

func recycleNode(node Node) *VNode {
	if node.NodeType() == TextNode {
		return text(node.NodeValue(), node)
	} else {
		childNodes := node.ChildNodes()
		children := make([]*VNode, len(childNodes))
		for i, childNode := range childNodes {
			children[i] = recycleNode(childNode)
		}
		return &VNode{
			tag:      strings.ToLower(node.NodeName()),
			children: children,
			kind:     SSRNode,
			node:     node,
		}
	}
}

func shouldRestart(a, b Payload) bool {
	return a != b // TODO implement
}

type subscriptions []Subscription

func (s subscriptions) Index(i int) Subscription {
	if i < len(s) {
		return s[i]
	}
	return Subscription{Disabled: true}
}

func patchSubs(oldSubs, newSubs subscriptions, dispatch Dispatch) []Subscription {
	var subs []Subscription
	for i := 0; i < len(oldSubs) || i < len(newSubs); i++ {
		oldSub := oldSubs.Index(i)
		newSub := newSubs.Index(i)
		var sub Subscription
		if !newSub.Disabled {
			if oldSub.Disabled || &newSub.Subscriber != &oldSub.Subscriber || shouldRestart(newSub.Payload, oldSub.Payload) {
				if !oldSub.Disabled {
					oldSub.unsubscribe()
				}
				sub = Subscription{
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
			sub = Subscription{
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

func stylePairKeys(x, y interface{}) []string {
	seen := map[string]struct{}{}
	var keys []string

	if a, ok := x.(map[string]string); ok {
		for key := range a {
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}

	if b, ok := y.(map[string]string); ok {
		for key := range b {
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				keys = append(keys, key)
			}
		}
	}

	return keys
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
		sort.Strings(parts)
	default:
		return fmt.Sprint(obj)
	}
	return strings.Join(parts, " ")
}

func isFalsy(v interface{}) bool {
	return v == false || v == 0 || v == 0.0 || v == "" || v == nil
}

func patchProperty(node Node, key string, oldValue, newValue interface{}, listener eventListenerGenerator, isSvg bool) {
	if key == "style" {
		newStyle, isNewStyle := newValue.(map[string]string)
		for _, k := range stylePairKeys(oldValue, newValue) {
			oldValue := ""
			if isNewStyle && newStyle != nil {
				oldValue = newStyle[k]
			}
			if k[0] == '-' {
				node.SetStyleProperty(k, oldValue)
			} else {
				node.SetStyle(k, oldValue)
			}
		}
	} else if key[0] == 'o' && key[1] == 'n' {
		key := key[2:]
		if isFalsy(newValue) {
			node.Events().Del(key)
			if id := node.EventListenerID(key); id != nil {
				node.RemoveEventListener(key, id)
			}
		} else {
			if _, ok := newValue.(Dispatchable); !ok {
				panic(fmt.Errorf("hypp: expected Dispatchable for key starting with 'on'. Key: %s, value: %+v of type %T, %s\n", key, newValue, newValue, newValue))
			}
			node.Events().Set(key, newValue.(Dispatchable))
			if isFalsy(oldValue) {
				id := node.AddEventListener(key, listener(node))
				node.SetEventListenerID(key, id)
			}
		}
	} else if !isSvg && key != "list" && key != "form" && node.In(key) {
		if isFalsy(newValue) {
			node.Set(key, "")
		} else {
			node.Set(key, newValue)
		}
	} else {
		if isFalsy(newValue) {
			node.RemoveAttribute(key)
		} else {
			node.SetAttribute(key, newValue)
		}
	}
}

func createNode(driver Driver, vdom *VNode, listener eventListenerGenerator, isSvg bool) Node {
	props := vdom.props
	var node Node
	if vdom.kind == TextNode {
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
		vdom.children[i] = maybeVNode(vdom.children[i], nil)
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

func equalProps(a, b Option[interface{}]) bool {
	if a.OK != b.OK {
		return false
	}
	switch v := a.V.(type) {
	case int:
		bInt, ok := b.V.(int)
		return ok && v == bInt
	case string:
		bString, ok := b.V.(string)
		return ok && v == bString
	case bool:
		bBool, ok := b.V.(bool)
		return ok && v == bBool
	default:
		return false
	}
}

func patch(
	driver Driver,
	parent Node,
	node Node,
	oldVNode *VNode,
	newVNode *VNode,
	listener eventListenerGenerator,
	isSvg bool,
) Node {
	if oldVNode == newVNode {
		// Do nothing
	} else if oldVNode != nil && oldVNode.kind == TextNode && newVNode.kind == TextNode {
		if oldVNode.tag != newVNode.tag {
			node.SetNodeValue(newVNode.tag)
		}
	} else if oldVNode == nil || oldVNode.tag != newVNode.tag {
		newVNode = maybeVNode(newVNode, nil)
		node = parent.InsertBefore(
			createNode(driver, newVNode, listener, isSvg),
			node,
		)
		if oldVNode != nil {
			parent.RemoveChild(oldVNode.node)
		}
	} else {
		var tmpVKid *VNode
		var oldVKid *VNode

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
			if !equalProps(cmpVal, newProps.Get(i)) {
				patchProperty(node, i, oldProps.Get(i).V, newProps.Get(i).V, listener, isSvg)
			}
		}

		for newHead <= newTail && oldHead <= oldTail {
			oldKey = oldVKids.getKey(oldHead)
			if !oldKey.OK || oldKey != newVKids.getKey(newHead) {
				break
			}

			newVKids[newHead] = maybeVNode(
				newVKids[newHead],
				oldVKids[oldHead],
			)
			patch(
				driver,
				node,
				oldVKids[oldHead].node,
				oldVKids[oldHead],
				newVKids[newHead],
				listener,
				isSvg,
			)
			newHead++
			oldHead++
		}

		for newHead <= newTail && oldHead <= oldTail {
			oldKey = oldVKids.getKey(oldTail)
			if !oldKey.OK || oldKey != newVKids.getKey(newTail) {
				break
			}
			newVKids[newTail] = maybeVNode(
				newVKids[newTail],
				oldVKids[oldTail],
			)
			patch(
				driver,
				node,
				oldVKids[oldTail].node,
				oldVKids[oldTail],
				newVKids[newTail],
				listener,
				isSvg,
			)
			newTail--
			oldTail--
		}

		if oldHead > oldTail {
			for newHead <= newTail {
				newVKids[newHead] = maybeVNode(newVKids[newHead], nil)
				oldVKid = nil
				var oldVKidNode Node
				if oldHead < len(oldVKids) {
					oldVKid = oldVKids[oldHead]
					oldVKidNode = oldVKid.node
				}
				node.InsertBefore(
					createNode(
						driver,
						newVKids[newHead],
						listener,
						isSvg,
					),
					oldVKidNode,
				)
				newHead++
			}
		} else if newHead > newTail {
			for oldHead <= oldTail {
				node.RemoveChild(oldVKids[oldHead].node)
				oldHead++
			}
		} else {
			keyed := map[string]*VNode{}
			newKeyed := Set[string]{}
			for i := oldHead; i <= oldTail; i++ {
				oldKey = oldVKids[i].key
				if oldKey.OK {
					keyed[oldKey.V] = oldVKids[i]
				}
			}

			for newHead <= newTail {
				oldVKid = oldVKids.get(oldHead)
				oldKey = oldVKids.getKey(oldHead)
				newVKids[newHead] = maybeVNode(newVKids[newHead], oldVKid)
				newKey = newVKids.getKey(newHead)

				if (newKeyed.Has(oldKey.V)) || (newKey.OK && newKey == oldVKids.getKey(oldHead+1)) {
					if !oldKey.OK {
						node.RemoveChild(oldVKid.node)
					}
					oldHead++
					continue
				}

				if !newKey.OK || oldVNode.kind == SSRNode {
					if !oldKey.OK {
						var oldVKidNode Node
						if oldVKid != nil {
							oldVKidNode = oldVKid.node
						}
						patch(
							driver,
							node,
							oldVKidNode,
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
						newKeyed.Add(newKey.V)
						oldHead++
					} else {
						tmpVKid = keyed[newKey.V]
						if tmpVKid != nil {
							patch(
								driver,
								node,
								node.InsertBefore(tmpVKid.node, oldVKid.node),
								tmpVKid,
								newVKids[newHead],
								listener,
								isSvg,
							)
							newKeyed.Add(newKey.V)
						} else {
							patch(
								driver,
								node,
								oldVKid.node,
								nil,
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
				oldVKid = oldVKids.get(oldHead)
				if !oldVKids.getKey(oldHead).OK {
					node.RemoveChild(oldVKid.node)
				}
				oldHead++
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

func propsChanged(a, b MemoData) bool {
	return a.Hash() != b.Hash()
}

func maybeVNode(newVNode, oldVNode *VNode) *VNode {
	if newVNode != nil {
		if newVNode.memoView != nil {
			if oldVNode == nil || oldVNode.memoData == nil || propsChanged(oldVNode.memoData, newVNode.memoData) {
				oldVNode = newVNode.memoView(newVNode.memoData)
				oldVNode.memoData = newVNode.memoData
			}
			return oldVNode
		} else {
			return newVNode
		}
	} else {
		return text("", nil)
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
			appProps.Driver.Window().RequestAnimationFrame(appProps.render)
		}
	}
}

type eventListenerGenerator func(this Node) EventListener

func app[S State](appProps AppProps[S]) Dispatch {
	appProps.init()

	appProps.vdom = recycleNode(appProps.Node)

	listener := func(this Node) EventListener {
		return func(event Event) {
			appProps.dispatch(this.Events().Get(event.Type()), event)
		}
	}

	appProps.render = func() {
		vdomOld := appProps.vdom
		appProps.vdom = appProps.View(appProps.state)
		appProps.busy = false
		appProps.Node = patch(
			appProps.Driver,
			appProps.Node.ParentNode(),
			appProps.Node,
			vdomOld,
			appProps.vdom,
			listener,
			appProps.busy,
		)
	}

	appProps.dispatch = func(dispatchable Dispatchable, props Payload) {
		switch v := dispatchable.(type) {
		case StateAndEffects[S]:
			update(&appProps, v.State)
			for _, effect := range v.Effects {
				effect.Effecter(appProps.dispatch, effect.Payload)
			}
		case Action[S]:
			appProps.dispatch(v(appProps.state, props), nil)
		case ActionAndPayload[S]:
			appProps.dispatch(v.Action, v.Payload)
		case S: // State
			update(&appProps, v)
		default:
			panic(fmt.Errorf("hypp: dispatchable has unexpected type '%[1]T'. Expected type 'StateAndEffects[%[2]T]', 'Action[%[2]T]', 'ActionAndPayload[%[2]T]' or '%[2]T'", dispatchable, appProps.state))
		}
	}
	appProps.dispatch = appProps.DispatchWrapper(appProps.dispatch)
	appProps.dispatch(appProps.Init, nil)

	return appProps.dispatch
}

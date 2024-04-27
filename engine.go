package hypp

// This file is based on https://github.com/jorgebucaran/hyperapp/blob/main/index.js

import (
	"fmt"
	"sort"
	"strings"

	"github.com/macabot/hypp/js"
	"github.com/macabot/hypp/window"
)

const svgNS = "http://www.w3.org/2000/svg"

// ValidateHProps validates the given [HProps].
// It returns an error if any of the keys have an invalid value.
// See the description of [HProps] for a table of the keys and their allowed value types.
func ValidateHProps(props HProps) error {
	for key, value := range props {
		if len(key) >= 2 && key[0] == 'o' && key[1] == 'n' {
			if _, ok := value.(Dispatchable); !ok {
				return fmt.Errorf("hypp: expected HProps key '%s' to have Dispatchable value. Got %#v", key, value)
			}
		} else if key == "class" {
			switch value.(type) {
			case bool, int, float64, string, []string, map[string]bool:
				// Do nothing
			default:
				return fmt.Errorf("hypp: expected HProps key '%s' to have value of type bool, int, float64, string, []string or map[string]bool. Got %#v", key, value)
			}
		} else if key == "style" {
			if _, ok := value.(map[string]string); !ok {
				return fmt.Errorf("hypp: expected HProps key '%s' to have value of type map[string]string. Got %#v", key, value)
			}
		} else {
			switch value.(type) {
			case bool, int, float64, string:
				// Do nothing
			default:
				return fmt.Errorf("hypp: expected HProps key '%s' to have value of type bool, int, float64 or string. Got %#v", key, value)
			}
		}
	}
	return nil
}

// MustValidateHProps calls [ValidateHProps] and panics if an error is returned.
func MustValidateHProps(props HProps) {
	if err := ValidateHProps(props); err != nil {
		panic(err)
	}
}

func h(tag string, props HProps, children vKids) *VNode {
	MustValidateHProps(props)
	props = props.clone()
	if classOption := props.get("class"); classOption.OK {
		props.Set("class", createClass(classOption.V))
	}
	return &VNode{
		tag:      tag,
		props:    props,
		children: children,
		kind:     ElementNode,
	}
}

func memo(view func(data MemoData) *VNode, data MemoData) *VNode {
	return &VNode{
		memoView: view,
		memoData: data,
	}
}

func text(value string, node window.Element) *VNode {
	return &VNode{
		tag:  value,
		kind: TextNode,
		node: node,
	}
}

func dispatchWrapperID(dispatch Dispatch) Dispatch {
	return dispatch
}

func recycleNode(node window.Element) *VNode {
	if node.NodeType() == int(TextNode) {
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
			kind:     ElementNode,
			node:     node,
		}
	}
}

// shouldRestart returns true if a is not equal to b.
// If a and b are not comparable it always returns false.
func shouldRestart(a, b Payload) bool {
	defer func() {
		recover()
	}()
	return a != b
}

type subscriptions []Subscription

func (s subscriptions) Index(i int) Subscription {
	if i < len(s) {
		return s[i]
	}
	return Subscription{Disabled: true}
}

// patchSubs loops over the oldSubs and newSubs and compares the subscriptions at the same index.
// It compares the Disabled field and Payload field. It does not compare the Subscriber field, because Go can't compare functions.
func patchSubs(oldSubs, newSubs subscriptions, dispatch Dispatch) []Subscription {
	var subs []Subscription
	for i := 0; i < len(oldSubs) || i < len(newSubs); i++ {
		oldSub := oldSubs.Index(i)
		newSub := newSubs.Index(i)
		var sub Subscription
		if !newSub.Disabled {
			if oldSub.Disabled || shouldRestart(newSub.Payload, oldSub.Payload) {
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

func mergeAndSortHPropsKeys(propsSlice ...HProps) []string {
	seen := map[string]struct{}{}
	var keys []string
	for _, props := range propsSlice {
		for key := range props {
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				keys = append(keys, key)
			}
		}
	}
	sort.Strings(keys)
	return keys
}

func stylePairKeys(x, y any) []string {
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

func createClass(obj any) string {
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

func isFalsy(v any) bool {
	return v == false || v == 0 || v == 0.0 || v == "" || v == nil
}

// propertyInObject returns true if the specified property is in the specified object or its prototype chain.
// See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/in
func propertyInObject(property string, object js.Value) bool {
	getPrototypeOf := js.Global().Get("Object").Get("getPrototypeOf")
	for v := object; !v.IsNull(); v = getPrototypeOf.Invoke(v) {
		if v.Call("hasOwnProperty", property).Bool() {
			return true
		}
	}
	return false
}

func patchProperty(node window.Element, key string, oldValue, newValue any, listener eventListenerGenerator, isSvg bool) {
	if key == "key" {
		// no-op
	} else if key == "style" {
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
	} else if len(key) >= 2 && key[0] == 'o' && key[1] == 'n' {
		key := key[2:]
		if isFalsy(newValue) {
			getEvents(node).del(key)
			if id := node.EventListenerID(key); !id.IsUndefined() {
				node.RemoveEventListener(key, id)
			}
		} else {
			getEvents(node).set(key, newValue.(Dispatchable))
			if isFalsy(oldValue) {
				id := node.AddEventListener(key, listener(node))
				node.SetEventListenerID(key, id)
			}
		}
	} else if !isSvg && key != "list" && key != "form" && propertyInObject(key, node.Value) {
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

func createNode(vdom *VNode, listener eventListenerGenerator, isSvg bool) window.Element {
	props := vdom.props
	var node window.Element
	if vdom.kind == TextNode {
		node = window.Document().CreateTextNode(vdom.tag)
	} else {
		isSvg = isSvg || vdom.tag == "svg"
		var options *window.ElementCreationOptions
		if props.Has("is") {
			options = &window.ElementCreationOptions{
				Is: fmt.Sprint(props.get("is").V),
			}
		}
		if isSvg {
			node = window.Document().CreateElementNS(svgNS, vdom.tag, options)
		} else {
			node = window.Document().CreateElement(vdom.tag, options)
		}
	}

	keys := mergeAndSortHPropsKeys(props)
	for _, k := range keys {
		patchProperty(node, k, nil, props[k], listener, isSvg)
	}

	for i := 0; i < len(vdom.children); i++ {
		vdom.children[i] = maybeVNode(vdom.children[i], nil)
		node.AppendChild(
			createNode(
				vdom.children[i],
				listener,
				isSvg,
			),
		)
	}
	vdom.node = node
	return node
}

func equalProps(a, b option[any]) bool {
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

func elementGet(e window.Element, name string) option[any] {
	if !propertyInObject(name, e.Value) {
		return option[any]{}
	}
	v := e.Value.Get(name)
	kind := v.Type()
	switch kind {
	case js.TypeUndefined, js.TypeNull:
		return option[any]{OK: true}
	case js.TypeBoolean:
		return option[any]{V: v.Bool(), OK: true}
	case js.TypeNumber:
		if js.Global().Get("Number").Call("isInteger", v).Bool() {
			return option[any]{V: v.Int(), OK: true}
		} else {
			return option[any]{V: v.Float(), OK: true}
		}
	case js.TypeString:
		return option[any]{V: v.String(), OK: true}
	default:
		panic(fmt.Errorf("hypp: cannot get property of window.Element with type '%s'", kind))
	}
}

func patch(
	parent window.Element,
	node window.Element,
	oldVNode *VNode,
	newVNode *VNode,
	listener eventListenerGenerator,
	isSvg bool,
) window.Element {
	if oldVNode == newVNode {
		// Do nothing
	} else if oldVNode != nil && oldVNode.kind == TextNode && newVNode.kind == TextNode {
		if oldVNode.tag != newVNode.tag {
			node.SetNodeValue(newVNode.tag)
		}
	} else if oldVNode == nil || oldVNode.tag != newVNode.tag {
		newVNode = maybeVNode(newVNode, nil)
		node = parent.InsertBefore(
			createNode(newVNode, listener, isSvg),
			node,
		)
		if oldVNode != nil {
			removeChild(parent, oldVNode.node)
		}
	} else {
		var tmpVKid *VNode
		var oldVKid *VNode

		var oldKey option[string]
		var newKey option[string]

		oldProps := oldVNode.props
		newProps := newVNode.props

		oldVKids := oldVNode.children
		newVKids := newVNode.children

		oldHead := 0
		newHead := 0
		oldTail := len(oldVKids) - 1
		newTail := len(newVKids) - 1

		isSvg := isSvg || newVNode.tag == "svg"

		allKeys := mergeAndSortHPropsKeys(oldProps, newProps)
		for _, i := range allKeys {
			var cmpVal option[any]
			if i == "value" || i == "selected" || i == "checked" {
				cmpVal = elementGet(node, i)
			} else {
				cmpVal = oldProps.get(i)
			}
			if !equalProps(cmpVal, newProps.get(i)) {
				patchProperty(node, i, oldProps.get(i).V, newProps.get(i).V, listener, isSvg)
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
				var oldVKidNode window.Element
				if oldHead < len(oldVKids) {
					oldVKid = oldVKids[oldHead]
					oldVKidNode = oldVKid.node
				}
				node.InsertBefore(
					createNode(
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
				removeChild(node, oldVKids[oldHead].node)
				oldHead++
			}
		} else {
			keyed := map[string]*VNode{}
			newKeyed := set[string]{}
			for i := oldHead; i <= oldTail; i++ {
				oldKey = oldVKids[i].key()
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
						removeChild(node, oldVKid.node)
					}
					oldHead++
					continue
				}

				if !newKey.OK || oldVNode.kind == ElementNode {
					if !oldKey.OK {
						var oldVKidNode window.Element
						if oldVKid != nil {
							oldVKidNode = oldVKid.node
						}
						patch(
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
					removeChild(node, oldVKid.node)
				}
				oldHead++
			}

			for i := range keyed {
				if !newKeyed.Has(i) {
					removeChild(node, keyed[i].node)
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
		return text("", window.Element{})
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
			window.RequestAnimationFrame(appProps.render)
		}
	}
}

type eventListenerGenerator func(this window.Element) window.EventListener

func app[S State](appProps AppProps[S]) Dispatch {
	appProps.init()

	appProps.vdom = recycleNode(appProps.Node)

	listener := func(this window.Element) window.EventListener {
		return func(event window.Event) {
			appProps.dispatch(getEvents(this).get(event.Type()), event)
		}
	}

	appProps.render = func() {
		vdomOld := appProps.vdom
		appProps.vdom = appProps.View(appProps.state)
		appProps.busy = false
		appProps.Node = patch(
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
		case func(S, Payload) Dispatchable:
			appProps.dispatch(v(appProps.state, props), nil)
		case ActionAndPayload[S]:
			appProps.dispatch(v.Action, v.Payload)
		case S: // State
			update(&appProps, v)
		default:
			panic(fmt.Errorf("hypp: dispatchable has unexpected type '%[1]T'. Expected type 'StateAndEffects[%[2]T]', 'Action[%[2]T]', 'func(%[2]T, Payload) Dispatchable', 'ActionAndPayload[%[2]T]' or '%[2]T'", dispatchable, appProps.state))
		}
	}
	appProps.dispatch = appProps.DispatchWrapper(appProps.dispatch)
	appProps.dispatch(appProps.Init, nil)

	return appProps.dispatch
}

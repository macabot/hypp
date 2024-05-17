package hypp

import (
	"github.com/macabot/hypp/window"
)

var (
	globalNodeEvents         = newRepo[Dispatchable]()
	globalNodeEventListeners = newRepo[window.EventListenerID]()
)

// removeChild deletes all events and eventListeners before removing the child from its parent.
func removeChild(parent, child window.Element) {
	getEvents(child).deleteAll()

	eventListeners := getEventListeners(child)
	for key, id := range eventListeners.toMap() {
		child.RemoveEventListener(key, id)
	}
	eventListeners.deleteAll()

	parent.RemoveChild(child)
}

// getEvents returns the node's "events" property.
func getEvents(node window.Element) *propertyRepo[Dispatchable] {
	if node.Value.Get("events").IsUndefined() {
		node.Value.Set("events", map[string]any{})
	}
	return newPropertyRepo(node.Value.Get("events"), globalNodeEvents)
}

// getEventListeners returns the node's "eventListeners" property.
func getEventListeners(node window.Element) *propertyRepo[window.EventListenerID] {
	if node.Value.Get("eventListeners").IsUndefined() {
		node.Value.Set("eventListeners", map[string]any{})
	}
	return newPropertyRepo(node.Value.Get("eventListeners"), globalNodeEventListeners)
}

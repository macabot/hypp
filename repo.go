package hypp

import (
	"sync"

	"github.com/macabot/hypp/js"
)

// repo stores values of type T.
// Values are indexed by an auto-incremented index.
type repo[T any] struct {
	mu sync.Mutex
	// i is set to the index of the next value.
	i int
	v map[int]T
}

// newRepo creates a new repo.
func newRepo[T any]() *repo[T] {
	return &repo[T]{v: map[int]T{}}
}

// add adds the given value to the repo and returns its index.
func (r *repo[T]) add(value T) int {
	r.mu.Lock()
	i := r.i
	r.v[i] = value
	r.i++
	r.mu.Unlock()
	return i
}

// set sets the given value at the given index.
func (r *repo[T]) set(i int, value T) {
	r.mu.Lock()
	r.v[i] = value
	r.mu.Unlock()
}

// del deletes the value at the given index.
// Nothing happens if no value was found.
func (r *repo[T]) del(i int) {
	r.mu.Lock()
	delete(r.v, i)
	r.mu.Unlock()
}

// get gets the value at the given index.
// If no value is found, it returns the empty value of T.
func (r *repo[T]) get(i int) T {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.v[i]
}

// propertyRepo is used to link a value of type T to a js.Value.
// A [repo] is used to map T to an index. This index is stored on the js.Value.
type propertyRepo[T any] struct {
	// value must be of type js.TypeObject.
	value js.Value
	repo  *repo[T]
}

// newPropertyRepo creates a new propertyRepo.
//
// The repo must not be nil.
func newPropertyRepo[T any](v js.Value, r *repo[T]) *propertyRepo[T] {
	return &propertyRepo[T]{value: v, repo: r}
}

// set sets the property with the given key to the index corresponding to the given value.
func (r *propertyRepo[T]) set(key string, value T) {
	v := r.value
	if p := v.Get(key); p.Type() == js.TypeNumber {
		r.repo.set(p.Int(), value)
	} else {
		i := r.repo.add(value)
		v.Set(key, i)
	}
}

// get gets the value linked to the property with the given key.
// If no property is found, it returns the empty value of T.
func (r *propertyRepo[T]) get(key string) T {
	v := r.value
	if p := v.Get(key); p.Type() == js.TypeNumber {
		return r.repo.get(p.Int())
	}
	var empty T
	return empty
}

// del deletes the property with the given key.
// Nothing happens if the property is not found.
func (r *propertyRepo[T]) del(key string) {
	v := r.value
	if p := v.Get(key); p.Type() == js.TypeNumber {
		r.repo.del(p.Int())
		v.Delete(key)
	}
}

// keys returns the keys.
func (r *propertyRepo[T]) keys() []string {
	keysValue := js.Global().Get("Object").Call("keys", r.value)
	keys := make([]string, keysValue.Length())
	for i := 0; i < len(keys); i++ {
		keys[i] = keysValue.Index(i).String()
	}
	return keys
}

// deleteAll deletes all properties.
func (r *propertyRepo[T]) deleteAll() {
	for _, key := range r.keys() {
		r.del(key)
	}
}

// toMap converts the keys and values into a map.
func (r *propertyRepo[T]) toMap() map[string]T {
	keys := r.keys()
	items := make(map[string]T, len(keys))
	for _, key := range keys {
		items[key] = r.get(key)
	}
	return items
}

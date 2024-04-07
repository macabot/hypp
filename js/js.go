package js

import "errors"

type Driver interface {
	CopyBytesToGo(dst []byte, src Value) int
	CopyBytesToJS(dst Value, src []byte) int
	FuncOf(fn func(this Value, args []Value) any) Func
	Global() Value
	Null() Value
	Undefined() Value
	ValueOf(x any) Value
	DefaultValueDriver() ValueDriver
	DefaultFuncDriver() FuncDriver
}

var driver Driver

func d() Driver {
	if driver == nil {
		panic(errors.New("hypp/js: driver is not set"))
	}
	return driver
}

// Register registers the driver that will be used by this package.
// A common use case is to register the jsd driver
//
//	import _ "github.com/macabot/hypp/jsd"
//
// Make sure to import it from your main package.
func Register(d Driver) {
	driver = d
}

// GetDriver returns the registered Driver or nil if no Driver has been registered.
func GetDriver() Driver {
	return driver
}

// CopyBytesToGo copies bytes from src to dst.
// It panics if src is not a Uint8Array or Uint8ClampedArray.
// It returns the number of bytes copied, which will be the minimum of the lengths of src and dst.
func CopyBytesToGo(dst []byte, src Value) int {
	return d().CopyBytesToGo(dst, src)
}

// CopyBytesToJS copies bytes from src to dst.
// It panics if dst is not a Uint8Array or Uint8ClampedArray.
// It returns the number of bytes copied, which will be the minimum of the lengths of src and dst.
func CopyBytesToJS(dst Value, src []byte) int {
	return d().CopyBytesToJS(dst, src)
}

// FuncOf returns a function to be used by JavaScript.
//
// The Go function fn is called with the value of JavaScript's "this" keyword and the
// arguments of the invocation. The return value of the invocation is
// the result of the Go function mapped back to JavaScript according to ValueOf.
//
// Invoking the wrapped Go function from JavaScript will
// pause the event loop and spawn a new goroutine.
// Other wrapped functions which are triggered during a call from Go to JavaScript
// get executed on the same goroutine.
//
// As a consequence, if one wrapped function blocks, JavaScript's event loop
// is blocked until that function returns. Hence, calling any async JavaScript
// API, which requires the event loop, like fetch (http.Client), will cause an
// immediate deadlock. Therefore a blocking function should explicitly start a
// new goroutine.
//
// Func.Release must be called to free up resources when the function will not be invoked any more.
func FuncOf(fn func(this Value, args []Value) any) Func {
	return d().FuncOf(fn)
}

// Global returns the JavaScript global object, usually "window" or "global".
func Global() Value {
	return d().Global()
}

// Null returns the JavaScript value "null".
func Null() Value {
	return d().Null()
}

// Undefined returns the JavaScript value "undefined".
func Undefined() Value {
	return d().Undefined()
}

// ValueOf returns x as a JavaScript value:
//
//	| Go                     | JavaScript             |
//	| ---------------------- | ---------------------- |
//	| js.Value               | [its value]            |
//	| js.Func                | function               |
//	| nil                    | null                   |
//	| bool                   | boolean                |
//	| integers and floats    | number                 |
//	| string                 | string                 |
//	| []any                  | new array              |
//	| map[string]any         | new object             |
//
// Panics if x is not one of the expected types.
func ValueOf(x any) Value {
	return d().ValueOf(x)
}

// Error wraps a JavaScript error.
type Error struct {
	Value
}

// Error implements the error interface.
func (e Error) Error() string {
	return "JavaScript error: " + e.Get("message").String()
}

type FuncDriver interface {
	ValueDriver() ValueDriver
	Release()
}

// Func is a wrapped Go function to be called by JavaScript.
type Func struct {
	Value
	driver FuncDriver
}

func MakeFunc(driver FuncDriver) Func {
	return Func{Value: MakeValue(driver.ValueDriver()), driver: driver}
}

func (f Func) Driver() FuncDriver {
	if f.driver == nil {
		return d().DefaultFuncDriver()
	}
	return f.driver
}

// Release frees up resources allocated for the function.
// The function must not be invoked after calling Release.
// It is allowed to call Release while the function is still running.
func (f Func) Release() {
	f.Driver().Release()
}

// Type is based on syscall/js.Type.
// See https://pkg.go.dev/syscall/js#Type
type Type int

// These are the valid Type values.
const (
	TypeUndefined Type = iota
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeSymbol
	TypeObject
	TypeFunction
)

// String returns the string value of Type t.
// It panics if t is not one of the valid Type values.
func (t Type) String() string {
	switch t {
	case TypeUndefined:
		return "undefined"
	case TypeNull:
		return "null"
	case TypeBoolean:
		return "boolean"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeSymbol:
		return "symbol"
	case TypeObject:
		return "object"
	case TypeFunction:
		return "function"
	default:
		panic("bad type")
	}
}

type ValueDriver interface {
	Bool() bool
	Call(m string, args ...any) Value
	Delete(p string)
	Equal(w Value) bool
	Float() float64
	Get(p string) Value
	Index(i int) Value
	InstanceOf(t Value) bool
	Int() int
	Invoke(args ...any) Value
	IsNaN() bool
	IsNull() bool
	IsUndefined() bool
	Length() int
	New(args ...any) Value
	Set(p string, x any)
	SetIndex(i int, x any)
	String() string
	Truthy() bool
	Type() Type
}

// Value represents a JavaScript value.
// It is based on type syscall/js.Value.
// It allows you to use the JavaScript environment without the js/wasm build constraint.
// See https://pkg.go.dev/syscall/js#Value for the method definitions.
type Value struct {
	driver ValueDriver
}

func MakeValue(driver ValueDriver) Value {
	return Value{driver}
}

func (v Value) Driver() ValueDriver {
	if v.driver == nil {
		return d().DefaultValueDriver()
	}
	return v.driver
}

// Bool returns the value v as a bool.
// It panics if v is not a JavaScript boolean.
func (v Value) Bool() bool {
	return v.Driver().Bool()
}

// Call does a JavaScript call to the method m of value v with the given arguments.
// It panics if v has no method m.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) Call(m string, args ...any) Value {
	return v.Driver().Call(m, args...)
}

// Delete deletes the JavaScript property p of value v.
// It panics if v is not a JavaScript object.
func (v Value) Delete(p string) {
	v.Driver().Delete(p)
}

// Equal reports whether v and w are equal according to JavaScript's === operator.
func (v Value) Equal(w Value) bool {
	return v.Driver().Equal(w)
}

// Float returns the value v as a float64.
// It panics if v is not a JavaScript number.
func (v Value) Float() float64 {
	return v.Driver().Float()
}

// Get returns the JavaScript property p of value v.
// It panics if v is not a JavaScript object.
func (v Value) Get(p string) Value {
	return v.Driver().Get(p)
}

// Index returns JavaScript index i of value v.
// It panics if v is not a JavaScript object.
func (v Value) Index(i int) Value {
	return v.Driver().Index(i)
}

// InstanceOf reports whether v is an instance of type t according to JavaScript's instanceof operator.
func (v Value) InstanceOf(t Value) bool {
	return v.Driver().InstanceOf(t)
}

// Int returns the value v truncated to an int.
// It panics if v is not a JavaScript number.
func (v Value) Int() int {
	return v.Driver().Int()
}

// Invoke does a JavaScript call of the value v with the given arguments.
// It panics if v is not a JavaScript function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) Invoke(args ...any) Value {
	return v.Driver().Invoke(args...)
}

// IsNaN reports whether v is the JavaScript value "NaN".
func (v Value) IsNaN() bool {
	return v.Driver().IsNaN()
}

// IsNull reports whether v is the JavaScript value "null".
func (v Value) IsNull() bool {
	return v.Driver().IsNull()
}

// IsUndefined reports whether v is the JavaScript value "undefined".
func (v Value) IsUndefined() bool {
	return v.Driver().IsUndefined()
}

// Length returns the JavaScript property "length" of v.
// It panics if v is not a JavaScript object.
func (v Value) Length() int {
	return v.Driver().Length()
}

// New uses JavaScript's "new" operator with value v as constructor and the given arguments.
// It panics if v is not a JavaScript function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) New(args ...any) Value {
	return v.Driver().New(args...)
}

// Set sets the JavaScript property p of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (v Value) Set(p string, x any) {
	v.Driver().Set(p, x)
}

// SetIndex sets the JavaScript index i of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (v Value) SetIndex(i int, x any) {
	v.Driver().SetIndex(i, x)
}

// String returns the value v as a string.
// String is a special case because of Go's String method convention. Unlike the other getters,
// it does not panic if v's Type is not TypeString. Instead, it returns a string of the form "<T>"
// or "<T: V>" where T is v's type and V is a string representation of v's value.
func (v Value) String() string {
	return v.Driver().String()
}

// Truthy returns the JavaScript "truthiness" of the value v. In JavaScript,
// false, 0, "", null, undefined, and NaN are "falsy", and everything else is
// "truthy". See https://developer.mozilla.org/en-US/docs/Glossary/Truthy.
func (v Value) Truthy() bool {
	return v.Driver().Truthy()
}

// Type returns the JavaScript type of the value v. It is similar to JavaScript's typeof operator,
// except that it returns TypeNull instead of TypeObject for null.
func (v Value) Type() Type {
	return v.Driver().Type()
}

// A ValueError occurs when a Value method is invoked on
// a Value that does not support it. Such cases are documented
// in the description of each method.
type ValueError struct {
	Method string
	Type   Type
}

func (e *ValueError) Error() string {
	return "hypp/js: call of " + e.Method + " on " + e.Type.String()
}

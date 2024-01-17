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

func CopyBytesToGo(dst []byte, src Value) int {
	return d().CopyBytesToGo(dst, src)
}

func CopyBytesToJS(dst Value, src []byte) int {
	return d().CopyBytesToJS(dst, src)
}

func FuncOf(fn func(this Value, args []Value) any) Func {
	return d().FuncOf(fn)
}

func Global() Value {
	return d().Global()
}

func Null() Value {
	return d().Null()
}

func Undefined() Value {
	return d().Undefined()
}

func ValueOf(x any) Value {
	return d().ValueOf(x)
}

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

func (v Value) Bool() bool {
	return v.Driver().Bool()
}

func (v Value) Call(m string, args ...any) Value {
	return v.Driver().Call(m, args...)
}

func (v Value) Delete(p string) {
	v.Driver().Delete(p)
}

func (v Value) Equal(w Value) bool {
	return v.Driver().Equal(w)
}

func (v Value) Float() float64 {
	return v.Driver().Float()
}

func (v Value) Get(p string) Value {
	return v.Driver().Get(p)
}

func (v Value) Index(i int) Value {
	return v.Driver().Index(i)
}

func (v Value) InstanceOf(t Value) bool {
	return v.Driver().InstanceOf(t)
}

func (v Value) Int() int {
	return v.Driver().Int()
}

func (v Value) Invoke(args ...any) Value {
	return v.Driver().Invoke(args...)
}

func (v Value) IsNaN() bool {
	return v.Driver().IsNaN()
}

func (v Value) IsNull() bool {
	return v.Driver().IsNull()
}

func (v Value) IsUndefined() bool {
	return v.Driver().IsUndefined()
}

func (v Value) Length() int {
	return v.Driver().Length()
}

func (v Value) New(args ...any) Value {
	return v.Driver().New(args...)
}

func (v Value) Set(p string, x any) {
	v.Driver().Set(p, x)
}

func (v Value) SetIndex(i int, x any) {
	v.Driver().SetIndex(i, x)
}

func (v Value) String() string {
	return v.Driver().String()
}

func (v Value) Truthy() bool {
	return v.Driver().Truthy()
}

func (v Value) Type() Type {
	return v.Driver().Type()
}

type ValueError struct {
	Method string
	Type   Type
}

func (e *ValueError) Error() string {
	return "hypp/js: call of " + e.Method + " on " + e.Type.String()
}

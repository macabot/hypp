package hypp

type JavaScript interface {
	CopyBytesToGo(dst []byte, src Value) int
	CopyBytesToJS(dst Value, src []byte) int
	FuncOf(fn func(this Value, args []Value) any) Func
	Global() Value
	Null() Value
	Undefined() Value
	ValueOf(x any) Value
}

type Error interface {
	Value
	Error() string
}

type Func interface {
	Value
	Release()
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

// Value represents a JavaScript value.
// It is based on type syscall/js.Value.
// It allows you to use the JavaScript environment without the js/wasm build constraint.
// See https://pkg.go.dev/syscall/js#Value for the method definitions.
type Value interface {
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

type ValueError struct {
	Method string
	Type   Type
}

func (e *ValueError) Error() string {
	return "hypp: call of " + e.Method + " on " + e.Type.String()
}

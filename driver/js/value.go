package js

import (
	"syscall/js"

	"github.com/macabot/hypp"
)

type JavaScript struct{}

var _ hypp.JavaScript = JavaScript{}

func (j JavaScript) CopyBytesToGo(dst []byte, src hypp.Value) int {
	return js.CopyBytesToGo(dst, src.(Value).Value)
}

func (j JavaScript) CopyBytesToJS(dst hypp.Value, src []byte) int {
	return js.CopyBytesToJS(dst.(Value).Value, src)
}

func (j JavaScript) FuncOf(fn func(hypp.Value, []hypp.Value) interface{}) hypp.Func {
	return Func{js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := make([]hypp.Value, len(args))
		for i, a := range args {
			input[i] = Value{a}
		}
		return fn(Value{this}, input)
	})}
}

func (j JavaScript) Global() hypp.Value {
	return Value{js.Global()}
}

func (j JavaScript) Null() hypp.Value {
	return Value{js.Null()}
}

func (j JavaScript) Undefined() hypp.Value {
	return Value{js.Undefined()}
}

func (j JavaScript) ValueOf(x any) hypp.Value {
	if v, ok := x.(hypp.Value); ok {
		return v
	}
	return Value{js.ValueOf(x)}
}

// Err is an alias for js.Error.
// When directly embedding js.Error in type Error we cannot have a method named Error: the name would already be taken by the embedded value.
// Instead we embed the alias Err.
type Err js.Error

type Error struct {
	Err
}

var _ hypp.Error = Error{}

func (e Error) Call(m string, args ...any) hypp.Value {
	return Value{e.Err.Value}.Call(m, args...)
}

func (e Error) Equal(w hypp.Value) bool {
	return Value{e.Err.Value}.Equal(w)
}

func (e Error) Get(p string) hypp.Value {
	return Value{e.Err.Value}.Get(p)
}

func (e Error) Index(i int) hypp.Value {
	return Value{e.Err.Value}.Index(i)
}

func (e Error) InstanceOf(t hypp.Value) bool {
	return Value{e.Err.Value}.InstanceOf(t)
}

func (e Error) Invoke(args ...any) hypp.Value {
	return Value{e.Err.Value}.Invoke(args...)
}

func (e Error) New(args ...any) hypp.Value {
	return Value{e.Err.Value}.New(args...)
}

func (e Error) Type() hypp.Type {
	return Value{e.Err.Value}.Type()
}

func (e Error) Error() string {
	return js.Error(e.Err).Error()
}

type Func struct {
	js.Func
}

var _ hypp.Func = Func{}

func (f Func) Call(m string, args ...any) hypp.Value {
	return Value{f.Func.Value}.Call(m, args...)
}

func (f Func) Equal(w hypp.Value) bool {
	return Value{f.Func.Value}.Equal(w)
}

func (f Func) Get(p string) hypp.Value {
	return Value{f.Func.Value}.Get(p)
}

func (f Func) Index(i int) hypp.Value {
	return Value{f.Func.Value}.Index(i)
}

func (f Func) InstanceOf(t hypp.Value) bool {
	return Value{f.Func.Value}.InstanceOf(t)
}

func (f Func) Invoke(args ...any) hypp.Value {
	return Value{f.Func.Value}.Invoke(args...)
}

func (f Func) New(args ...any) hypp.Value {
	return Value{f.Func.Value}.New(args...)
}

func (f Func) Type() hypp.Type {
	return Value{f.Func.Value}.Type()
}

func (f Func) Release() {
	f.Func.Release()
}

type Value struct {
	js.Value
}

var _ hypp.Value = Value{}

func convertArg(arg any) any {
	switch v := arg.(type) {
	case hypp.Error:
		return v.(Error).Error
	case hypp.Func:
		return v.(Func).Func
	case hypp.Value:
		return v.(Value).Value
	case hypp.ValueError:
		return js.ValueError{
			Method: v.Method,
			Type:   js.Type(v.Type),
		}
	case []interface{}:
		l := make([]interface{}, len(v))
		for i, x := range v {
			l[i] = convertArg(x)
		}
		return l
	case map[string]interface{}:
		m := make(map[string]interface{}, len(v))
		for k, v := range v {
			m[k] = convertArg(v)
		}
		return m
	default:
		return arg
	}
}

func (v Value) Call(m string, args ...any) hypp.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return Value{v.Value.Call(m, converted...)}
}

func (v Value) Equal(w hypp.Value) bool {
	return v.Value.Equal(w.(Value).Value)
}

func (v Value) Get(p string) hypp.Value {
	return Value{v.Value.Get(p)}
}

func (v Value) Index(i int) hypp.Value {
	return Value{v.Value.Index(i)}
}

func (v Value) InstanceOf(t hypp.Value) bool {
	return v.Value.InstanceOf(t.(Value).Value)
}

func (v Value) Invoke(args ...any) hypp.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return Value{v.Value.Invoke(converted...)}
}

func (v Value) New(args ...any) hypp.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return Value{v.Value.New(converted...)}
}

func (v Value) Set(p string, x any) {
	v.Value.Set(p, convertArg(x))
}

func (v Value) SetIndex(i int, x any) {
	v.Value.SetIndex(i, convertArg(x))
}

func (v Value) Type() hypp.Type {
	return hypp.Type(v.Value.Type())
}

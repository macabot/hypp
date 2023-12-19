package jsd

import (
	"syscall/js"

	hyppjs "github.com/macabot/hypp/js"
)

func init() {
	hyppjs.Register(Driver{})
}

type Driver struct{}

var _ hyppjs.Driver = Driver{}

func (_ Driver) CopyBytesToGo(dst []byte, src hyppjs.Value) int {
	return js.CopyBytesToGo(dst, src.(Value).Value)
}

func (_ Driver) CopyBytesToJS(dst hyppjs.Value, src []byte) int {
	return js.CopyBytesToJS(dst.(Value).Value, src)
}

func (_ Driver) FuncOf(fn func(hyppjs.Value, []hyppjs.Value) interface{}) hyppjs.Func {
	return Func{js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := make([]hyppjs.Value, len(args))
		for i, a := range args {
			input[i] = Value{a}
		}
		return fn(Value{this}, input)
	})}
}

func (_ Driver) Global() hyppjs.Value {
	return Value{js.Global()}
}

func (_ Driver) Null() hyppjs.Value {
	return Value{js.Null()}
}

func (_ Driver) Undefined() hyppjs.Value {
	return Value{js.Undefined()}
}

func (_ Driver) ValueOf(x any) hyppjs.Value {
	return Value{js.ValueOf(convertArg(x))}
}

type Func struct {
	js.Func
}

var _ hyppjs.Func = Func{}

func (f Func) Call(m string, args ...any) hyppjs.Value {
	return Value{f.Func.Value}.Call(m, args...)
}

func (f Func) Equal(w hyppjs.Value) bool {
	return Value{f.Func.Value}.Equal(w)
}

func (f Func) Get(p string) hyppjs.Value {
	return Value{f.Func.Value}.Get(p)
}

func (f Func) Index(i int) hyppjs.Value {
	return Value{f.Func.Value}.Index(i)
}

func (f Func) InstanceOf(t hyppjs.Value) bool {
	return Value{f.Func.Value}.InstanceOf(t)
}

func (f Func) Invoke(args ...any) hyppjs.Value {
	return Value{f.Func.Value}.Invoke(args...)
}

func (f Func) New(args ...any) hyppjs.Value {
	return Value{f.Func.Value}.New(args...)
}

func (f Func) Type() hyppjs.Type {
	return Value{f.Func.Value}.Type()
}

func (f Func) Release() {
	f.Func.Release()
}

type Value struct {
	js.Value
}

var _ hyppjs.Value = Value{}

func convertArg(arg any) any {
	switch v := arg.(type) {
	case hyppjs.Error:
		return v.Value
	case hyppjs.Func:
		return v.(Func).Func
	case hyppjs.Value:
		return v.(Value).Value
	case hyppjs.ValueError:
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

func (v Value) Call(m string, args ...any) hyppjs.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return Value{v.Value.Call(m, converted...)}
}

func (v Value) Equal(w hyppjs.Value) bool {
	return v.Value.Equal(w.(Value).Value)
}

func (v Value) Get(p string) hyppjs.Value {
	return Value{v.Value.Get(p)}
}

func (v Value) Index(i int) hyppjs.Value {
	return Value{v.Value.Index(i)}
}

func (v Value) InstanceOf(t hyppjs.Value) bool {
	return v.Value.InstanceOf(t.(Value).Value)
}

func (v Value) Invoke(args ...any) hyppjs.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return Value{v.Value.Invoke(converted...)}
}

func (v Value) New(args ...any) hyppjs.Value {
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

func (v Value) Type() hyppjs.Type {
	return hyppjs.Type(v.Value.Type())
}

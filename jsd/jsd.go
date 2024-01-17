package jsd

import (
	"syscall/js"

	hyppjs "github.com/macabot/hypp/js"
)

func init() {
	hyppjs.Register(Driver{})
}

func convertArg(arg any) any {
	switch v := arg.(type) {
	case hyppjs.Error:
		return js.Error{Value: v.Value.Driver().(Value).Value}
	case hyppjs.Func:
		return v.Driver().(Func).Func
	case hyppjs.Value:
		return v.Driver().(Value).Value
	case hyppjs.ValueError:
		return js.ValueError{
			Method: v.Method,
			Type:   js.Type(v.Type),
		}
	case []any:
		l := make([]any, len(v))
		for i, x := range v {
			l[i] = convertArg(x)
		}
		return l
	case map[string]any:
		m := make(map[string]any, len(v))
		for k, v := range v {
			m[k] = convertArg(v)
		}
		return m
	default:
		return arg
	}
}

type Driver struct{}

var _ hyppjs.Driver = Driver{}

func (Driver) CopyBytesToGo(dst []byte, src hyppjs.Value) int {
	return js.CopyBytesToGo(dst, src.Driver().(Value).Value)
}

func (Driver) CopyBytesToJS(dst hyppjs.Value, src []byte) int {
	return js.CopyBytesToJS(dst.Driver().(Value).Value, src)
}

func (Driver) FuncOf(fn func(hyppjs.Value, []hyppjs.Value) any) hyppjs.Func {
	return hyppjs.MakeFunc(Func{js.FuncOf(func(this js.Value, args []js.Value) any {
		input := make([]hyppjs.Value, len(args))
		for i, a := range args {
			input[i] = hyppjs.MakeValue(Value{a})
		}
		return fn(hyppjs.MakeValue(Value{this}), input)
	})})
}

func (Driver) Global() hyppjs.Value {
	return hyppjs.MakeValue(Value{js.Global()})
}

func (Driver) Null() hyppjs.Value {
	return hyppjs.MakeValue(Value{js.Null()})
}

func (Driver) Undefined() hyppjs.Value {
	return hyppjs.MakeValue(Value{js.Undefined()})
}

func (Driver) ValueOf(x any) hyppjs.Value {
	return hyppjs.MakeValue(Value{js.ValueOf(convertArg(x))})
}

func (Driver) DefaultValueDriver() hyppjs.ValueDriver {
	return Value{}
}

func (Driver) DefaultFuncDriver() hyppjs.FuncDriver {
	return Func{}
}

type Func struct {
	js.Func
}

var _ hyppjs.FuncDriver = Func{}

func (f Func) ValueDriver() hyppjs.ValueDriver {
	return Value{f.Func.Value}
}

func (f Func) Release() {
	f.Func.Release()
}

type Value struct {
	js.Value
}

var _ hyppjs.ValueDriver = Value{}

func (v Value) Call(m string, args ...any) hyppjs.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return hyppjs.MakeValue(Value{v.Value.Call(m, converted...)})
}

func (v Value) Equal(w hyppjs.Value) bool {
	return v.Value.Equal(w.Driver().(Value).Value)
}

func (v Value) Get(p string) hyppjs.Value {
	return hyppjs.MakeValue(Value{v.Value.Get(p)})
}

func (v Value) Index(i int) hyppjs.Value {
	return hyppjs.MakeValue(Value{v.Value.Index(i)})
}

func (v Value) InstanceOf(t hyppjs.Value) bool {
	return v.Value.InstanceOf(t.Driver().(Value).Value)
}

func (v Value) Invoke(args ...any) hyppjs.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return hyppjs.MakeValue(Value{v.Value.Invoke(converted...)})
}

func (v Value) New(args ...any) hyppjs.Value {
	converted := make([]any, len(args))
	for i, a := range args {
		converted[i] = convertArg(a)
	}
	return hyppjs.MakeValue(Value{v.Value.New(converted...)})
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

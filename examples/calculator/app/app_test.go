package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputer(t *testing.T) {
	assert.Equal(t, 2.0, computer["+"](1, 1))
	assert.Equal(t, 2.0, computer["-"](3, 1))
	assert.Equal(t, 2.0, computer["×"](2, 1))
	assert.Equal(t, 2.0, computer["÷"](4, 2))
}

func TestNewDigit(t *testing.T) {
	assert.Equal(t, &MyState{
		hasCarry: false,
		value:    23,
	}, newDigit(&MyState{hasCarry: false, value: 2}, 3.0))

	assert.Equal(t, &MyState{
		hasCarry: false,
		value:    3,
	}, newDigit(&MyState{hasCarry: true, value: 2}, 3.0))
}

func TestNewFn(t *testing.T) {
	assert.Equal(t, &MyState{
		fn:       "+",
		hasCarry: true,
		carry:    3,
		value:    3,
	}, newFn(&MyState{
		value:    3,
		hasCarry: true,
	}, "+"))

	assert.Equal(t, &MyState{
		fn:       "-",
		hasCarry: true,
		carry:    2,
		value:    5,
	}, newFn(&MyState{
		fn:       "+",
		hasCarry: false,
		carry:    3,
		value:    2,
	}, "-"))
}

func TestEqual(t *testing.T) {
	assert.Equal(t, &MyState{
		fn:       "",
		hasCarry: true,
		carry:    3,
		value:    2,
	}, equal(&MyState{
		fn:       "",
		hasCarry: true,
		carry:    3,
		value:    2,
	}, nil))

	assert.Equal(t, &MyState{
		fn:       "-",
		hasCarry: true,
		carry:    3,
		value:    2,
	}, equal(&MyState{
		fn:       "-",
		hasCarry: false,
		carry:    5,
		value:    3,
	}, nil))

	assert.Equal(t, &MyState{
		fn:       "-",
		hasCarry: true,
		carry:    5,
		value:    -2,
	}, equal(&MyState{
		fn:       "-",
		hasCarry: true,
		carry:    5,
		value:    3,
	}, nil))
}

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
	assert.Equal(t, &State{
		hasCarry: false,
		value:    23,
	}, newDigit(&State{hasCarry: false, value: 2}, 3.0))

	assert.Equal(t, &State{
		hasCarry: false,
		value:    3,
	}, newDigit(&State{hasCarry: true, value: 2}, 3.0))
}

func TestNewFn(t *testing.T) {
	assert.Equal(t, &State{
		fn:       "+",
		hasCarry: true,
		carry:    3,
		value:    3,
	}, newFn(&State{
		value:    3,
		hasCarry: true,
	}, "+"))

	assert.Equal(t, &State{
		fn:       "-",
		hasCarry: true,
		carry:    2,
		value:    5,
	}, newFn(&State{
		fn:       "+",
		hasCarry: false,
		carry:    3,
		value:    2,
	}, "-"))
}

func TestEqual(t *testing.T) {
	assert.Equal(t, &State{
		fn:       "",
		hasCarry: true,
		carry:    3,
		value:    2,
	}, equal(&State{
		fn:       "",
		hasCarry: true,
		carry:    3,
		value:    2,
	}, nil))

	assert.Equal(t, &State{
		fn:       "-",
		hasCarry: true,
		carry:    3,
		value:    2,
	}, equal(&State{
		fn:       "-",
		hasCarry: false,
		carry:    5,
		value:    3,
	}, nil))

	assert.Equal(t, &State{
		fn:       "-",
		hasCarry: true,
		carry:    5,
		value:    -2,
	}, equal(&State{
		fn:       "-",
		hasCarry: true,
		carry:    5,
		value:    3,
	}, nil))
}

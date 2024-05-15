package hypp

import (
	"testing"

	mocks "github.com/macabot/hypp/internal/mocks/hypp/js"
	"github.com/macabot/hypp/js"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:generate mockery

// When creating AppProps, the state is initialized to its default value.
// The dispatch function only renders if the state has changed.
// This test ensures that the render function is called at least once, even if AppProps.Init resolves to the state's default value.
//
// E.g. this allows the state to be an empty struct.
// The default value of 'struct{}' is 'struct{}{}'. Since it has no fields, it will never change. Nevertheless, we should expect the render function to be called at least once.
func TestUpdateWillRenderAtLeastOnce(t *testing.T) {
	driver := mocks.NewMockDriver(t)
	js.Register(driver)

	global := mocks.NewMockValueDriver(t)
	globalValue := js.MakeValue(global)
	driver.EXPECT().Global().Return(globalValue)
	jsFunc := mocks.NewMockFuncDriver(t)
	jsFuncValue := mocks.NewMockValueDriver(t)
	jsFunc.EXPECT().ValueDriver().Return(jsFuncValue)
	driver.EXPECT().FuncOf(mock.Anything).Return(js.MakeFunc(jsFunc))
	requestID := mocks.NewMockValueDriver(t)
	global.EXPECT().Call("requestAnimationFrame", mock.Anything).Return(js.MakeValue(requestID))
	requestID.EXPECT().Int().Return(1)

	// The String() method is called due to the testify library.
	// When comparing the arguments, it calls fmt.Sprintf.
	// See https://github.com/stretchr/testify/blob/v1.9.0/mock/mock.go#L939.
	// Sprintf checks if jsFuncValue implements the Stringer interface, which it does, and then calls it.
	// Therefore we need to mock the call to String() as well.
	jsFuncValue.EXPECT().String().Return("test func")

	appProps := &AppProps[struct{}]{
		View:               func(state struct{}) *VNode { return nil },
		hasRequestedRender: false, // This is the default value when initiating AppProps.
	}
	update(appProps, struct{}{})
	assert.True(t, appProps.hasRequestedRender)
}

package hypp

import (
	"testing"

	mocks "github.com/macabot/hypp/internal/mocks/hypp/js"
	"github.com/macabot/hypp/js"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// When creating AppProps, the state is initialized to its default value.
// The dispatch function only renders if the state has changed.
// This test ensures that the render function is called at least once, even if AppProps.Init resolves to the state's default value.
//
// E.g. this allows the state to be an empty struct.
// The default value of 'struct{}' is 'struct{}{}'. Since it has no fields, it will never change. Nevertheless, we should expect the render function to be called at least once.
func TestUpdateWillRenderIfNeverRenderedBefore(t *testing.T) {
	mockDriver := mocks.NewMockDriver(t)
	js.Register(mockDriver)

	global := mocks.NewMockValueDriver(t)
	globalValue := js.MakeValue(global)
	mockDriver.EXPECT().Global().Return(globalValue)
	mockDriver.EXPECT().FuncOf(mock.Anything).Return(js.Func{})
	requestID := mocks.NewMockValueDriver(t)
	global.EXPECT().Call("requestAnimationFrame", mock.Anything).Return(js.MakeValue(requestID))
	requestID.EXPECT().Int().Return(1)

	// TODO why is this needed?
	defaultValue := mocks.NewMockValueDriver(t)
	mockDriver.EXPECT().DefaultValueDriver().Return(js.MakeValue(defaultValue))
	defaultValue.EXPECT().String().Return("DEBUG")

	appProps := &AppProps[struct{}]{
		View:               func(state struct{}) *VNode { return nil },
		hasRequestedRender: false, // This is the default value when initiating AppProps.
	}
	update(appProps, struct{}{})
	assert.True(t, appProps.hasRequestedRender)
}

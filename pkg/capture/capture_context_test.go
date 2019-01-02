package capture_test

import (
	"context"
	"testing"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/capture"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestContextManagerGetCaptureOK(t *testing.T) {
	ctx := context.Background()

	capt := pkg.Capture{ID: kallax.NewULID()}
	ctx = capture.WithCapture(ctx, &capt)

	capt2, err := capture.GetFromContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, capt.ID, capt2.ID)
}

func TestContextManagerGetCaptureMissingCapture(t *testing.T) {
	ctx := context.Background()

	_, err := capture.GetFromContext(ctx)
	assert.EqualError(t, err, "capture not found in context")
}
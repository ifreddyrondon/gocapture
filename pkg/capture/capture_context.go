package capture

import (
	"context"
	"errors"

	"github.com/ifreddyrondon/capture/pkg"
)

type ctxKey string

const (
	captureKey ctxKey = "capture"
)

var (
	errMissingCapture    = errors.New("capture not found in context")
	errWrongCaptureValue = errors.New("capture value set incorrectly in context")
)

// WithCapture will return a new context with the capture value added to it.
func WithCapture(ctx context.Context, capt *pkg.Capture) context.Context {
	return context.WithValue(ctx, captureKey, capt)
}

// GetFromContext will return the capture assigned to the context,
// or nil if there is any error or there isn't a capture.
func GetFromContext(ctx context.Context) (*pkg.Capture, error) {
	tmp := ctx.Value(captureKey)
	if tmp == nil {
		return nil, errMissingCapture
	}
	capt, ok := tmp.(*pkg.Capture)
	if !ok {
		return nil, errWrongCaptureValue
	}
	return capt, nil
}
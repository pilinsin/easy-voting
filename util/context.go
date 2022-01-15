package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewContext() context.Context {
	return context.Background()
}
func CancelContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
func CancelTimerContext(t time.Duration) (context.Context, context.CancelFunc){
	return context.WithTimeout(context.Background(), t)
}
func SignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
func WithSignal(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
}

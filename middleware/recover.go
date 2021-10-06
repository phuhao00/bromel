package middleware

import (
	"context"
	"errors"
	"github.com/phuhao00/bromel/logger"
	"go.uber.org/zap"
	"runtime"
)

// HandlerFunc is recovery handler func.
type HandlerFunc func(ctx context.Context, req, err interface{}) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
	logger  logger.Logger
}

// WithHandler with recovery handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithLogger with recovery logger.
func WithLogger(logger logger.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Recovery is a server middleware that recovers from any panics.
func Recovery(opts ...Option) Middleware {
	options := options{
		handler: func(ctx context.Context, req, err interface{}) error {
			return errors.New("RECOVERY" + zap.Any("err", err).String)
		},
	}
	for _, o := range opts {
		o(&options)
	}
	logger := options.logger
	return func(handler Handler) Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					logger.Error("Recovery", zap.Any("err", rerr))

					err = options.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}

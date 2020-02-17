package util

import (
	"context"
	"fmt"

	"logur.dev/logur"
)

type ErrorHandler interface {
	Handle(err error)
	HandleContext(ctx context.Context, err error)
}

type AppErrorHandler struct {
	Logger logur.LoggerFacade
}

func (a AppErrorHandler) Handle(err error) {
	if err != nil {
		a.Logger.Error(fmt.Sprintf("Error: %v", err))
	}
}
func (AppErrorHandler) HandleContext(_ context.Context, _ error) {}

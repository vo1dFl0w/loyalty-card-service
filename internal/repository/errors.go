package repository

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNotFound                     = errors.New("not found")
	ErrClientClosedRequest          = errors.New("client closed request")
	ErrGatewayTimeout               = errors.New("gateway timeout")
	ErrCurrentBalanceLessThanAmount = errors.New("current balance less than amount")
	ErrLoyaltyCardAlreadyBlocked    = errors.New("loyalty card already blocked")
	ErrNoRowsDeleted                = errors.New("no rows deleted")
)

func MapContextOnly(err error, ctx context.Context) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrGatewayTimeout
	}
	if errors.Is(err, context.Canceled) {
		return ErrClientClosedRequest
	}

	if ctx != nil {
		if errCtx := ctx.Err(); errCtx != nil {
			if errors.Is(errCtx, context.DeadlineExceeded) {
				return ErrGatewayTimeout
			}
			if errors.Is(errCtx, context.Canceled) {
				return ErrClientClosedRequest
			}
		}
	}

	return nil
}

func WrapError(prefix string, err error, ctx context.Context) error {
	if err == nil {
		if ctxErr := MapContextOnly(nil, ctx); ctxErr != nil {
			return ctxErr
		}
		return nil
	}
	if ctxErr := MapContextOnly(err, ctx); ctxErr != nil {
		return ctxErr
	}
	return fmt.Errorf("%s: %w", prefix, err)
}

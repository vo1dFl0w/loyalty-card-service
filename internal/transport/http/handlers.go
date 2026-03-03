package http

import (
	"context"
	"net/http"

	"github.com/vo1dFl0w/loyalty-card-service/internal/config"
	"github.com/vo1dFl0w/loyalty-card-service/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/loyalty-card-service/internal/usecase"
	"github.com/vo1dFl0w/loyalty-card-service/pkg/logger"
)

type Handler struct {
	cfg            *config.Config
	logger         logger.Logger
	loyaltyCardSrv usecase.LoyaltyCardService
}

func NewHandler(cfg *config.Config, logger logger.Logger, loyaltyCardSrv usecase.LoyaltyCardService) *Handler {
	return &Handler{cfg: cfg, logger: logger, loyaltyCardSrv: loyaltyCardSrv}

}
func (h *Handler) APIV1LoyaltyCreatePost(ctx context.Context, req *httpgen.CreateLoyaltyCardRequest) (httpgen.APIV1LoyaltyCreatePostRes, error) {
	res, err := h.loyaltyCardSrv.NewLoyaltyCard(ctx, req.UserID)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToLoyaltyCreatePostRes(), nil
	}

	return &httpgen.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   res.Balance,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (h *Handler) APIV1LoyaltyGet(ctx context.Context, params httpgen.APIV1LoyaltyGetParams) (httpgen.APIV1LoyaltyGetRes, error) {
	res, err := h.loyaltyCardSrv.GetLoyaltyCardByUserID(ctx, params.UserID.Value)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToLoyaltyGetRes(), nil
	}

	return &httpgen.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   res.Balance,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (h *Handler) APIV1LoyaltyBalancePatch(ctx context.Context, req *httpgen.UpdateBalanceRequest) (httpgen.APIV1LoyaltyBalancePatchRes, error) {
	res, err := h.loyaltyCardSrv.ChangeLoyaltyCardBalance(ctx, req.UserID, req.Amount, string(req.Mode))
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToLoayltyBalancePatchRes(), nil
	}

	return &httpgen.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   res.Balance,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (h *Handler) APIV1LoyaltyIsBlockedPatch(ctx context.Context, req *httpgen.UpdateIsBlockRequest) (httpgen.APIV1LoyaltyIsBlockedPatchRes, error) {
	res, err := h.loyaltyCardSrv.ChangeLoyaltyCardIsBlocked(ctx, req.UserID, req.IsBlocked)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToLoyaltyIsBlockedPatchRes(), nil
	}

	return &httpgen.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   res.Balance,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (h *Handler) APIV1LoyaltyDeleteDelete(ctx context.Context, req *httpgen.DeleteLoyaltyCardRequest) (httpgen.APIV1LoyaltyDeleteDeleteRes, error) {
	_, err := h.loyaltyCardSrv.DeleteLoyaltyCard(ctx, req.UserID)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToLoyaltyDeleteRes(), nil
	}

	return &httpgen.APIV1LoyaltyDeleteDeleteNoContent{}, nil
}

func (h *Handler) LogHTTPError(ctx context.Context, err error, httpErr *HTTPError) {
	attrs := []any{
		"error", err,
		"status", httpErr.Status,
		"message", httpErr.Message,
	}

	switch {
	case httpErr.Status >= 500:
		switch httpErr.Status {
		case http.StatusGatewayTimeout:
			h.logger.Error("http_request_failed", append(attrs, "reason", "dependency_timeout")...)
		default:
			h.logger.Error("http_request_failed", append(attrs, "reason", "internal_server_error")...)
		}
	case httpErr.Status >= 400:
		h.logger.Warn("http_request_failed", append(attrs, "reason", "client_error")...)
	}
}

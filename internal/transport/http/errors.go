package http

import (
	"errors"
	"net/http"

	"github.com/vo1dFl0w/loyalty-card-service/internal/domain"
	"github.com/vo1dFl0w/loyalty-card-service/internal/transport/http/httpgen"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrClientClosedRequest = errors.New("client closed request")
	ErrGatewayTimeout      = errors.New("gateway timeout")
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("not found")
)

type HTTPError struct {
	Message string
	Status  int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) ToLoyaltyCreatePostRes() httpgen.APIV1LoyaltyCreatePostRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.APIV1LoyaltyCreatePostBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1LoyaltyCreatePostGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.APIV1LoyaltyCreatePostInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToLoyaltyGetRes() httpgen.APIV1LoyaltyGetRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.APIV1LoyaltyGetBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.APIV1LoyaltyGetNotFound{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1LoyaltyGetGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.APIV1LoyaltyGetInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToLoayltyBalancePatchRes() httpgen.APIV1LoyaltyBalancePatchRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.APIV1LoyaltyBalancePatchBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.APIV1LoyaltyBalancePatchNotFound{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1LoyaltyBalancePatchGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.APIV1LoyaltyBalancePatchInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToLoyaltyIsBlockedPatchRes() httpgen.APIV1LoyaltyIsBlockedPatchRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.APIV1LoyaltyIsBlockedPatchBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.APIV1LoyaltyIsBlockedPatchNotFound{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1LoyaltyIsBlockedPatchGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.APIV1LoyaltyIsBlockedPatchInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToLoyaltyDeleteRes() httpgen.APIV1LoyaltyDeleteDeleteRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.APIV1LoyaltyDeleteDeleteBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.APIV1LoyaltyDeleteDeleteNotFound{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1LoyaltyDeleteDeleteGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.APIV1LoyaltyDeleteDeleteInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func MapError(err error) *HTTPError {
	switch {
	case errors.Is(err, domain.ErrInvalidMode) || errors.Is(err, domain.ErrInvalidAmount) ||
		errors.Is(err, domain.ErrCurrentBalanceLessThanAmount) || errors.Is(err, domain.ErrLoyaltyCardAlreadyBlocked):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrNotFound):
		return &HTTPError{Message: ErrNotFound.Error(), Status: http.StatusNotFound}
	case errors.Is(err, domain.ErrGatewayTimeout):
		return &HTTPError{Message: ErrGatewayTimeout.Error(), Status: http.StatusGatewayTimeout}
	default:
		return &HTTPError{Message: ErrInternalServerError.Error(), Status: http.StatusInternalServerError}
	}
}

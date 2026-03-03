package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/loyalty-card-service/internal/domain"
	"github.com/vo1dFl0w/loyalty-card-service/internal/repository"
)

type LoyaltyCardService interface {
	NewLoyaltyCard(ctx context.Context, userID uuid.UUID) (loyaltyCard *domain.LoyaltyCard, err error)
	GetLoyaltyCardByUserID(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error)
	ChangeLoyaltyCardBalance(ctx context.Context, userID uuid.UUID, amount float64, mode string) (loyaltyCard *domain.LoyaltyCard, err error)
	ChangeLoyaltyCardIsBlocked(ctx context.Context, userID uuid.UUID, isBlocked bool) (loyaltyCard *domain.LoyaltyCard, err error)
	DeleteLoyaltyCard(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error)
}

type loyaltyCardSvc struct {
	loyaltyCardRepo repository.LoyaltyCardRepository
}

func NewLoyaltyCardService(loyaltyCardRepo repository.LoyaltyCardRepository) *loyaltyCardSvc {
	return &loyaltyCardSvc{loyaltyCardRepo: loyaltyCardRepo}
}

func (l *loyaltyCardSvc) NewLoyaltyCard(ctx context.Context, userID uuid.UUID) (loyaltyCard *domain.LoyaltyCard, err error) {
	res, err := l.loyaltyCardRepo.Create(ctx, userID)
	if err != nil {
		return nil, domain.WrapError("create", err, ctx)
	}

	return res, nil
}

func (l *loyaltyCardSvc) GetLoyaltyCardByUserID(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error) {
	res, err := l.loyaltyCardRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.WrapError("create", err, ctx)
	}

	return res, nil
}

func (l *loyaltyCardSvc) ChangeLoyaltyCardBalance(ctx context.Context, userID uuid.UUID, amount float64, mode string) (loyaltyCard *domain.LoyaltyCard, err error) {
	if mode != "add" && mode != "withdraw" {
		return nil, domain.ErrInvalidMode
	}

	if amount <= 0.0 {
		return nil, domain.ErrInvalidAmount
	}

	res, err := l.loyaltyCardRepo.UpdateBalance(ctx, userID, amount, mode)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		} else if errors.Is(err, repository.ErrCurrentBalanceLessThanAmount) {
			return nil, domain.ErrCurrentBalanceLessThanAmount
		}
		return nil, domain.WrapError("update balance", err, ctx)
	}

	return res, nil
}

func (l *loyaltyCardSvc) ChangeLoyaltyCardIsBlocked(ctx context.Context, userID uuid.UUID, isBlocked bool) (loyaltyCard *domain.LoyaltyCard, err error) {
	res, err := l.loyaltyCardRepo.UpdateIsBlocked(ctx, userID, isBlocked)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		} else if errors.Is(err, repository.ErrLoyaltyCardAlreadyBlocked) {
			return nil, domain.ErrLoyaltyCardAlreadyBlocked
		}
		return nil, domain.WrapError("update is blocked", err, ctx)
	}

	return res, nil
}

func (l *loyaltyCardSvc) DeleteLoyaltyCard(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error) {
	res, err := l.loyaltyCardRepo.Delete(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRowsDeleted) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.WrapError("delete", err, ctx)
	}

	return res, nil
}

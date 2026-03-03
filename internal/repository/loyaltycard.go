package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/loyalty-card-service/internal/domain"
)

type LoyaltyCardRepository interface {
	Create(ctx context.Context, userID uuid.UUID) (loyaltyCard *domain.LoyaltyCard, err error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, amount float64, mode string) (loyaltyCard *domain.LoyaltyCard, err error)
	UpdateIsBlocked(ctx context.Context, userID uuid.UUID, isBlocked bool) (loyaltyCard *domain.LoyaltyCard, err error)
	Delete(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error)
}

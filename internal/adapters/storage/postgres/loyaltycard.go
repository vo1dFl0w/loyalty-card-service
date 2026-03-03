package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/loyalty-card-service/internal/adapters/storage/postgres/pggen"
	"github.com/vo1dFl0w/loyalty-card-service/internal/domain"
	"github.com/vo1dFl0w/loyalty-card-service/internal/repository"
)

type loyaltyCard struct {
	db      *sql.DB
	queries *pggen.Queries
}

func NewLoyaltyCard(db *sql.DB, q *pggen.Queries) *loyaltyCard {
	return &loyaltyCard{db: db, queries: q}
}

func (lc *loyaltyCard) Create(ctx context.Context, userID uuid.UUID) (loyaltyCard *domain.LoyaltyCard, err error) {
	tx, err := lc.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, repository.WrapError("begin transaction", err, ctx)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		commitErr := tx.Commit()
		if commitErr != nil {
			err = repository.WrapError("commit", err, ctx)
		}
	}()

	qtx := lc.queries.WithTx(tx)

	res, err := qtx.Create(ctx, userID)
	if err != nil {
		return nil, repository.WrapError("create", err, ctx)
	}

	points, err := ParseFromStringToFloat64(res.Balance.String)
	if err != nil {
		return nil, repository.WrapError("parse from string to float64", err, ctx)
	}

	return &domain.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   points,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (lc *loyaltyCard) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error) {
	res, err := lc.queries.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError("find by user id", err, ctx)
	}

	points, err := ParseFromStringToFloat64(res.Balance.String)
	if err != nil {
		return nil, repository.WrapError("parse from string to float64", err, ctx)
	}

	return &domain.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   points,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (lc *loyaltyCard) UpdateBalance(ctx context.Context, userID uuid.UUID, amount float64, mode string) (loyaltyCard *domain.LoyaltyCard, err error) {
	tx, err := lc.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, repository.WrapError("begin transaction", err, ctx)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		commitErr := tx.Commit()
		if commitErr != nil {
			err = repository.WrapError("commit", err, ctx)
		}
	}()

	qtx := lc.queries.WithTx(tx)

	userCard, err := qtx.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError("find by user id", err, ctx)
	}

	userBalance, err := ParseFromStringToFloat64(userCard.Balance.String)
	if err != nil {
		return nil, repository.WrapError("parse from string to float64", err, ctx)
	}

	var newBalance string
	if mode == "withdraw" {
		if userBalance-amount < 0 {
			return nil, repository.ErrCurrentBalanceLessThanAmount
		}
		newBalance = ParseFloat64ToString(userBalance - amount)
	}
	if mode == "add" {
		newBalance = ParseFloat64ToString(userBalance + amount)
	}

	res, err := qtx.UpdateBalance(ctx, pggen.UpdateBalanceParams{
		Balance: sql.NullString{
			String: newBalance,
			Valid:  newBalance != "",
		},
		UserID: userID,
	})
	if err != nil {
		return nil, repository.WrapError("update balance", err, ctx)
	}

	points, err := ParseFromStringToFloat64(res.Balance.String)
	if err != nil {
		return nil, repository.WrapError("parse from string to float64", err, ctx)
	}

	return &domain.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   points,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (lc *loyaltyCard) UpdateIsBlocked(ctx context.Context, userID uuid.UUID, isBlocked bool) (loyaltyCard *domain.LoyaltyCard, err error) {
	tx, err := lc.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, repository.WrapError("begin transaction", err, ctx)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		commitErr := tx.Commit()
		if commitErr != nil {
			err = repository.WrapError("commit", err, ctx)
		}
	}()

	qtx := lc.queries.WithTx(tx)

	userCard, err := qtx.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError("find by user id", err, ctx)
	}

	if isBlocked == true && userCard.IsBlocked == true {
		return nil, repository.ErrLoyaltyCardAlreadyBlocked
	}

	res, err := qtx.UpdateIsBlocked(ctx, pggen.UpdateIsBlockedParams{
		IsBlocked: isBlocked,
		UserID:    userID,
	})
	if err != nil {
		return nil, repository.WrapError("update is blocked", err, ctx)
	}

	points, err := ParseFromStringToFloat64(res.Balance.String)
	if err != nil {
		return nil, repository.WrapError("parse from string to float64", err, ctx)
	}

	return &domain.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   points,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

func (lc *loyaltyCard) Delete(ctx context.Context, userID uuid.UUID) (*domain.LoyaltyCard, error) {
	res, err := lc.queries.Delete(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRowsDeleted
		}
		return nil, repository.WrapError("delete", err, ctx)
	}

	points, err := ParseFromStringToFloat64(res.Balance.String)
	if err != nil {
		return nil, repository.WrapError("parse from string to float64", err, ctx)
	}

	return &domain.LoyaltyCard{
		ID:        res.ID,
		UserID:    res.UserID,
		Balance:   points,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		IsBlocked: res.IsBlocked,
	}, nil
}

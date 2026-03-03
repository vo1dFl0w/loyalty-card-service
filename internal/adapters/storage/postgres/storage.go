package postgres

import (
	"database/sql"

	"github.com/vo1dFl0w/loyalty-card-service/internal/adapters/storage/postgres/pggen"
	"github.com/vo1dFl0w/loyalty-card-service/internal/repository"
)

type Storage struct {
	db              *sql.DB
	loyaltyCardRepo repository.LoyaltyCardRepository
}

func NewStorage(db *sql.DB) *Storage {
	q := pggen.New(db)

	return &Storage{
		db: db,
		loyaltyCardRepo: NewLoyaltyCard(db, q),
	}
}

func (s *Storage) LoyaltyCardRepo() repository.LoyaltyCardRepository {
	return s.loyaltyCardRepo
}

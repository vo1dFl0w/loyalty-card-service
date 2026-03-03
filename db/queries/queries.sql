-- name: Create :one
INSERT INTO loyalty_cards (user_id)
VALUES ($1)
RETURNING id, user_id, balance, created_at, updated_at, is_blocked;

-- name: FindByUserID :one
SELECT id, user_id, balance, created_at, updated_at, is_blocked
FROM loyalty_cards
WHERE user_id = $1;

-- name: UpdateBalance :one
UPDATE loyalty_cards
SET balance = $1, updated_at = NOW()
WHERE user_id = $2
RETURNING id, user_id, balance, created_at, updated_at, is_blocked;

-- name: UpdateIsBlocked :one
UPDATE loyalty_cards
SET is_blocked = $1, updated_at = NOW()
WHERE user_id = $2
RETURNING id, user_id, balance, created_at, updated_at, is_blocked;

-- name: Delete :one
DELETE FROM loyalty_cards
WHERE user_id = $1
RETURNING id, user_id, balance, created_at, updated_at, is_blocked;
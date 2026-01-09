package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Challenge struct {
	ID          int64
	ChallengeID uuid.UUID
	DeviceID    uuid.UUID
	Challenge   string
	ExpiresAt   time.Time
	Used        bool
	CreatedAt   time.Time
}

func CreateChallenge(ctx context.Context, c *Challenge) error {
	query := `
		INSERT INTO challenges (challenge_id, device_id, challenge, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return Pool.QueryRow(
		ctx,
		query,
		c.ChallengeID,
		c.DeviceID,
		c.Challenge,
		c.ExpiresAt,
	).Scan(&c.ID, &c.CreatedAt)
}

func GetValidChallenge(ctx context.Context, challengeID uuid.UUID) (*Challenge, error) {
	query := `
		SELECT id, challenge_id, device_id, challenge, expires_at, used, created_at
		FROM challenges
		WHERE challenge_id = $1
		  AND used = false
		  AND expires_at > now()
	`

	var c Challenge
	err := Pool.QueryRow(ctx, query, challengeID).Scan(
		&c.ID,
		&c.ChallengeID,
		&c.DeviceID,
		&c.Challenge,
		&c.ExpiresAt,
		&c.Used,
		&c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func MarkChallengeUsed(ctx context.Context, challengeID uuid.UUID) error {
	_, err := Pool.Exec(
		ctx,
		`UPDATE challenges SET used = true WHERE challenge_id = $1`,
		challengeID,
	)
	return err
}

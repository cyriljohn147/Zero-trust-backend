package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"github.com/cyriljohn147/zero-trust-backend/internal/db"
)

func GenerateChallenge(
	ctx context.Context,
	deviceID uuid.UUID,
) (*db.Challenge, error) {

	randomBytes := make([]byte, 32)
	_, _ = rand.Read(randomBytes)

	challenge := &db.Challenge{
		ChallengeID: uuid.New(),
		DeviceID:    deviceID,
		Challenge:   base64.StdEncoding.EncodeToString(randomBytes),
		ExpiresAt:   time.Now().Add(2 * time.Minute),
	}

	if err := db.CreateChallenge(ctx, challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}

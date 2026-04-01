package services

import (
	"context"

	"github.com/google/uuid"

	"errors"

	"github.com/cyriljohn147/zero-trust-backend/internal/auth"
	"github.com/cyriljohn147/zero-trust-backend/internal/crypto"
	"github.com/cyriljohn147/zero-trust-backend/internal/db"
)

var (
	ErrDeviceRevoked = errors.New("device revoked")
	ErrHighRisk      = errors.New("high risk")
)

func VerifyChallenge(
	ctx context.Context,
	challengeID uuid.UUID,
	signature string,
	ip string,
) (string, error) {

	challenge, err := db.GetValidChallenge(ctx, challengeID)
	if err != nil {
		return "", err
	}

	device, err := db.GetDeviceByDeviceID(ctx, challenge.DeviceID)
	if err != nil {
		return "", err
	}

	if device.Status == "revoked" {
		return "", ErrDeviceRevoked
	}

	// 🔥 RISK CALCULATION
	risk := CalculateRisk(ip, device.LastSeen)

	// 🔥 RISK DECISION
	if risk >= 50 {
		return "", ErrHighRisk
	}

	if err := crypto.VerifySignature(
		device.PublicKey,
		challenge.Challenge,
		signature,
	); err != nil {
		return "", err
	}

	_ = db.MarkChallengeUsed(ctx, challengeID)
	_ = db.UpdateLastSeen(ctx, device.DeviceID)

	return auth.GenerateToken(
		device.DeviceID.String(),
		device.UserID,
	)
}

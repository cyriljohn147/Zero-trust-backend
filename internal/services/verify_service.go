package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/cyriljohn147/zero-trust-backend/internal/auth"
	"github.com/cyriljohn147/zero-trust-backend/internal/crypto"
	"github.com/cyriljohn147/zero-trust-backend/internal/db"
)

func VerifyChallenge(
	ctx context.Context,
	challengeID uuid.UUID,
	signature string,
) (string, error) {

	challenge, err := db.GetValidChallenge(ctx, challengeID)
	if err != nil {
		return "", err
	}

	device, err := db.GetDeviceByDeviceID(ctx, challenge.DeviceID)
	if err != nil || device.Status != "active" {
		return "", err
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

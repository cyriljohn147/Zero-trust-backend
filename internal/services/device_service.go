package services

import (
	"context"

	"github.com/cyriljohn147/zero-trust-backend/internal/db"
	"github.com/google/uuid"
)

func RegisterDevice(
	ctx context.Context,
	userID int64,
	publicKey string,
) (*db.Device, error) {

	device := &db.Device{
		DeviceID:  uuid.New(),
		UserID:    userID,
		PublicKey: publicKey,
		Status:    "active",
	}

	if err := db.CreateDevice(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

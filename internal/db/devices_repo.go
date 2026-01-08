package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID           int64
	DeviceID     uuid.UUID
	UserID       int64
	PublicKey    string
	Status       string
	RegisteredAt time.Time
	LastSeen     *time.Time
}

var ErrDeviceNotFound = errors.New("device not found")

// CreateDevice inserts a new device record
func CreateDevice(ctx context.Context, device *Device) error {
	query := `
		INSERT INTO devices (device_id, user_id, public_key, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id, registered_at
	`

	return Pool.QueryRow(
		ctx,
		query,
		device.DeviceID,
		device.UserID,
		device.PublicKey,
	).Scan(&device.ID, &device.RegisteredAt)
}

// GetDeviceByDeviceID fetches a device using device_id
func GetDeviceByDeviceID(ctx context.Context, deviceID uuid.UUID) (*Device, error) {
	query := `
		SELECT id, device_id, user_id, public_key, status, registered_at, last_seen
		FROM devices
		WHERE device_id = $1
	`

	var d Device
	err := Pool.QueryRow(ctx, query, deviceID).Scan(
		&d.ID,
		&d.DeviceID,
		&d.UserID,
		&d.PublicKey,
		&d.Status,
		&d.RegisteredAt,
		&d.LastSeen,
	)

	if err != nil {
		return nil, ErrDeviceNotFound
	}

	return &d, nil
}

// UpdateLastSeen updates last_seen timestamp
func UpdateLastSeen(ctx context.Context, deviceID uuid.UUID) error {
	query := `
		UPDATE devices
		SET last_seen = now()
		WHERE device_id = $1
	`

	cmd, err := Pool.Exec(ctx, query, deviceID)
	if err != nil || cmd.RowsAffected() == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

// RevokeDevice marks a device as revoked
func RevokeDevice(ctx context.Context, deviceID uuid.UUID) error {
	query := `
		UPDATE devices
		SET status = 'revoked'
		WHERE device_id = $1
	`

	cmd, err := Pool.Exec(ctx, query, deviceID)
	if err != nil || cmd.RowsAffected() == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

package db

import (
	"context"

	"github.com/google/uuid"
)

func CreateAuditLog(
	ctx context.Context,
	auditID uuid.UUID,
	userID *int64,
	deviceID *uuid.UUID,
	eventType string,
	eventStatus string,
	ip string,
	userAgent string,
) error {

	query := `
		INSERT INTO audit_logs (
			audit_id, user_id, device_id,
			event_type, event_status,
			ip_address, user_agent
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := Pool.Exec(
		ctx,
		query,
		auditID,
		userID,
		deviceID,
		eventType,
		eventStatus,
		ip,
		userAgent,
	)

	return err
}

package services

import (
	"time"
)

func CalculateRisk(ip string, lastSeen *time.Time) int {
	risk := 0

	// 1. New IP (simple assumption)
	if ip != "127.0.0.1" {
		risk += 30
	}

	// 2. Unusual time (night access)
	hour := time.Now().Hour()
	if hour < 6 || hour > 23 {
		risk += 20
	}

	// 3. Inactive device (not seen recently)
	if lastSeen != nil {
		if time.Since(*lastSeen).Hours() > 24 {
			risk += 20
		}
	}

	return risk
}

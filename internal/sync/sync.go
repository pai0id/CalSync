package sync

import (
	"CalSync/internal/alfa"
	"fmt"
)

func SyncCalendars(gCalService []byte, email, alfaApiKey string) error {
	token, err := alfa.GetAPIToken(email, alfaApiKey)
	if err != nil {
		return fmt.Errorf("Error at SyncCalendars: %w", err)
	}

	lessons, err := alfa.GetLessons(token)
	if err != nil {
		return fmt.Errorf("Error at SyncCalendars: %w", err)
	}

	regLessons, err := alfa.GetRegularLessons(token)
	if err != nil {
		return fmt.Errorf("Error at SyncCalendars: %w", err)
	}

	fmt.Println(lessons)
	fmt.Println(regLessons)
	return nil
}

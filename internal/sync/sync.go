package sync

import (
	"CalSync/internal/alfa"
	"CalSync/internal/gcal"
	"CalSync/internal/logic"
	"context"
	"fmt"
)

func SyncCalendars(gCalCreds []byte, email, alfaApiKey string) error {
	token, err := alfa.GetAPIToken(email, alfaApiKey)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	err = alfa.UpdateCustomers(token)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	err = alfa.UpdateGroups(token)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	aLessons, err := alfa.GetLessons(token, -1)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	aRegLessons, err := alfa.GetRegularLessons(token, -1)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	aLessons = append(aLessons, aRegLessons...)

	gcalService, err := gcal.GetService(context.Background(), gCalCreds)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	gLessons, err := gcal.GetLessons(gcalService, -1)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	gAdd, aAdd := logic.RemoveCommonElements(aLessons, gLessons)

	return nil
}

package sync

import (
	"CalSync/internal/alfa"
	"fmt"
)

func SyncCalendars(gCalCreds []byte, email, alfaApiKey string) error {
	token, err := alfa.GetAPIToken(email, alfaApiKey)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	// alessons, err := alfa.GetLessons(token, -1)
	// if err != nil {
	// 	return fmt.Errorf("error at SyncCalendars: %w", err)
	// }
	// fmt.Println(alessons)

	regaLessons, err := alfa.GetRegularLessons(token, -1)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}
	fmt.Println(regaLessons)

	// gcalService, err := gcal.GetService(context.Background(), gCalCreds)
	// if err != nil {
	// 	return fmt.Errorf("error at SyncCalendars: %w", err)
	// }

	// glessons, err := gcal.GetLessons(gcalService, -1)
	// if err != nil {
	// 	return fmt.Errorf("error at SyncCalendars: %w", err)
	// }

	// fmt.Println(glessons)

	return nil
}

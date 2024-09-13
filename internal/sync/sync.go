package sync

import (
	"CalSync/internal/gcal"
	"context"
	"fmt"
)

func SyncCalendars(gCalCreds []byte, email, alfaApiKey string) error {
	// token, err := alfa.GetAPIToken(email, alfaApiKey)
	// if err != nil {
	// 	return fmt.Errorf("error at SyncCalendars: %w", err)
	// }

	// lessons, err := alfa.GetLessons(token)
	// if err != nil {
	// 	return fmt.Errorf("error at SyncCalendars: %w", err)
	// }

	// regLessons, err := alfa.GetRegularLessons(token)
	// if err != nil {
	// 	return fmt.Errorf("error at SyncCalendars: %w", err)
	// }

	// fmt.Println(lessons)
	// fmt.Println(regLessons)

	gcalService, err := gcal.GetService(context.Background(), gCalCreds)
	if err != nil {
		return fmt.Errorf("error at SyncCalendars: %w", err)
	}

	gcal.GetLessons(gcalService)

	return nil
}

package main

import (
	"CalSync/internal/sync"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load("env/alfacreds.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Получение email и api_key из переменных окружения
	email := os.Getenv("EMAIL")
	alfaApiKey := os.Getenv("API_KEY")
	if email == "" || alfaApiKey == "" {
		log.Fatalf("Email or API key is missing in the env/alfacred.env file")
	}

	// Аутентификация Google API
	ctx := context.Background()
	clientSecretFile := "env/credentials.json"
	gCalService, err := calendar.NewService(ctx, option.WithCredentialsFile(clientSecretFile))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	s := gocron.NewScheduler(time.UTC)
	if err != nil {
		log.Fatalf("Unable to create new scheduler: %v", err)
	}

	// Cинхронизация каждые 30 минут
	s.Every(30).Minutes().Do(func() {
		err := sync.SyncCalendars(gCalService, email, alfaApiKey)
		if err != nil {
			log.Fatalf("Failed to sync calendars: %v", err)
		} else {
			fmt.Println("Calendars synced successfully!")
		}
	})

	s.StartBlocking()
}

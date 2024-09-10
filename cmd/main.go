package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"CalSync/internal/alfa"
	"CalSync/internal/gcal"
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

	// Аутентификация Alfa календаря
	_, err = alfa.GetAlfaCalendarToken(email, alfaApiKey)
	if err != nil {
		log.Fatalf("Unable to authenticate to alfa calendar: %v", err)
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
		err := gcal.SyncGoogleCalendar(gCalService, email, alfaApiKey)
		if err != nil {
			log.Fatalf("Failed to sync googlr calendar: %v", err)
		} else {
			err = alfa.SyncAlfaCalendar(gCalService, email, alfaApiKey)
			if err != nil {
				log.Fatalf("Failed to sync googlr calendar: %v", err)
			} else {
				fmt.Println("Calendars synced successfully!")
			}
		}
	})

	s.StartBlocking()
}

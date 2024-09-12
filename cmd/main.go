package main

import (
	"CalSync/internal/sync"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

const (
	alfaCredsFile = "env/alfacreds.env"
	gCalCredsFile = "env/credentials.json"
)

const minutesPeriod = 30

func main() {
	err := godotenv.Load(alfaCredsFile)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Получение email и api_key из переменных окружения
	email := os.Getenv("EMAIL")
	alfaApiKey := os.Getenv("API_KEY")
	if email == "" || alfaApiKey == "" {
		log.Fatalf("Email or API key is missing in the env/alfacred.env file")
	}

	// Получение google calendar credentials
	gCalCreds, err := os.ReadFile(gCalCredsFile)
	if err != nil {
		log.Fatalf("Unable to read gCal secret file: %v", err)
	}

	s := gocron.NewScheduler(time.UTC)
	if err != nil {
		log.Fatalf("Unable to create new scheduler: %v", err)
	}

	// Cинхронизация каждые 30 минут
	s.Every(minutesPeriod).Minutes().Do(func() {
		err := sync.SyncCalendars(gCalCreds, email, alfaApiKey)
		if err != nil {
			log.Fatalf("Failed to sync calendars: %v", err)
		} else {
			fmt.Println("Calendars synced successfully!")
		}
	})

	s.StartBlocking()
}

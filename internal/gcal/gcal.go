package gcal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const tokenFile = "env/token.json"

var calendarIDs = []string{
	"c3aa27d84ad2921ebe2e97f163d69cf3b930292822267de4a2c808661cda8fcf@group.calendar.google.com", // 01 - Pro
	"8b218d7a47971ddddd851a867cfb558cfaf555d189e7555c1229aea777c9a03f@group.calendar.google.com", // 02 - G Labs
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // 01 - V Labs
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

func getClient(creds []byte) (*http.Client, error) {
	config, err := google.ConfigFromJSON(creds, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	tokFile := tokenFile
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, fmt.Errorf("unable to get token from web: %w", err)
		}
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok), nil
}

func GetService(ctx context.Context, creds []byte) (*calendar.Service, error) {
	client, err := getClient(creds)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %w", err)
	}

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Calendar client: %w", err)
	}

	return srv, nil
}

// Получение занятий из календаря по его id (если calId == -1 -> все занятия)
func GetLessons(srv *calendar.Service, calId int) ([]*calendar.Event, error) {
	var lessons []*calendar.Event
	t := time.Now().Format(time.RFC3339)

	if calId == -1 {
		for _, id := range calendarIDs {
			events, err := srv.Events.List(id).ShowDeleted(false).
				SingleEvents(true).TimeMin(t).MaxResults(100).OrderBy("startTime").Do()
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve events from calendar: %w", err)
			}

			lessons = append(lessons, events.Items...)
		}
	} else {
		events, err := srv.Events.List(calendarIDs[calId]).ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(100).OrderBy("startTime").Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve events from calendar: %w", err)
		}

		lessons = events.Items
	}
	return lessons, nil
}

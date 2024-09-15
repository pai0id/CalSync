package gcal

import (
	"CalSync/internal/logic"
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
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // 03 - V Labs
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // 04 - V+Pro
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // 06 - VP+Pro
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // 07 - VGP
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // 08 - Vocal + G
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // Зал loft
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // ОПЛОТ - игровая
	"d777f17fae186c081e02122ac5142100ffb861293a21b5b9606a370bf05439b0@group.calendar.google.com", // ОПЛОТ - мастеровая
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

func convertGoogleEventToLesson(event *calendar.Event, calId int) (*logic.Lesson, error) {
	startTime, err := time.Parse(time.RFC3339, event.Start.DateTime)
	if err != nil {
		return nil, fmt.Errorf("error parsing start time: %v", err)
	}

	endTime, err := time.Parse(time.RFC3339, event.End.DateTime)
	if err != nil {
		return nil, fmt.Errorf("error parsing end time: %v", err)
	}

	lesson := &logic.Lesson{
		Name:     event.Summary,
		Date:     startTime,
		CalId:    calId,
		TimeFrom: startTime,
		TimeTo:   endTime,
	}

	return lesson, nil
}

// Получение занятий из календаря по его id (если calId == -1 -> все занятия)
func GetLessons(srv *calendar.Service, calId int) ([]logic.Lesson, error) {
	var lessons []logic.Lesson

	t := time.Now().Format(time.RFC3339)

	if calId == -1 {
		for i, id := range calendarIDs {
			events, err := srv.Events.List(id).ShowDeleted(false).
				SingleEvents(true).TimeMin(t).MaxResults(100).OrderBy("startTime").Do()
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve events from calendar: %w", err)
			}

			for _, e := range events.Items {
				lesson, err := convertGoogleEventToLesson(e, i)
				if err != nil {
					return nil, fmt.Errorf("error converting event to lesson: %w", err)
				}
				lessons = append(lessons, *lesson)
			}
		}
	} else {
		events, err := srv.Events.List(calendarIDs[calId]).ShowDeleted(false).
			SingleEvents(true).TimeMin(t).OrderBy("startTime").Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve events from calendar: %w", err)
		}

		for _, e := range events.Items {
			lesson, err := convertGoogleEventToLesson(e, calId)
			if err != nil {
				return nil, fmt.Errorf("error converting event to lesson: %w", err)
			}
			lessons = append(lessons, *lesson)
		}
	}

	return lessons, nil
}

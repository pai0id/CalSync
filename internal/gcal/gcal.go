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
	"50c82426236347be525723d287c83a20de8b336689c7bcbdc62c2197249f2d7b@group.calendar.google.com", // 01 - Pro
	"e094039b8c6a95961d567523fdf8586f984c5d445b4131d9c1be62511cbbbe84@group.calendar.google.com", // 02 - G Labs
	"5bce8f738c60fc06c5b4eba02c21b38b55f00146c069f8f1447b0dd27e562a8f@group.calendar.google.com", // 03 - V Labs
	"9e1ad33080757c217c4648f6d8c28b23f652a2c8419533fe77d93465088b91de@group.calendar.google.com", // 04 - V+Pro
	"3ddce7e25bee8475a47b21fa7f19e3d6188b8501ee1b0eafe90ce561737e6ef2@group.calendar.google.com", // 06 - VP+Pro
	"3da901d09b590abbbbdb77651ae2e8ff1ace46c7e51d10a32d9509f66c55ac33@group.calendar.google.com", // 07 - VGP
	"5cfd1da735d9c354ac4c45024056289e29c580d6effe60a2139819c24b14bbc9@group.calendar.google.com", // 08 - Vocal + G
	"3e9dca44aec311959f94fadf8bb356c2dec2bdf6c82c44ac8af545357d86a2f8@group.calendar.google.com", // Зал loft
	"896e972cd329cc65d4fb0320a1fba1fe3a0e71804e7c64563da54d4558819787@group.calendar.google.com", // ОПЛОТ - игровая
	"2d62001c42ee5c946a2febf2438f51e3ec3a3530319b23f681ded5797db63c47@group.calendar.google.com", // ОПЛОТ - мастеровая
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
	config, err := google.ConfigFromJSON(creds, calendar.CalendarScope)
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
		Date:     startTime.Truncate(24 * time.Hour).Add(time.Hour * -3),
		CalId:    calId,
		TimeFrom: startTime.Truncate(time.Minute),
		TimeTo:   endTime.Truncate(time.Minute),
	}

	return lesson, nil
}

// Получение занятий из календаря по его id (если calId == -1 -> все занятия)
func GetLessons(srv *calendar.Service, calId int) ([]logic.Lesson, error) {
	var lessons []logic.Lesson

	t := time.Now().Truncate(24 * time.Hour)

	if calId == -1 {
		for i, id := range calendarIDs {
			events, err := srv.Events.List(id).ShowDeleted(false).
				SingleEvents(true).TimeMin(t.Format(time.RFC3339)).TimeMax(t.AddDate(0, 0, 8).Format(time.RFC3339)).OrderBy("startTime").Do()
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
			SingleEvents(true).TimeMin(t.Format(time.RFC3339)).TimeMax(t.AddDate(0, 0, 8).Format(time.RFC3339)).OrderBy("startTime").Do()
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

func AddEvents(ctx context.Context, service *calendar.Service, lessons []logic.Lesson) error {
	for _, lesson := range lessons {
		calendarID := calendarIDs[lesson.CalId]

		event := &calendar.Event{
			Summary:     lesson.Name,
			Start:       &calendar.EventDateTime{DateTime: lesson.TimeFrom.Format(time.RFC3339), TimeZone: "UTC"},
			End:         &calendar.EventDateTime{DateTime: lesson.TimeTo.Format(time.RFC3339), TimeZone: "UTC"},
			Description: fmt.Sprintf("Event scheduled for %s", lesson.Date.Format("2006-01-02")),
		}

		_, err := service.Events.Insert(calendarID, event).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("could not insert event %s into calendar %s: %v", lesson.Name, calendarID, err)
		}

		fmt.Printf("Added event '%s' to calendar '%s'\n", lesson.Name, calendarID)
	}

	return nil
}

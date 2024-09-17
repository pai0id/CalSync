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
	"aa7dfbd43b5f0ffa3a63a57faae946d8bd317f3c45bfe221919ce5f7c90852fc@group.calendar.google.com", // 01 - Pro
	"567bfc5ca2de6bc54583aa4744cf049a65c451d2bd897d9dc003aa9dfa287a87@group.calendar.google.com", // 02 - G Labs
	"59bfd0df6be1efe40b5f10c45d663cd1d14b924ebf3b709cfcd69fb09897138f@group.calendar.google.com", // 03 - V Labs
	"f1e41b8ff60c54795946e150920b2971cdf929a015fee172824fb0e961accd42@group.calendar.google.com", // 04 - V+Pro
	"efee209e8160bf8bb891ab36abc86d7e3a8ea82dd2fd6a7d5cf81ac1a61fb42d@group.calendar.google.com", // 06 - VP+Pro
	"e9c080d310635936318c85aa7c308baabbfb251eabd05342e52f47e20eae4816@group.calendar.google.com", // 07 - VGP
	"73982ef0ab637a3fc9d5645be9018b330669937add99b6825c7a09725bbe6df8@group.calendar.google.com", // 08 - Vocal + G
	"6d8844a5a5e8b4bb4fcc83c7e0571640e11ff7c2e2f1989f79bee62fda5f57cb@group.calendar.google.com", // Зал loft
	"6bf6d9cce226a2072bbbf0c16a0c0f96af99c424d167f353d9ae0b57eb6da23f@group.calendar.google.com", // ОПЛОТ - игровая
	"9ca57980980ce6b06a929b6eedbc624b8a1d3227733fb86946d65242e5752ae4@group.calendar.google.com", // ОПЛОТ - мастеровая
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
				SingleEvents(true).TimeMin(t.Format(time.RFC3339)).TimeMax(t.AddDate(0, 1, 1).Format(time.RFC3339)).OrderBy("startTime").Do()
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
			SingleEvents(true).TimeMin(t.Format(time.RFC3339)).TimeMax(t.AddDate(0, 1, 1).Format(time.RFC3339)).OrderBy("startTime").Do()
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

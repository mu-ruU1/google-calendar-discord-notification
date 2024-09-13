package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Calendar struct {
	Summary     string
	Description string
	Start       string
	End         string
}

var CalendarEvents []Calendar

func loadEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s not set", key)
	}

	return value, nil
}

func formatKey(key string) string {
	return strings.ToLower(key[2:])
}

func cal() {
	credentialKey := [...]string{"G_TYPE", "G_PROJECT_ID", "G_PRIVATE_KEY_ID", "G_PRIVATE_KEY", "G_CLIENT_EMAIL", "G_CLIENT_ID", "G_AUTH_URL", "G_TOKEN_URL", "G_AUTH_PROVIDER_X509_CERT_URL", "G_CLIENT_X509_CERT_URL", "G_UNIVERSE_DOMAIN"}

	credentials := make(map[string]string)

	for _, key := range credentialKey {
		if value, err := loadEnv(key); err != nil {
			fmt.Println("Error loading environment variable:", err)
			os.Exit(1)
		} else {
			credentials[formatKey(key)] = value
		}
	}

	// Convert map to JSON
	credentialJSON, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println("Error converting credentials to JSON:", err)
		os.Exit(1)
	}

	// カレンダーIDを取得
	calId, err := loadEnv("G_CALENDAR_ID")
	if err != nil {
		calId = "primary"
	}

	ctx := context.Background()

	// GoogleカレンダーAPIサービスを作成
	srv, err := calendar.NewService(ctx, option.WithCredentialsJSON(credentialJSON))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// 日時を取得
	now := time.Now().Format(time.RFC3339)
	tommorow := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)

	// カレンダーの予定を取得
	events, err := srv.Events.List(calId).ShowDeleted(false).SingleEvents(true).TimeMin(now).TimeMax(tommorow).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}

	// 予定を抽出
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		fmt.Println("Upcoming events:")
		for _, item := range events.Items {
			var start string
			var end string

			if item.Start.DateTime != "" {
				start = item.Start.DateTime
				end = item.Start.DateTime
			} else {
				start = item.Start.Date
				end = item.End.Date
			}

			calendarEvents := Calendar{
				Summary:     item.Summary,
				Description: item.Description,
				Start:       start,
				End:         end,
			}

			CalendarEvents = append(CalendarEvents, calendarEvents)
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

// func formatKey(key string) string {
// 	return strings.ToLower(key[2:])
// }

func formatTime(timeStr string) string {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Println("Error converting time:", err)
		os.Exit(1)
	}

	ft := t.Format("01/02 15:04")
	return ft
}

func cal() {
	// Convert map to JSON
	credentialStr, _ := loadEnv("G_CREDENTIALS")
	credentialByte := []byte(credentialStr)

	// カレンダーIDを取得
	calId, err := loadEnv("G_CALENDAR_ID")
	if err != nil {
		calId = "primary"
	}

	ctx := context.Background()

	// GoogleカレンダーAPIサービスを作成
	srv, err := calendar.NewService(ctx, option.WithCredentialsJSON(credentialByte))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// 日時を取得
	nowRfc := time.Now().Format(time.RFC3339)
	tommorowRfc := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)

	// カレンダーの予定を取得
	events, err := srv.Events.List(calId).ShowDeleted(false).SingleEvents(true).TimeMin(nowRfc).TimeMax(tommorowRfc).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}

	// 予定を抽出
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			var start string
			var end string

			if item.Start.DateTime != "" {
				start = formatTime(item.Start.DateTime)
				end = formatTime(item.End.DateTime)
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

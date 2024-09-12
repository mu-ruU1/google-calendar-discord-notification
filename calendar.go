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

func loadEnv(key string) (value string) {
	value, ok := os.LookupEnv(key)

	if !ok {
		os.Exit(1)
	}

	return
}

func cal() {
	// 認証済みのクライアントを作成
	ctx := context.Background()

	// GoogleカレンダーAPIサービスを作成
	srv, err := calendar.NewService(ctx, option.WithCredentialsFile("./credential.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// 現在の日時を取得
	t := time.Now().Format(time.RFC3339)

	// カレンダーの予定を取得
	events, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}

	// 予定があるか確認し、出力
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		fmt.Println("Upcoming events:")
		for _, item := range events.Items {
			var start string
			if item.Start.DateTime != "" {
				start = item.Start.DateTime
			} else {
				start = item.Start.Date
			}
			fmt.Printf("%s (%s)\n", item.Summary, start)
		}
	}
}

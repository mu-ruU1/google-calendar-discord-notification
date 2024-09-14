package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	roleId, _        = loadEnv("D_ROLE_ID")
	channelL2L3, _   = loadEnv("D_A1_CHANNEL_ID")
	channelServer, _ = loadEnv("D_A2_CHANNEL_ID")
)

func msgCreate(c *Calendar) (embed *discordgo.MessageEmbed, channelId string, err error) {
	var (
		color int
		team  string
	)

	// チームIDとイベント名に分割
	summary := strings.SplitN(c.Summary, "_", 3)
	if len(summary) < 2 {
		err = fmt.Errorf("invalid summary: %s", c.Summary)
		return
	}

	teamId := summary[0] + "_" + summary[1]
	event := summary[2]

	switch teamId {
	case "学生_L2L3":
		channelId = channelL2L3
		color = 0x00ff00
		team = "L2L3チーム"
	case "学生_server":
		channelId = channelServer
		color = 0xff0000
		team = "サーバチーム"
	}

	embed = NewEmbed().
		SetTitle(":calendar:"+event+" ( "+team+" ) ").
		SetDescription(c.Description).
		AddField("開始時刻", c.Start).
		AddField("終了時刻", c.End).
		SetColor(color).MessageEmbed

	return
}

func discord() {
	// var embed *discordgo.MessageEmbed

	token, err := loadEnv("D_BOT_TOKEN")
	if err != nil {
		fmt.Println("Error loading environment variable:", err)
		os.Exit(1)
	}

	token = "Bot " + token

	discord, err := discordgo.New(token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
	}

	for i := range CalendarEvents {
		embed, channelId, err := msgCreate(&CalendarEvents[i])
		if err != nil {
			fmt.Println("Error creating message:", err)
			continue
		}
		discord.ChannelMessageSend(channelId, "<@"+roleId+"> 本日 予定があります")
		discord.ChannelMessageSendEmbed(channelId, embed)
	}

	err = discord.Close()

	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}

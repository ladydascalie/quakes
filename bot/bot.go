package bot

import (
	"fmt"
	"os"

	"github.com/tucnak/telebot"
	"github.com/ladydascalie/earthquakes/domain"
)

var (
	botToken string
	// Bot represents the telegram bot
	Bot *telebot.Bot

	chat = telebot.Chat{
		Type:     telebot.ChatChannel,
		Username: "usgs_quakes",
	}
)

// Begin the connection to the telegram API
func Begin() {
	botToken = os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic("canont load bot token")
	}
	var err error
	Bot, err = telebot.NewBot(botToken)
	if err != nil {
		panic(err)
	}
}

// NotifyTelegramChannel produces and sends 2 formatted messages
// to the telegram channel
func NotifyTelegramChannel(alert *domain.Alert) {
	msg := fmt.Sprintf(`
	*%.1f - %s*
	- [View online](https://quakes.cable.fyi/en/%s)
	- [Official USGS report](%s)
	`, alert.Magnitude, alert.Place, alert.ID, alert.URL)

	Bot.SendLocation(chat, &telebot.Location{
		Longitude: float32(alert.Lng),
		Latitude:  float32(alert.Lat),
	}, nil)

	Bot.SendMessage(chat, msg, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

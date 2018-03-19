package bot

import (
	"fmt"

	"github.com/ladydascalie/earthquakes/config"
	"github.com/ladydascalie/earthquakes/domain"
	"github.com/tucnak/telebot"
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
	var err error
	Bot, err = telebot.NewBot(config.BotToken)
	if err != nil {
		panic(err)
	}
}

// NotifyTelegramChannel produces and sends 2 formatted messages
// to the telegram channel
func NotifyTelegramChannel(alert *domain.Alert) {
	msg := fmt.Sprintf(`
	*%.1f - %s*
	- [View online](%s/en/%s)
	- [Official USGS report](%s)
	`, alert.Magnitude, alert.Place, config.BaseURL, alert.ID, alert.URL)

	Bot.SendLocation(chat, &telebot.Location{
		Longitude: float32(alert.Lng),
		Latitude:  float32(alert.Lat),
	}, nil)

	Bot.SendMessage(chat, msg, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

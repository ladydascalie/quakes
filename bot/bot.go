package bot

import (
	"fmt"

	"github.com/ladydascalie/earthquakes/config"
	"github.com/ladydascalie/earthquakes/domain"
	"github.com/tucnak/telebot"
)

var (
	// bot represents the telegram bot
	bot *telebot.Bot

	// chat represents the telegram channel
	// you want to broadcast to
	chat = telebot.Chat{
		Type:     telebot.ChatChannel,
		Username: config.QuakesChannel,
	}
)

// Begin the connection to the telegram API
func Begin() {
	var err error
	bot, err = telebot.NewBot(config.BotToken)
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

	bot.SendLocation(chat, &telebot.Location{
		Longitude: float32(alert.Lng),
		Latitude:  float32(alert.Lat),
	}, nil)

	bot.SendMessage(chat, msg, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

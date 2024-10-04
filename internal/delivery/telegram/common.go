package telegram

import (
	"time"

	"github.com/AsaHero/just-ask-bot/internal/inerr"
	telegram_bot "github.com/AsaHero/just-ask-bot/pkg/telegram-bot"
	"gopkg.in/telebot.v3"
)

func VoiceRecordingAction(user telebot.Recipient, bot *telegram_bot.TelegramBot, done chan struct{}, failed chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	err := bot.Notify(user, telebot.RecordingAudio)
	if err != nil {
		inerr.Err(err, "failed to send action", "global", "VoiceRecordingAction", "bot.Notify")
		return
	}

	for {
		select {
		case <-ticker.C:
			err := bot.Notify(user, telebot.RecordingAudio)
			if err != nil {
				inerr.Err(err, "failed to send action", "global", "VoiceRecordingAction", "bot.Notify")
				return
			}

		case <-done:
			return
		case <-failed:
			return
		}
	}

}

package telegram

import (
	"context"

	"github.com/AsaHero/just-ask-bot/internal/delivery/telegram/models"
	"github.com/AsaHero/just-ask-bot/internal/inerr"
	"gopkg.in/telebot.v3"
)

func (h *Router) TextMessage(c telebot.Context) error {
	user := c.Sender()
	if user == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	message := c.Message()
	if message == nil {
		return nil
	}

	question := message.Text
	if question == "" {
		return nil
	}

	figure, ok := h.personas[user.ID]
	if !ok {
		return c.Send(models.MessageFiguresFirstChoose)
	}

	// Channels to manage the VoiceRecordingAction
	done := make(chan struct{})
	failed := make(chan struct{})
	defer close(failed)

	go VoiceRecordingAction(user, h.TelegramBot, done, failed)

	// Handle the message
	err := h.conversationService.HandleMessage(ctx, user.ID, figure, question)
	if err != nil {
		failed <- struct{}{}
		return inerr.Err(err, "failed to handle message", "Router", "TextMessage", "conversationService.HandleMessage")
	}

	done <- struct{}{}
	return nil
}

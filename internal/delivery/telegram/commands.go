package telegram

import (
	"context"

	"github.com/AsaHero/just-ask-bot/internal/delivery/telegram/models"
	"github.com/AsaHero/just-ask-bot/internal/entity"
	"github.com/AsaHero/just-ask-bot/internal/inerr"
	"gopkg.in/telebot.v3"
)

func (h *Router) StartCommand(c telebot.Context) error {
	user := c.Sender()
	if user == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := h.userService.Upsert(ctx, &entity.Users{
		ExternalID: user.ID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	})
	if err != nil {
		return inerr.Err(err, "failed to upsert user", "Router", "StartCommand", "userService.Upsert")
	}

	return c.Send(models.MessageHomePage, models.KeyboardHomePageChoose)
}

func (h *Router) ChooseCommand(c telebot.Context) error {
	user := c.Sender()
	if user == nil {
		return nil
	}

	keyboard := getPersonaKeyboard(h.figures)

	return c.Send(models.MessageFiguresChoose, keyboard)
}

package telegram

import (
	"fmt"

	"github.com/AsaHero/just-ask-bot/internal/delivery/telegram/models"
	"github.com/AsaHero/just-ask-bot/internal/entity"
	"gopkg.in/telebot.v3"
)

func (h *Router) FiguresChoose(c telebot.Context) error {
	user := c.Sender()
	if user == nil {
		return nil
	}

	data := c.Callback()
	if data == nil {
		return nil
	}

	figureID := data.Data
	figure, ok := h.figures[figureID]
	if !ok {
		return c.Send(models.MessageFigureNotAvailable)
	}

	h.personas[user.ID] = figure

	return c.Send(fmt.Sprintf(models.MessageFigureAvailable, figure.Name), telebot.ModeHTML)
}

// getPersonaKeyboard возвращает клавиатуру с историческими личностями
func getPersonaKeyboard(figures map[string]*entity.Figures) *telebot.ReplyMarkup {
	var keyboard [][]telebot.InlineButton
	for k, p := range figures {
		btn := telebot.InlineButton{
			Unique: models.ButtonFiguresChoose.Unique,
			Text:   p.Name,
			Data:   k,
		}
		keyboard = append(keyboard, []telebot.InlineButton{btn})
	}

	return &telebot.ReplyMarkup{InlineKeyboard: keyboard}
}

package models

import "gopkg.in/telebot.v3"

// Messages
var (
	MessageHomePage = `Добро пожаловать! Используйте /choose для выбора исторической личности.`
)

// Buttons
var (
	ButtonHomePage       = &telebot.Btn{Unique: "hp_Home", Text: "Меню"}
	ButtonHomePageChoose = &telebot.Btn{Unique: "hp_Choose", Text: "Выбрать личность"}
)

// KeyBoard
var KeyboardHomePageChoose = &telebot.ReplyMarkup{
	InlineKeyboard: [][]telebot.InlineButton{
		{
			*ButtonHomePageChoose.Inline(),
		},
	},
}

var KeyboardHomePage = &telebot.ReplyMarkup{
	InlineKeyboard: [][]telebot.InlineButton{
		{
			*ButtonHomePage.Inline(),
		},
	},
}

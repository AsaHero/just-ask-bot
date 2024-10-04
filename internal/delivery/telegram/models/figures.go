package models

import "gopkg.in/telebot.v3"

// Messages
var (
	MessageFiguresFirstChoose = `Сначала выберите историческую личность с помощью /choose`
	MessageFiguresChoose      = `Выберите историческую личность:`
	MessageFigureNotAvailable = `Выбранная личность временно не доступна. Пожалуйста выберите другую или попробуйте позже.`
	MessageFigureAvailable    = `Вы теперь общаетесь с <b>%s</b>.

<b>Задайте свой вопрос:</b>`
)

// Buttons
var (
	ButtonFiguresChoose = &telebot.Btn{Unique: "f_Choose"}
)

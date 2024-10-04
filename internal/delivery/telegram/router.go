package telegram

import (
	"encoding/json"
	"log"
	"os"

	"github.com/AsaHero/just-ask-bot/internal/delivery/telegram/models"
	"github.com/AsaHero/just-ask-bot/internal/entity"
	"github.com/AsaHero/just-ask-bot/internal/usecase/conversation"
	"github.com/AsaHero/just-ask-bot/internal/usecase/users"
	"github.com/AsaHero/just-ask-bot/pkg/config"
	telegram_bot "github.com/AsaHero/just-ask-bot/pkg/telegram-bot"
	"gopkg.in/telebot.v3"
	telebot_middleware "gopkg.in/telebot.v3/middleware"
)

type Router struct {
	Config              *config.Config
	TelegramBot         *telegram_bot.TelegramBot
	figures             map[string]*entity.Figures
	personas            map[int64]*entity.Figures
	userService         users.Users
	conversationService conversation.Conversation
}

func NewRouter(config *config.Config, telegramBot *telegram_bot.TelegramBot, userService users.Users, conversationService conversation.Conversation) {
	var figures []*entity.Figures

	file, err := os.Open("./static/figures.json")
	if err != nil {
		log.Fatalf("failed to open figures.json: %v", err)
	}

	if err := json.NewDecoder(file).Decode(&figures); err != nil {
		log.Fatalf("failed to parse figures.json: %v", err)
	}

	var figuresMap = make(map[string]*entity.Figures, len(figures))
	for _, v := range figures {
		figuresMap[v.GUID] = v
	}

	telegramBot.Use(telebot_middleware.Recover())
	telegramBot.Use(telebot_middleware.AutoRespond())
	telegramBot.Use(telebot_middleware.Logger(log.Default()))

	router := &Router{
		Config:              config,
		TelegramBot:         telegramBot,
		figures:             figuresMap,
		userService:         userService,
		personas:            make(map[int64]*entity.Figures),
		conversationService: conversationService,
	}

	telegramBot.Handle("/start", router.StartCommand)
	telegramBot.Handle("/choose", router.ChooseCommand)
	telegramBot.Handle(models.ButtonHomePage, router.StartCommand)
	telegramBot.Handle(models.ButtonHomePageChoose, router.ChooseCommand)
	telegramBot.Handle(models.ButtonFiguresChoose, router.FiguresChoose)
	telegramBot.Handle(telebot.OnText, router.TextMessage)
}

package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AsaHero/just-ask-bot/api"
	"github.com/AsaHero/just-ask-bot/internal/delivery/telegram"
	"github.com/AsaHero/just-ask-bot/internal/infrastructure/openai"
	users_repo "github.com/AsaHero/just-ask-bot/internal/infrastructure/repository/users"
	"github.com/AsaHero/just-ask-bot/internal/usecase/conversation"
	"github.com/AsaHero/just-ask-bot/internal/usecase/users"

	"github.com/AsaHero/just-ask-bot/pkg/config"
	"github.com/AsaHero/just-ask-bot/pkg/database/postgres"
	"github.com/AsaHero/just-ask-bot/pkg/logger"
	telegram_bot "github.com/AsaHero/just-ask-bot/pkg/telegram-bot"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	server *http.Server
	bot    *telegram_bot.TelegramBot
	logger *logrus.Logger
	db     *gorm.DB
}

func NewApp(cfg *config.Config) *App {

	logger := logger.Init(cfg, cfg.APP+".log")

	telegramBot := telegram_bot.NewTelegramBot(cfg.Telegram.BotToken, cfg.Telegram.WebhookURL)

	db, err := postgres.NewGORMPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	return &App{
		Config: cfg,
		logger: logger,
		bot:    telegramBot,
		db:     db,
	}
}

func (a *App) Run() error {
	contextDeadline, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	openaiAPI, err := openai.New(a.Config)
	if err != nil {
		return err
	}

	// repos init
	userRepo := users_repo.NewUsersRepository(a.db)

	// services init
	usersService := users.NewUsersService(contextDeadline, userRepo)
	conversationService := conversation.NewConversationService(contextDeadline, a.bot, usersService, openaiAPI)

	// bot router init
	telegram.NewRouter(a.Config, a.bot, usersService, conversationService)

	// api router init
	handler := api.NewRouter(a.Config, a.bot)

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	// start telegram bot listen updates
	go a.bot.Start()

	return a.server.ListenAndServe()
}

func (a *App) Stop() {
	a.server.Close()

	a.bot.Stop()

	sqlDB, _ := a.db.DB()

	sqlDB.Close()

	a.logger.Writer().Close()
}

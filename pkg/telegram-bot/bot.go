package telegram_bot

import (
	"log"
	"net/http"

	"gopkg.in/telebot.v3"
)

type BotState string
type BotSessionType string

type TelegramBot struct {
	*telebot.Bot
	webhook         *telebot.Webhook
	sessionHandlers map[BotSessionType]map[BotState]func(telebot.Context) error
}

func NewTelegramBot(token string, webhookURL string) *TelegramBot {
	webhook := &telebot.Webhook{
		AllowedUpdates: []string{
			"callback_query",
			"edited_message",
			"message",
			"pre_checkout_query",
		},
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: webhookURL,
		},
		DropUpdates: true,
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:       token,
		Poller:      webhook,
		Synchronous: false,
		Verbose:     false,
		ParseMode:   telebot.ModeMarkdown,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &TelegramBot{
		Bot:             bot,
		webhook:         webhook,
		sessionHandlers: make(map[BotSessionType]map[BotState]func(telebot.Context) error),
	}
}

func (b TelegramBot) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	b.webhook.ServeHTTP(w, r)
}

// Adds Handler to handle specific session type and state
func (b *TelegramBot) Subscribe(sessionType BotSessionType, sessionState BotState, handler func(telebot.Context) error) {
	if _, ok := b.sessionHandlers[sessionType]; !ok {
		b.sessionHandlers[sessionType] = make(map[BotState]func(telebot.Context) error)
	}

	if _, ok := b.sessionHandlers[sessionType][sessionState]; !ok {
		b.sessionHandlers[sessionType][sessionState] = handler
	} else {
		log.Fatalf("type: %s. state: %s handler already subscribed", sessionType, sessionState)
	}
}

// Handle specific session type and state
func (b *TelegramBot) Publish(sessionType BotSessionType, sessionState BotState, c telebot.Context) {
	if handlers, found := b.sessionHandlers[sessionType]; found {
		if handler, found := handlers[sessionState]; found {
			go handler(c)
		}
	}
}

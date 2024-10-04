package conversation

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/AsaHero/just-ask-bot/internal/entity"
	"github.com/AsaHero/just-ask-bot/internal/inerr"
	openai_api "github.com/AsaHero/just-ask-bot/internal/infrastructure/openai"
	"github.com/AsaHero/just-ask-bot/internal/usecase/users"
	telegram_bot "github.com/AsaHero/just-ask-bot/pkg/telegram-bot"
	"github.com/AsaHero/just-ask-bot/pkg/tts"
	"github.com/AsaHero/just-ask-bot/pkg/utility"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/telebot.v3"
)

type conversationService struct {
	contextTimeout time.Duration
	bot            *telegram_bot.TelegramBot
	userService    users.Users
	openaiAPI      openai_api.OpenaiAPI
}

func NewConversationService(contextTimeout time.Duration, bot *telegram_bot.TelegramBot, userService users.Users, openaiAPI openai_api.OpenaiAPI) Conversation {
	return &conversationService{
		contextTimeout: contextTimeout,
		userService:    userService,
		openaiAPI:      openaiAPI,
		bot:            bot,
	}
}

func (s *conversationService) HandleMessage(ctx context.Context, userID int64, figure *entity.Figures, question string) error {
	message := formatMessagePrompt(figure, question)

	fmt.Println(entity.SystemPrompt)
	fmt.Println(message)

	answer, err := s.openaiAPI.ChatCompletion(ctx, entity.SystemPrompt, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: question,
		},
	})
	if err != nil {
		return inerr.Err(err, "failed to get openai chat completion", "conversationService", "HandleMessage", "openaiAPI.ChatCompletion")
	}

	fmt.Println(answer)

	audioData, err := tts.TextToSpeech(answer, figure.Name)
	if err != nil {
		return inerr.Err(err, "failed to get speech to text", "conversationService", "HandleMessage", "tts.TextToSpeech")
	}

	audioReader := bytes.NewReader(audioData)
	voiceMsg := &telebot.Voice{File: telebot.FromReader(audioReader)}

	_, err = s.bot.Send(&telebot.User{ID: userID}, voiceMsg)
	if err != nil {
		return inerr.Err(err, "error sending voice message", "conversationService", "HandleMessage", "bot.Send")
	}

	return nil
}

func formatMessagePrompt(figure *entity.Figures, question string) string {
	var messagePrompt string

	traits := utility.DashFormat(figure.KeyTraits)
	achievements := utility.DashFormat(figure.NotableAchievements)
	expertise := utility.DashFormat(figure.AreasOfExpertise)

	messagePrompt = fmt.Sprintf(entity.MessagePrompt,
		figure.Name,
		figure.Era,
		figure.BriefDescription,
		figure.BirthYear,
		figure.DeathYear,
		traits,
		achievements,
		expertise,
		figure.DeathYear,
		question,
		figure.Name,
	)

	return messagePrompt
}

package conversation

import (
	"context"

	"github.com/AsaHero/just-ask-bot/internal/entity"
)

type Conversation interface {
	HandleMessage(ctx context.Context, userID int64, figure *entity.Figures, question string) error
}

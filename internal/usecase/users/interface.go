package users

import (
	"context"

	"github.com/AsaHero/just-ask-bot/internal/entity"
)

type Users interface {
	Upsert(ctx context.Context, user *entity.Users) error
}
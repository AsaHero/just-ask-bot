package users

import (
	"github.com/AsaHero/just-ask-bot/internal/entity"
	"github.com/AsaHero/just-ask-bot/internal/infrastructure/repository"
)

type Repository interface {
	repository.BaseRepository[entity.Users]
}

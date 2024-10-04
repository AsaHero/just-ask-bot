package users

import (
	"github.com/AsaHero/just-ask-bot/internal/entity"
	"github.com/AsaHero/just-ask-bot/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type usersRepository struct {
	repository.BaseRepository[entity.Users]
	DB *gorm.DB
}

func NewUsersRepository(db *gorm.DB) Repository {
	return &usersRepository{
		BaseRepository: repository.NewBaseRepository[entity.Users](db),
		DB:             db,
	}
}

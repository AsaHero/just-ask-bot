package users

import (
	"context"
	"time"

	"github.com/AsaHero/just-ask-bot/internal/entity"
	"github.com/AsaHero/just-ask-bot/internal/inerr"
	"github.com/AsaHero/just-ask-bot/internal/infrastructure/repository/users"
	"github.com/google/uuid"
)

type usersService struct {
	contextTimeout time.Duration
	repo           users.Repository
}

func NewUsersService(contextTimeout time.Duration, repo users.Repository) Users {
	return &usersService{
		contextTimeout: contextTimeout,
		repo:           repo,
	}
}

func (s *usersService) GetByExternalID(ctx context.Context, id int64) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repo.FindOne(
		ctx,
		map[string]any{
			"external_id": id,
		},
	)
	if err != nil {
		return nil, inerr.Err(err, "failed to get user", "usersService", "GetByExternalID", "repo.FindOne")
	}

	return user, nil
}

func (s *usersService) Upsert(ctx context.Context, user *entity.Users) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	exUser, err := s.GetByExternalID(ctx, user.ExternalID)
	if err != nil && !inerr.IsErrNotFound(err) {
		return inerr.Err(err, "failed to get user", "walletsService", "Upsert", "GetByUserID")
	}

	if exUser == nil {
		user.GUID = uuid.New().String()
		err := s.repo.Create(ctx, user)
		if err != nil {
			return inerr.Err(err, "failed to create user", "walletsService", "Upsert", "walletsRepo.Create")
		}
	}

	return nil
}

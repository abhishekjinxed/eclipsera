package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	//CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, filter bson.M) (*User, error)
	UpsertUser(ctx context.Context, user *User) error
	AddFriend(ctx context.Context, userID, friendID string) error
	AcceptFriend(ctx context.Context, userID, friendID string) error
}

type service struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &service{repo: repo}
}

// func (s *service) CreateUser(ctx context.Context, user *User) error {
// 	return s.repo.CreateUser(ctx, user)
// }

func (s *service) GetUser(ctx context.Context, filter bson.M) (*User, error) {
	return s.repo.GetUser(ctx, filter)
}

func (s *service) UpsertUser(ctx context.Context, user *User) error {
	return s.repo.UpsertUser(ctx, user)
}

func (s *service) AddFriend(ctx context.Context, userID, friendID string) error {
	return s.repo.AddFriend(ctx, userID, friendID)
}
func (s *service) AcceptFriend(ctx context.Context, userID, friendID string) error {
	return s.repo.AcceptFriend(ctx, userID, friendID)
}

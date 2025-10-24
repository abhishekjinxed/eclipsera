package chat

import "context"

type Service interface {
	SaveMessage(ctx context.Context, msg *Message) error
	GetMessages(ctx context.Context, user1, user2 string) ([]Message, error)
	MarkRead(ctx context.Context, receiverID, senderID string) error
	GetUnreadMessages(ctx context.Context, receiverID string) ([]Message, error)
}

type service struct {
	repo Repository
}

func NewChatService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SaveMessage(ctx context.Context, msg *Message) error {
	return s.repo.SaveMessage(ctx, msg)
}

func (s *service) GetMessages(ctx context.Context, user1, user2 string) ([]Message, error) {
	return s.repo.GetMessages(ctx, user1, user2)
}

func (s *service) MarkRead(ctx context.Context, receiverID, senderID string) error {
	return s.repo.MarkRead(ctx, receiverID, senderID)
}

func (s *service) GetUnreadMessages(ctx context.Context, receiverID string) ([]Message, error) {
	return s.repo.GetUnreadMessages(ctx, receiverID)
}

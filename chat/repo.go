package chat

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	SaveMessage(ctx context.Context, msg *Message) error
	GetMessages(ctx context.Context, user1, user2 string) ([]Message, error)
	MarkRead(ctx context.Context, receiverID, senderID string) error
	GetUnreadMessages(ctx context.Context, receiverID string) ([]Message, error)
}

type repository struct {
	col *mongo.Collection
}

func NewChatRepository(client *mongo.Client) Repository {
	return &repository{
		col: client.Database("lumora").Collection("messages"),
	}
}

func (r *repository) SaveMessage(ctx context.Context, msg *Message) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msg.Timestamp = time.Now()
	_, err := r.col.InsertOne(dbCtx, msg)
	return err
}

func (r *repository) GetMessages(ctx context.Context, user1, user2 string) ([]Message, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": user1, "receiver_id": user2},
			{"sender_id": user2, "receiver_id": user1},
		},
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *repository) MarkRead(ctx context.Context, receiverID, senderID string) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"receiver_id": receiverID,
		"sender_id":   senderID,
		"read":        false,
	}

	update := bson.M{
		"$set": bson.M{
			"read":    true,
			"read_at": time.Now().UTC(),
		},
	}

	_, err := r.col.UpdateMany(dbCtx, filter, update)
	return err
}

func (r *repository) GetUnreadMessages(ctx context.Context, receiverID string) ([]Message, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"receiver_id": receiverID, "read": false}
	cur, err := r.col.Find(dbCtx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(dbCtx)

	var msgs []Message
	if err := cur.All(dbCtx, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

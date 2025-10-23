package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	//CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, filter bson.M) (*User, error)
	UpsertUser(ctx context.Context, user *User) error
	AddFriend(ctx context.Context, userID, friendID string) error
	AcceptFriend(ctx context.Context, userID, friendID string) error
}

type repository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) Repository {
	col := client.Database("lumora").Collection("users") // collection names are usually lowercase plural
	fmt.Println("Mongo collection:", col.Name())
	return &repository{collection: col}
}

// func (r *repository) CreateUser(ctx context.Context, user *User) error {
// 	_, err := r.collection.InsertOne(ctx, user)
// 	return err
// }

func (r *repository) GetUser(ctx context.Context, filter bson.M) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpsertUser(ctx context.Context, user *User) error {
	filter := bson.M{"google_id": user.GoogleID} // match by GoogleID
	update := bson.M{"$set": user}               // update all fields
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *repository) AddFriend(ctx context.Context, userID, friendID string) error {
	filter := bson.M{"google_id": userID, "friends.user_id": bson.M{"$ne": friendID}}
	update := bson.M{"$push": bson.M{"friends": bson.M{"user_id": friendID, "status": "pending"}}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("friend request already sent or user not found")
	}
	return nil
}
func (r *repository) AcceptFriend(ctx context.Context, userID, friendID string) error {
	// 1️⃣ Update current user's friend entry to "accepted"
	filterUser := bson.M{"google_id": userID, "friends.user_id": friendID, "friends.status": "pending"}
	updateUser := bson.M{"$set": bson.M{"friends.$.status": "accepted"}}
	result, err := r.collection.UpdateOne(ctx, filterUser, updateUser)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no pending friend request found")
	}

	// 2️⃣ Ensure the other user also has the requester as accepted friend
	filterFriend := bson.M{"google_id": friendID}
	updateFriend := bson.M{
		"$set": bson.M{
			"updated_at": time.Now().Unix(),
		},
		"$addToSet": bson.M{
			"friends": bson.M{
				"user_id": userID,
				"status":  "accepted",
			},
		},
	}
	_, err = r.collection.UpdateOne(ctx, filterFriend, updateFriend)
	if err != nil {
		return fmt.Errorf("failed to update friend's list: %v", err)
	}

	return nil
}

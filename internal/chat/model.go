package chat

import "time"

type Message struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	SenderID   string    `bson:"sender_id" json:"sender_id"`
	ReceiverID string    `bson:"receiver_id" json:"receiver_id"`
	Content    string    `bson:"content" json:"content"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	Delivered  bool      `bson:"delivered" json:"delivered"`
	Read       bool      `bson:"read" json:"read"`
	ReadAt     time.Time `bson:"read_at,omitempty" json:"read_at,omitempty"`
}

type WSMessage struct {
	Type string  `json:"type"` // "chat", "read_ack"
	Data Message `json:"data,omitempty"`
	From string  `json:"from,omitempty"`
	To   string  `json:"to,omitempty"`
}

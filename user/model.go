package user

type User struct {
	ID           string   `bson:"_id,omitempty" json:"id,omitempty"`
	GoogleID     string   `bson:"google_id,omitempty" json:"google_id,omitempty"`
	Email        string   `bson:"email,omitempty" json:"email,omitempty"`
	Name         string   `bson:"name,omitempty" json:"name,omitempty"`
	Picture      string   `bson:"picture,omitempty" json:"picture,omitempty"`
	AccessToken  string   `bson:"access_token,omitempty" json:"access_token,omitempty"`
	RefreshToken string   `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	CreatedAt    int64    `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt    int64    `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Friends      []Friend `bson:"friends,omitempty" json:"friends,omitempty"`
}

type Friend struct {
	UserID string `bson:"user_id" json:"user_id"`
	Status string `bson:"status" json:"status"` // "pending", "accepted", "rejected"
}

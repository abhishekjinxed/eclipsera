package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"lumora/internal/user"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuth struct {
	conf        *oauth2.Config
	logger      *zap.Logger
	userService user.Service //
}
type GoogleUser struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewGoogleAuth(userService user.Service, logger *zap.Logger) *GoogleAuth {
	_ = godotenv.Load("../../.env") // Loads .env if it exists

	return &GoogleAuth{
		conf: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		logger:      logger,
		userService: userService, // ✅ injected
	}
}

func (ga *GoogleAuth) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/google/login", ga.HandleGoogleLogin)
	mux.HandleFunc("/auth/google/callback", ga.HandleGoogleCallback)
	mux.HandleFunc("/auth/me", ga.HandleAuthMe)
	mux.HandleFunc("/auth/logout", ga.HandleAuthLogout)
	//mux.Handle("/auth/me", ValidateJWT(http.HandlerFunc(ga.HandleAuthMe)))

}

func (ga *GoogleAuth) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := ga.conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *GoogleAuth) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	token, err := a.conf.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := a.conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var gUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// ✅ Step 1: Check if user already exists
	existingUser, err := a.userService.GetUser(r.Context(), bson.M{"google_id": gUser.ID})
	if err != nil && err != mongo.ErrNoDocuments {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	//fmt.Print("Google User:", existingUser)
	// ✅ Step 2: Build or update user record
	u := user.User{
		GoogleID:     gUser.ID,
		Email:        gUser.Email,
		Name:         gUser.Name,
		Picture:      gUser.Picture,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		UpdatedAt:    time.Now().Unix(),
	}

	if existingUser != nil {
		// ✅ Preserve existing friends
		u.Friends = existingUser.Friends
		u.CreatedAt = existingUser.CreatedAt
	} else {
		// ✅ New user — add self as default friend
		u.CreatedAt = time.Now().Unix()
		u.Friends = []user.Friend{
			{
				UserID: "10000000000",
				Status: "accepted",
			},
		}
	}

	// ✅ Step 3: Upsert user
	if err := a.userService.UpsertUser(r.Context(), &u); err != nil {
		http.Error(w, "Failed to upsert user", http.StatusInternalServerError)
		return
	}

	// ✅ Step 4: Generate JWT (now includes existing friends)
	jwtToken, err := GenerateJWT(u.Email, u.Name, u.Friends, u.Picture, u.GoogleID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// ✅ Step 5: Redirect to frontend
	redirectURL := fmt.Sprintf("http://localhost:5173/login/success?token=%s", jwtToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (a *GoogleAuth) HandleAuthMe(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Now handle actual GET request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := ValidateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user := map[string]interface{}{
		"name":      claims.Name,
		"email":     claims.Email,
		"picture":   claims.Picture,
		"friends":   claims.Friends, // include friends array
		"google_id": claims.GoogleID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}
func (a *GoogleAuth) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Optionally: invalidate refresh token in DB here

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out successfully",
	})
}

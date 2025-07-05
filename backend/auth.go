package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	_, err = DB.Exec(`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`, req.Username, req.Email, string(hash))
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Database error during registration: %v\n", err)
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// LoginHandler handles user login and session creation.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	var user User
	err := DB.QueryRow(`SELECT id, username, email, password_hash FROM users WHERE username=$1`, req.Username).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Create session
	sessionID := uuid.New().String()
	expires := time.Now().Add(7 * 24 * time.Hour)
	_, err = DB.Exec(`INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`, sessionID, user.ID, expires)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// LogoutHandler clears the session cookie and deletes the session.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		DB.Exec(`DELETE FROM sessions WHERE id=$1`, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// MeHandler returns the current logged-in user's info.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"user":   user,
	})
}

// SessionMiddleware loads the user from the session cookie, if present.
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		var user User
		var expires time.Time
		err = DB.QueryRow(`
			SELECT users.id, users.username, users.email, users.password_hash, sessions.expires_at
			FROM sessions
			JOIN users ON sessions.user_id = users.id
			WHERE sessions.id=$1`, cookie.Value).
			Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &expires)
		if err == nil && expires.After(time.Now()) {
			ctx := context.WithValue(r.Context(), "user", &user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
} 
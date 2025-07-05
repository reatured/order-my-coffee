package main

import (
	"encoding/json"
	"net/http"
	"time"
	"database/sql"
)

// OrderHandler handles order creation for both guests and logged-in users.
func OrderHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		CoffeeID int    `json:"coffeeId"`
		Notes    string `json:"notes"`
		Quantity int    `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name == "" || req.Email == "" || req.CoffeeID == 0 || req.Quantity == 0 {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	var userID sql.NullInt64
	user := r.Context().Value("user")
	if user != nil {
		userID.Int64 = int64(user.(*User).ID)
		userID.Valid = true
	}

	_, err := DB.Exec(
		`INSERT INTO orders (user_id, name, email, coffee_id, notes, quantity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, req.Name, req.Email, req.CoffeeID, req.Notes, req.Quantity, time.Now(),
	)
	if err != nil {
		http.Error(w, "Could not save order", http.StatusInternalServerError)
		return
	}

	// Optionally: send confirmation and admin emails here using SendEmail()

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// OrderHandler handles order creation (stub)
func OrderHandlerStub(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("OrderHandler not implemented yet"))
}
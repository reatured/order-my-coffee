package main

import (
	"database/sql"
	"time"
)

// User represents a user in the system.
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
}

// Session represents a login session.
type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

// Order represents a coffee order.
type Order struct {
	ID        int
	UserID    sql.NullInt64 // Nullable for guest orders
	Name      string
	Email     string
	CoffeeID  int
	Notes     string
	Quantity  int
	CreatedAt time.Time
}
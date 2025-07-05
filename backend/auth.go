package main

import (
	"net/http"
)

// RegisterHandler handles user registration (stub)
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("RegisterHandler not implemented yet"))
}

// LoginHandler handles user login (stub)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("LoginHandler not implemented yet"))
}

// LogoutHandler handles user logout (stub)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("LogoutHandler not implemented yet"))
}

// MeHandler returns current user info (stub)
func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("MeHandler not implemented yet"))
}

// SessionMiddleware is a stub middleware
func SessionMiddleware(next http.Handler) http.Handler {
	return next
} 
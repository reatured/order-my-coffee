package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

// Improved CORS middleware for local development and production
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "https://reatured.github.io" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	godotenv.Load(".env")
	InitDB()
	defer DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/register", withCORS(RegisterHandler))
	mux.HandleFunc("/login", withCORS(LoginHandler))
	mux.HandleFunc("/logout", withCORS(LogoutHandler))
	mux.HandleFunc("/me", withCORS(MeHandler))
	mux.HandleFunc("/order", withCORS(OrderHandler))
	// Add your /coffees handler as needed, e.g.:
	// mux.HandleFunc("/coffees", withCORS(CoffeesHandler))

	// Wrap with session middleware
	handler := SessionMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running at http://localhost:%s/\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
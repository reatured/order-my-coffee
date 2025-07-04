package main

import (
	"database/sql"
	"encoding/json"
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
	
	// Debug route to check database
	mux.HandleFunc("/debug/db", withCORS(func(w http.ResponseWriter, r *http.Request) {
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}
		
		// Check table structure
		rows, err := DB.Query(`
			SELECT column_name, data_type, is_nullable, column_default
			FROM information_schema.columns 
			WHERE table_name = 'users' 
			ORDER BY ordinal_position
		`)
		if err != nil {
			http.Error(w, fmt.Sprintf("Schema query error: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		
		var columns []map[string]interface{}
		for rows.Next() {
			var colName, dataType, isNullable, colDefault sql.NullString
			rows.Scan(&colName, &dataType, &isNullable, &colDefault)
			columns = append(columns, map[string]interface{}{
				"name": colName.String,
				"type": dataType.String,
				"nullable": isNullable.String,
				"default": colDefault.String,
			})
		}
		
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_count": count,
			"status": "Database connected successfully",
			"table_structure": columns,
		})
	}))
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
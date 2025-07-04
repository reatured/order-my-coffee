// main.go - Coffee Order Backend API
//
// This Go program provides a simple backend for a coffee ordering app.
// It supports listing available coffees, accepting orders, saving them,
// and sending email notifications to both the admin and the customer.
//
// Key concepts: HTTP handlers, JSON, file I/O, environment variables, SMTP email, CORS, Go structs.

package main

import (
    "encoding/json"   // For encoding/decoding JSON
    "fmt"             // For formatted I/O
    "io/ioutil"       // For reading/writing files
    "log"             // For logging errors/info
    "net/http"        // For HTTP server
    "os"              // For file and environment variable access
    "net/smtp"        // For sending emails
    "github.com/joho/godotenv" // For loading .env files
)

// Coffee represents a coffee drink in the menu.
type Coffee struct {
    ID   int    `json:"id"`   // Unique ID for the coffee
    Name string `json:"name"` // Name of the coffee
}

// Order represents a customer's order.
type Order struct {
    Name     string `json:"name"`      // Customer's name
    CoffeeID int    `json:"coffeeId"` // ID of the coffee ordered
    Notes    string `json:"notes"`    // Any notes for the order
    Email    string `json:"email,omitempty"` // Customer's email (optional)
}

// withCORS is a middleware that adds CORS headers to HTTP responses.
// This allows the frontend (on a different domain) to call the backend API.
func withCORS(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        h(w, r)
    }
}

// getCoffeesHandler handles GET requests to /coffees.
// It reads the coffee menu from coffees.json and returns it as JSON.
func getCoffeesHandler(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadFile("coffees.json")
    if err != nil {
        http.Error(w, "Could not read coffees", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

// orderHandler handles POST requests to /order.
// It decodes the order, saves it, and sends emails to admin and customer.
func orderHandler(w http.ResponseWriter, r *http.Request) {
    var order Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, "Invalid order", http.StatusBadRequest)
        return
    }

    // Optionally save order to orders.json
    orders := []Order{}
    ordersFile := "orders.json"
    if _, err := os.Stat(ordersFile); err == nil {
        data, _ := ioutil.ReadFile(ordersFile)
        json.Unmarshal(data, &orders)
    }
    orders = append(orders, order)
    ordersData, _ := json.MarshalIndent(orders, "", "  ")
    ioutil.WriteFile(ordersFile, ordersData, 0644)

    // Send email to admin
    emailSent, emailErr := sendOrderEmail(order)

    // Send confirmation email to customer if email is provided
    customerEmailSent := false
    customerEmailErr := error(nil)
    if order.Email != "" {
        customerEmailSent, customerEmailErr = sendConfirmationEmail(order)
    }

    w.Header().Set("Content-Type", "application/json")
    resp := map[string]interface{}{
        "status":    "ok",
        "adminEmailSent": emailSent,
        "customerEmailSent": customerEmailSent,
    }
    if emailErr != nil {
        resp["adminError"] = emailErr.Error()
    }
    if customerEmailErr != nil {
        resp["customerError"] = customerEmailErr.Error()
    }
    json.NewEncoder(w).Encode(resp)
}

// sendOrderEmail sends an email to the admin with the order details.
// Returns true if sent successfully, or false and an error if failed.
func sendOrderEmail(order Order) (bool, error) {
    // Load coffee name from coffees.json
    data, _ := ioutil.ReadFile("coffees.json")
    var coffees []Coffee
    json.Unmarshal(data, &coffees)
    coffeeName := "Unknown"
    for _, c := range coffees {
        if c.ID == order.CoffeeID {
            coffeeName = c.Name
            break
        }
    }

    // Email config (should use environment variables for security)
    from := os.Getenv("SMTP_USER")         // Sender email
    pass := os.Getenv("SMTP_PASSWORD")     // Sender password
    to := os.Getenv("CONTACT_EMAIL")       // Admin/recipient email
    smtpHost := os.Getenv("SMTP_HOST")     // SMTP server host
    smtpPort := os.Getenv("SMTP_PORT")     // SMTP server port

    subject := "New Coffee Order"
    body := fmt.Sprintf("Name: %s\nCoffee: %s\nNotes: %s", order.Name, coffeeName, order.Notes)
    if order.Email != "" {
        body += fmt.Sprintf("\nCustomer Email: %s", order.Email)
    }
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    // Set up authentication and send the email
    auth := smtp.PlainAuth("", from, pass, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        log.Println("Failed to send email:", err)
        return false, err
    }
    return true, nil
}

// sendConfirmationEmail sends a confirmation email to the customer.
// Returns true if sent successfully, or false and an error if failed.
func sendConfirmationEmail(order Order) (bool, error) {
    // Load coffee name from coffees.json
    data, _ := ioutil.ReadFile("coffees.json")
    var coffees []Coffee
    json.Unmarshal(data, &coffees)
    coffeeName := "Unknown"
    for _, c := range coffees {
        if c.ID == order.CoffeeID {
            coffeeName = c.Name
            break
        }
    }

    from := os.Getenv("SMTP_USER")         // Sender email
    pass := os.Getenv("SMTP_PASSWORD")     // Sender password
    smtpHost := os.Getenv("SMTP_HOST")     // SMTP server host
    smtpPort := os.Getenv("SMTP_PORT")     // SMTP server port

    to := order.Email
    subject := "Your Coffee Order Confirmation"
    body := fmt.Sprintf("Hi %s,\n\nThank you for your order!\n\nOrder Details:\nCoffee: %s\nNotes: %s\n\nWe will process your order soon.\n\nBest regards,\nCoffee Shop", order.Name, coffeeName, order.Notes)
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    // Set up authentication and send the email
    auth := smtp.PlainAuth("", from, pass, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        log.Println("Failed to send confirmation email to customer:", err)
        return false, err
    }
    return true, nil
}

// main is the entry point of the program.
// It loads environment variables, sets up HTTP routes, and starts the server.
func main() {
    godotenv.Load(".env") // Load environment variables from .env file
    http.HandleFunc("/coffees", withCORS(getCoffeesHandler)) // GET /coffees
    http.HandleFunc("/order", withCORS(orderHandler))         // POST /order
    fmt.Println("Server running at http://localhost:8080/")
    log.Fatal(http.ListenAndServe(":8080", nil)) // Start HTTP server
}
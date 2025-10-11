package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Rsvp represents a single RSVP response
type Rsvp struct {
	ID         int
	Name       string
	Email      string
	Phone      string
	WillAttend bool
	CreatedAt  time.Time
}

// Global variables
var (
	db        *sql.DB
	templates = make(map[string]*template.Template, 5)
	dbMutex   sync.RWMutex // Mutex for database operations
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Phone validation regex (allows digits, spaces, dashes, parentheses, plus)
var phoneRegex = regexp.MustCompile(`^[\d\s\-\+\(\)]+$`)

// initDatabase initializes the SQLite database and creates tables
func initDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./rsvp.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Create the rsvps table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS rsvps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		phone TEXT NOT NULL,
		will_attend BOOLEAN NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_email ON rsvps(email);
	CREATE INDEX IF NOT EXISTS idx_will_attend ON rsvps(will_attend);
	`

	if _, err = db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// loadTemplates loads all HTML templates
func loadTemplates() {
	templateNames := []string{"welcome", "form", "thanks", "sorry", "list"}
	for _, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			log.Printf("Loaded template: %s", name)
		} else {
			log.Printf("Failed to load template %s: %v", name, err)
		}
	}
}

// validateEmail validates email format
func validateEmail(email string) (bool, string) {
	email = strings.TrimSpace(email)
	if email == "" {
		return false, "Email address is required"
	}
	if !emailRegex.MatchString(email) {
		return false, "Please enter a valid email address"
	}
	return true, ""
}

// validatePhone validates phone number
func validatePhone(phone string) (bool, string) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false, "Phone number is required"
	}
	if !phoneRegex.MatchString(phone) {
		return false, "Please enter a valid phone number"
	}
	// Check if phone has at least 10 digits
	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	if len(digits) < 10 {
		return false, "Phone number must contain at least 10 digits"
	}
	return true, ""
}

// validateName validates name field
func validateName(name string) (bool, string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, "Name is required"
	}
	if len(name) < 2 {
		return false, "Name must be at least 2 characters long"
	}
	return true, ""
}

// checkDuplicateEmail checks if email already exists
func checkDuplicateEmail(email string) (bool, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM rsvps WHERE LOWER(email) = LOWER(?)", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// saveRsvp saves an RSVP to the database
func saveRsvp(rsvp *Rsvp) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	result, err := db.Exec(
		"INSERT INTO rsvps (name, email, phone, will_attend) VALUES (?, ?, ?, ?)",
		rsvp.Name, rsvp.Email, rsvp.Phone, rsvp.WillAttend,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	rsvp.ID = int(id)
	return nil
}

// getAllRsvps retrieves all RSVPs from the database
func getAllRsvps() ([]*Rsvp, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	rows, err := db.Query("SELECT id, name, email, phone, will_attend, created_at FROM rsvps ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rsvps []*Rsvp
	for rows.Next() {
		var rsvp Rsvp
		err := rows.Scan(&rsvp.ID, &rsvp.Name, &rsvp.Email, &rsvp.Phone, &rsvp.WillAttend, &rsvp.CreatedAt)
		if err != nil {
			return nil, err
		}
		rsvps = append(rsvps, &rsvp)
	}

	return rsvps, rows.Err()
}

// formData holds form data and validation errors
type formData struct {
	*Rsvp
	Errors []string
}

// welcomeHandler handles the home page
func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	if err := templates["welcome"].Execute(writer, nil); err != nil {
		log.Printf("Error executing welcome template: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// listHandler handles the guest list page
func listHandler(writer http.ResponseWriter, request *http.Request) {
	rsvps, err := getAllRsvps()
	if err != nil {
		log.Printf("Error retrieving RSVPs: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := templates["list"].Execute(writer, rsvps); err != nil {
		log.Printf("Error executing list template: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// formHandler handles the RSVP form (both GET and POST)
func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		// Show empty form
		if err := templates["form"].Execute(writer, formData{
			Rsvp:   &Rsvp{},
			Errors: []string{},
		}); err != nil {
			log.Printf("Error executing form template: %v", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
	} else if request.Method == http.MethodPost {
		handleFormSubmission(writer, request)
	} else {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handleFormSubmission processes the form submission
func handleFormSubmission(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract and sanitize form data
	name := strings.TrimSpace(request.Form.Get("name"))
	email := strings.TrimSpace(request.Form.Get("email"))
	phone := strings.TrimSpace(request.Form.Get("phone"))
	willAttendStr := request.Form.Get("willAttend")

	responseData := Rsvp{
		Name:       name,
		Email:      email,
		Phone:      phone,
		WillAttend: willAttendStr == "true",
	}

	// Validate all fields
	errors := []string{}

	if valid, msg := validateName(name); !valid {
		errors = append(errors, msg)
	}

	if valid, msg := validateEmail(email); !valid {
		errors = append(errors, msg)
	} else {
		// Check for duplicate email
		duplicate, err := checkDuplicateEmail(email)
		if err != nil {
			log.Printf("Error checking duplicate email: %v", err)
			errors = append(errors, "An error occurred. Please try again.")
		} else if duplicate {
			errors = append(errors, "This email address has already been used for an RSVP")
		}
	}

	if valid, msg := validatePhone(phone); !valid {
		errors = append(errors, msg)
	}

	// If there are validation errors, show form again with errors
	if len(errors) > 0 {
		if err := templates["form"].Execute(writer, formData{
			Rsvp:   &responseData,
			Errors: errors,
		}); err != nil {
			log.Printf("Error executing form template: %v", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Save to database
	if err := saveRsvp(&responseData); err != nil {
		log.Printf("Error saving RSVP: %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			errors = append(errors, "This email address has already been used for an RSVP")
			if err := templates["form"].Execute(writer, formData{
				Rsvp:   &responseData,
				Errors: errors,
			}); err != nil {
				log.Printf("Error executing form template: %v", err)
			}
		} else {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("New RSVP saved: %s (%s) - Attending: %v", responseData.Name, responseData.Email, responseData.WillAttend)

	// Show appropriate thank you page
	if responseData.WillAttend {
		if err := templates["thanks"].Execute(writer, responseData.Name); err != nil {
			log.Printf("Error executing thanks template: %v", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		if err := templates["sorry"].Execute(writer, responseData.Name); err != nil {
			log.Printf("Error executing sorry template: %v", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// healthHandler handles health check requests
func healthHandler(writer http.ResponseWriter, request *http.Request) {
	// Check database connection
	if err := db.Ping(); err != nil {
		log.Printf("Health check failed: database error: %v", err)
		http.Error(writer, "Database connection failed", http.StatusServiceUnavailable)
		return
	}

	writer.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(writer, "OK - Server is running\nDatabase: Connected")
}

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next(w, r)
		log.Printf("Completed in %v", time.Since(start))
	}
}

// main is the entry point of the application
func main() {
	// Initialize database
	if err := initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Load templates
	loadTemplates()

	// Setup static file server for CSS and JS
	fs := http.FileServer(http.Dir("."))
	http.Handle("/styles.css", fs)
	http.Handle("/app.js", fs)

	// Setup routes with logging middleware
	http.HandleFunc("/", loggingMiddleware(welcomeHandler))
	http.HandleFunc("/list", loggingMiddleware(listHandler))
	http.HandleFunc("/form", loggingMiddleware(formHandler))
	http.HandleFunc("/health", healthHandler)

	// Get port from environment variable (Railway)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // fallback for local development
	}

	// Start server
	log.Printf("Starting server on port %s", port)
	log.Printf("Visit http://localhost:%s to view the application", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

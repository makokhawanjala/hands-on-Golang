// Package main is the entry point for our web application
package main

// Import necessary packages for our web server
import (
	"fmt"           // For printing to console and string formatting
	"html/template" // For parsing and executing HTML templates
	"net/http"      // For creating HTTP server and handling requests
	"os"            // For reading environment variables
)

// Define a struct to represent an RSVP response
// Structs in Go group related data together
type Rsvp struct {
	Name, Email, Phone string // Guest's contact information
	WillAttend         bool   // Whether the guest will attend (must be capitalized to be exported/accessible from templates)
}

// Global variables to store our application data
// In a real application, you'd use a database instead of in-memory storage
var responses = make([]*Rsvp, 0, 10)                   // Slice to store RSVP responses (initially empty, capacity of 10)
var templates = make(map[string]*template.Template, 3) // Map to cache parsed templates for better performance

// loadTemplates function reads and parses all HTML template files
// This is called once at startup to avoid parsing templates on every request
func loadTemplates() {
	// Array of template names that correspond to our HTML files
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}

	// Loop through each template name
	for index, name := range templateNames {
		// ParseFiles reads the layout.html (base template) and the specific template file
		// layout.html contains the common HTML structure, individual files define the "body" block
		t, err := template.ParseFiles("layout.html", name+".html")

		// Check if parsing was successful
		if err == nil {
			// Store the parsed template in our map for later use
			templates[name] = t
			fmt.Println("Loaded template", index, name) // Debug output to console
		} else {
			// If template parsing fails, crash the program (panic)
			// This ensures we catch template errors at startup, not during user requests
			panic(err)
		}
	}
}

// welcomeHandler serves the welcome page (home page)
// HTTP handlers in Go take a ResponseWriter (to send response) and Request (incoming request data)
func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	// Execute the welcome template and send it to the client
	// The second parameter (nil) means we're not passing any data to the template
	templates["welcome"].Execute(writer, nil)
}

// listHandler shows the list of people who are attending the party
func listHandler(writer http.ResponseWriter, request *http.Request) {
	// Execute the list template and pass our responses slice as data
	// The template will loop through this data to display attendees
	templates["list"].Execute(writer, responses)
}

// formData struct is used to pass both RSVP data and validation errors to the form template
type formData struct {
	*Rsvp           // Embedded struct - includes Name, Email, Phone, WillAttend fields
	Errors []string // Slice of error messages to display to the user
}

// formHandler handles both displaying the RSVP form (GET) and processing submissions (POST)
func formHandler(writer http.ResponseWriter, request *http.Request) {
	// Check the HTTP method to determine what to do
	if request.Method == http.MethodGet {
		// GET request: Show empty form
		templates["form"].Execute(writer, formData{
			Rsvp: &Rsvp{}, Errors: []string{}, // Empty RSVP data and no errors
		})
	} else if request.Method == http.MethodPost {
		// POST request: Process form submission

		// ParseForm() parses the form data from the request body
		// This populates request.Form with the submitted values
		request.ParseForm()

		// Safely extract form values with existence and length checks
		// This prevents "index out of range" panics if fields are missing or empty

		name := ""
		if vals, exists := request.Form["name"]; exists && len(vals) > 0 {
			name = vals[0] // Take the first value (forms can have multiple values for same name)
		}

		email := ""
		if vals, exists := request.Form["email"]; exists && len(vals) > 0 {
			email = vals[0]
		}

		phone := ""
		if vals, exists := request.Form["phone"]; exists && len(vals) > 0 {
			phone = vals[0]
		}

		// Handle the attendance dropdown (willattend field)
		willAttend := false
		if vals, exists := request.Form["willattend"]; exists && len(vals) > 0 {
			// Convert string "true" to boolean true, anything else becomes false
			willAttend = vals[0] == "true"
		}

		// Create a new RSVP struct with the extracted values
		responseData := Rsvp{
			Name:       name,
			Email:      email,
			Phone:      phone,
			WillAttend: willAttend,
		}

		// Add the response to our global responses slice
		// In a real app, this would be saved to a database
		responses = append(responses, &responseData)

		// Show different thank you pages based on attendance choice
		if responseData.WillAttend {
			// Guest is attending - show thanks page with their name
			templates["thanks"].Execute(writer, responseData.Name)
		} else {
			// Guest is not attending - show sorry page with their name
			templates["sorry"].Execute(writer, responseData.Name)
		}
	}
}

// main function is the entry point of our program
func main() {
	// Load and parse all templates at startup
	loadTemplates()

	// Set up HTTP routes - map URL patterns to handler functions
	http.HandleFunc("/", welcomeHandler)  // Home page
	http.HandleFunc("/list", listHandler) // Attendee list page
	http.HandleFunc("/form", formHandler) // RSVP form page

	// Determine which port to run the server on
	// Cloud platforms set the PORT environment variable, locally we default to 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // Default port for local development
	}

	// Print startup message to console
	fmt.Printf("Server starting on port %s\n", port)

	// Start the HTTP server and listen for incoming requests
	// ListenAndServe blocks here - the program will run until interrupted
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		// If server fails to start, print the error
		fmt.Println(err)
	}
}

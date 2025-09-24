package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Define a struct to group a set of related values
// This holds all the info we need for each person's RSVP
type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool // Changed from willAttend to WillAttend (exported)
}

// Initialize a new slice with make function
// and give it initial size and capacity
// This will store all the RSVP responses we receive
var responses = make([]*Rsvp, 0, 10)

// This map will hold our HTML templates so we can reuse them
var templates = make(map[string]*template.Template, 3)

// Load all the HTML templates from files
func loadTemplates() {
	// Array of template names that match our HTML files
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	// Loop through each template name
	for index, name := range templateNames {
		// Try to parse the layout.html and the specific template file
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			// If successful, store the template in our map
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			// If there's an error, crash the program
			// panic(err)
			fmt.Printf("Failed to load template %s: %v\n", name, err)
		}
	}
}

// Handle requests to the home page
func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	// Execute the welcome template and send it to the browser
	templates["welcome"].Execute(writer, nil)
}

// Handle requests to view all RSVP responses
func listHandler(writer http.ResponseWriter, request *http.Request) {
	// Execute the list template and pass it all the responses we've collected
	templates["list"].Execute(writer, responses)
}

// Struct to hold form data and any validation errors
type formData struct {
	*Rsvp           // Embed the Rsvp struct
	Errors []string // Slice to hold error messages
}

// Handle both showing the form and processing form submissions
func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		// If it's a GET request, show the empty form
		templates["form"].Execute(writer, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		// If it's a POST request, process the form data
		request.ParseForm() // Parse the form data from the request

		// Create a new Rsvp struct with the form data
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willAttend"][0] == "true", // Convert string to boolean
		}

		// Check for validation errors
		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Please Enter your name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter your email address")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter your phone number")
		}

		if len(errors) > 0 {
			// If there are errors, show the form again with error messages
			templates["form"].Execute(writer, formData{
				Rsvp: &responseData, Errors: errors,
			})
		} else {
			// If no errors, save the response and show appropriate thank you page
			responses = append(responses, &responseData)
			if responseData.WillAttend {
				// Show thanks page if they're attending
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				// Show sorry page if they're not attending
				templates["sorry"].Execute(writer, responseData.Name)
			}

		}

	}
}

func healthHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "OK - Server is running")
}

// Main function - entry point of the program
func main() {
	// Load all the HTML templates first
	loadTemplates()

	// Set up URL routes and their corresponding handler functions
	http.HandleFunc("/", welcomeHandler)  // Home page
	http.HandleFunc("/list", listHandler) // View all responses
	http.HandleFunc("/form", formHandler) // RSVP form
	http.HandleFunc("/health", healthHandler)

	// Start the web server on port 5000
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		// Print any errors if the server fails to start
		fmt.Println(err)
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// MustGet - Original function from screenshot
// This function will panic if any error occurs
func MustGet(url string) string {
	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		panic(err) // This will crash the program on error
	}

	// Ensure response body is closed when function exits
	// This prevents memory leaks
	defer func() { _ = resp.Body.Close() }()

	// Read entire response body into memory
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Convert byte slice to string and return
	return string(body)
}

// SafeGet - Better version that handles errors gracefully
func SafeGet(url string) (string, error) {
	// Create HTTP client with timeout to prevent hanging
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make HTTP GET request with custom client
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}

	// CRITICAL: Always close the response body
	defer func() { _ = resp.Body.Close() }()

	// Check if the HTTP status is successful (200-299)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// GetWithHeaders - Advanced version with custom headers
func GetWithHeaders(url string, headers map[string]string) (string, error) {
	// Create new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create client with timeout
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and return body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}

func main() {
	// Example 1: Using MustGet (dangerous - will panic on error)
	fmt.Println("=== Example 1: MustGet (Dangerous) ===")
	// Uncomment to test, but it might crash!
	// result := MustGet("https://httpbin.org/get")
	// fmt.Println("Result:", result[:100] + "...")

	// Example 2: Using SafeGet (recommended)
	fmt.Println("\n=== Example 2: SafeGet (Recommended) ===")
	result, err := SafeGet("https://httpbin.org/get")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Success! Response length: %d characters\n", len(result))
		fmt.Printf("First 100 chars: %s...\n", result[:100])
	}

	// Example 3: Using GetWithHeaders
	fmt.Println("\n=== Example 3: GetWithHeaders ===")
	headers := map[string]string{
		"User-Agent": "MyGoApp/1.0",
		"Accept":     "application/json",
	}

	result, err = GetWithHeaders("https://httpbin.org/headers", headers)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Response with custom headers:\n%s\n", result)
	}
}

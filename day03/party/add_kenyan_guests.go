package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Kenyan names database
var kenyanFirstNames = []string{
	// Male names
	"Kamau", "Mwangi", "Otieno", "Ochieng", "Kipchoge", "Kimani", "Njoroge",
	"Karanja", "Mutua", "Wanjiru", "Kiprotich", "Onyango", "Omondi", "Wekesa",
	"Korir", "Cheruiyot", "Rotich", "Kiprop", "Kiptoo", "Maina", "Gitau",
	"Mbugua", "Ndungu", "Kariuki", "Wairimu", "Chege", "Kabiru", "Mugo",
	"Wafula", "Mulongo", "Barasa", "Makau", "Musyoka", "Mutiso", "Kilonzo",
	"Kipkemboi", "Kigen", "Biwott", "Sang", "Koech", "Langat", "Kirui",
	// Female names
	"Achieng", "Atieno", "Akinyi", "Adhiambo", "Wanjiku", "Njeri", "Nyambura",
	"Wangari", "Wambui", "Wairimu", "Njoki", "Nyokabi", "Chemutai", "Chepkoech",
	"Jebet", "Chepkemoi", "Faith", "Grace", "Mary", "Jane", "Lucy", "Susan",
	"Agnes", "Rose", "Anne", "Elizabeth", "Margaret", "Catherine", "Joyce",
	"Mercy", "Esther", "Sarah", "Rebecca", "Ruth", "Rachel", "Lydia",
}

var kenyanLastNames = []string{
	"Kamau", "Mwangi", "Otieno", "Ochieng", "Kipchoge", "Kimani", "Njoroge",
	"Karanja", "Mutua", "Wanjiru", "Kiprotich", "Onyango", "Omondi", "Wekesa",
	"Korir", "Cheruiyot", "Rotich", "Kiprop", "Kiptoo", "Maina", "Gitau",
	"Mbugua", "Ndungu", "Kariuki", "Wairimu", "Chege", "Kabiru", "Mugo",
	"Wafula", "Mulongo", "Barasa", "Makau", "Musyoka", "Mutiso", "Kilonzo",
	"Kipkemboi", "Kigen", "Biwott", "Sang", "Koech", "Langat", "Kirui",
	"Nyambura", "Wangari", "Wambui", "Njeri", "Njoki", "Nyokabi", "Achieng",
	"Atieno", "Akinyi", "Adhiambo", "Chemutai", "Chepkoech", "Jebet",
	"Chepkemoi", "Kiptanui", "Kiplagat", "Kibet", "Chepchirchir", "Jepkosgei",
}

// Kenyan email domains
var emailDomains = []string{
	"gmail.com", "yahoo.com", "outlook.com", "hotmail.com",
	"safaricom.co.ke", "icloud.com", "live.com", "protonmail.com",
}

// Kenyan mobile prefixes
var kenyanPrefixes = []string{
	// Safaricom prefixes (7xx)
	"710", "711", "712", "713", "714", "715", "716", "717", "718", "719",
	"720", "721", "722", "723", "724", "725", "726", "727", "728", "729",
	"740", "741", "742", "743", "745", "746", "748",
	"757", "758", "759",
	"768", "769",
	"790", "791", "792", "793", "794", "795", "796", "797", "798", "799",
	// Airtel prefixes (7xx)
	"730", "731", "732", "733", "734", "735", "736", "737", "738", "739",
	"750", "751", "752", "753", "754", "755", "756",
	"762", "763", "764", "765", "766", "767",
	"780", "781", "782", "783", "784", "785", "786", "787", "788", "789",
	// Telekom prefixes (77x)
	"770", "771", "772", "773", "774", "775", "776", "777", "778", "779",
}

// Rsvp represents a single RSVP response
type Rsvp struct {
	Name       string
	Email      string
	Phone      string
	WillAttend bool
	CreatedAt  time.Time
}

// generateKenyanName generates a random Kenyan name
func generateKenyanName() string {
	firstName := kenyanFirstNames[rand.Intn(len(kenyanFirstNames))]
	lastName := kenyanLastNames[rand.Intn(len(kenyanLastNames))]

	// Sometimes add middle name
	if rand.Float32() < 0.3 {
		middleName := kenyanFirstNames[rand.Intn(len(kenyanFirstNames))]
		return fmt.Sprintf("%s %s %s", firstName, middleName, lastName)
	}

	return fmt.Sprintf("%s %s", firstName, lastName)
}

// generateKenyanEmail generates a realistic Kenyan email
func generateKenyanEmail(name string) string {
	domain := emailDomains[rand.Intn(len(emailDomains))]

	// Clean name for email
	cleanName := name
	cleanName = fmt.Sprintf("%s%d", cleanName, rand.Intn(100))

	// Create email variations
	variations := []string{
		fmt.Sprintf("%s@%s", cleanName, domain),
		fmt.Sprintf("%s.ke@%s", cleanName, domain),
		fmt.Sprintf("%s_%d@%s", cleanName, rand.Intn(999), domain),
	}

	email := variations[rand.Intn(len(variations))]

	// Replace spaces with dots or underscores
	if rand.Float32() < 0.5 {
		email = replaceSpaces(email, ".")
	} else {
		email = replaceSpaces(email, "_")
	}

	return toLowercase(email)
}

// generateKenyanPhone generates a realistic Kenyan phone number
func generateKenyanPhone() string {
	// Pick a random valid prefix
	prefix := kenyanPrefixes[rand.Intn(len(kenyanPrefixes))]

	// Generate 6 random digits for the rest of the number
	digit1 := rand.Intn(10)
	digit2 := rand.Intn(10)
	digit3 := rand.Intn(10)
	digit4 := rand.Intn(10)
	digit5 := rand.Intn(10)
	digit6 := rand.Intn(10)

	// Pick a random format and build the number
	formatType := rand.Intn(6)

	switch formatType {
	case 0:
		// +254 712 345 678
		return fmt.Sprintf("+254 %s %d%d%d %d%d%d", prefix, digit1, digit2, digit3, digit4, digit5, digit6)
	case 1:
		// 0712 345 678
		return fmt.Sprintf("0%s %d%d%d %d%d%d", prefix, digit1, digit2, digit3, digit4, digit5, digit6)
	case 2:
		// +254-712-345-678
		return fmt.Sprintf("+254-%s-%d%d%d-%d%d%d", prefix, digit1, digit2, digit3, digit4, digit5, digit6)
	case 3:
		// (+254) 712 345678
		return fmt.Sprintf("(+254) %s %d%d%d%d%d%d", prefix, digit1, digit2, digit3, digit4, digit5, digit6)
	case 4:
		// +254 712345678
		return fmt.Sprintf("+254 %s%d%d%d%d%d%d", prefix, digit1, digit2, digit3, digit4, digit5, digit6)
	default:
		// 0712345678
		return fmt.Sprintf("0%s%d%d%d%d%d%d", prefix, digit1, digit2, digit3, digit4, digit5, digit6)
	}
}

// Helper functions
func replaceSpaces(s, replacement string) string {
	result := ""
	for _, char := range s {
		if char == ' ' {
			result += replacement
		} else {
			result += string(char)
		}
	}
	return result
}

func toLowercase(s string) string {
	result := ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(char + 32)
		} else {
			result += string(char)
		}
	}
	return result
}

// saveRsvp saves an RSVP to the database
func saveRsvp(db *sql.DB, rsvp *Rsvp) error {
	_, err := db.Exec(
		"INSERT INTO rsvps (name, email, phone, will_attend, created_at) VALUES (?, ?, ?, ?, ?)",
		rsvp.Name, rsvp.Email, rsvp.Phone, rsvp.WillAttend, rsvp.CreatedAt,
	)
	return err
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Open database
	db, err := sql.Open("sqlite3", "./rsvp.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")
	log.Println("Starting to add 100 Kenyan guests...")

	successCount := 0
	duplicateCount := 0
	errorCount := 0

	// Generate and add 100 guests
	for i := 0; i < 100; i++ {
		name := generateKenyanName()
		email := generateKenyanEmail(name)
		phone := generateKenyanPhone()
		willAttend := rand.Float32() < 0.85 // 85% will attend

		// Random creation time (within last 30 days)
		daysAgo := rand.Intn(30)
		hoursAgo := rand.Intn(24)
		minutesAgo := rand.Intn(60)
		createdAt := time.Now().AddDate(0, 0, -daysAgo).
			Add(time.Duration(-hoursAgo) * time.Hour).
			Add(time.Duration(-minutesAgo) * time.Minute)

		rsvp := &Rsvp{
			Name:       name,
			Email:      email,
			Phone:      phone,
			WillAttend: willAttend,
			CreatedAt:  createdAt,
		}

		err := saveRsvp(db, rsvp)
		if err != nil {
			if contains(err.Error(), "UNIQUE constraint failed") {
				duplicateCount++
				// Try again with a different email
				i--
				continue
			} else {
				log.Printf("Error saving RSVP #%d: %v", i+1, err)
				errorCount++
			}
		} else {
			successCount++
			status := "✓ Attending"
			if !willAttend {
				status = "✗ Not Attending"
			}
			log.Printf("[%d/100] Added: %-30s | %-40s | %-20s | %s",
				successCount, name, email, phone, status)
		}
	}

	// Summary
	log.Println("\n" + strings.Repeat("=", 70))
	log.Printf("✅ Successfully added: %d guests", successCount)
	log.Printf("⚠️  Duplicates skipped: %d", duplicateCount)
	log.Printf("❌ Errors: %d", errorCount)
	log.Println(strings.Repeat("=", 70))

	// Show attendance stats
	var attending, notAttending int
	err = db.QueryRow("SELECT COUNT(*) FROM rsvps WHERE will_attend = 1").Scan(&attending)
	if err == nil {
		db.QueryRow("SELECT COUNT(*) FROM rsvps WHERE will_attend = 0").Scan(&notAttending)
		log.Printf("\n📊 Database Stats:")
		log.Printf("   Total Guests: %d", attending+notAttending)
		log.Printf("   Attending: %d (%.1f%%)", attending, float64(attending)/float64(attending+notAttending)*100)
		log.Printf("   Not Attending: %d (%.1f%%)", notAttending, float64(notAttending)/float64(attending+notAttending)*100)
	}

	log.Println("\n✨ Done! Your RSVP database is now populated with Kenyan guests!")
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

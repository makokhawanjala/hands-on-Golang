package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ============================================
// SIMULATED SLOW DATABASE
// ============================================

// Database simulates a slow database with realistic delays
type Database struct {
	// In real life, this would be PostgreSQL, MySQL, etc.
	// We're using a simple map to simulate it
	users map[string]string
}

// NewDatabase creates a new database with sample users
func NewDatabase() *Database {
	return &Database{
		users: map[string]string{
			"user:1": "Alice Johnson - Software Engineer",
			"user:2": "Bob Smith - Product Manager",
			"user:3": "Charlie Brown - Designer",
			"user:4": "Diana Prince - CEO",
			"user:5": "Ethan Hunt - Security Specialist",
		},
	}
}

// GetUser simulates a slow database query
func (db *Database) GetUser(userID string) (string, error) {
	// Simulate database latency (network delay, disk I/O, query processing)
	time.Sleep(2 * time.Second) // This is the "slow" part!

	user, exists := db.users[userID]
	if !exists {
		return "", fmt.Errorf("user not found: %s", userID)
	}

	return user, nil
}

// ============================================
// WITHOUT CACHING - Direct Database Access
// ============================================

// UserServiceWithoutCache always hits the database
type UserServiceWithoutCache struct {
	db *Database
}

func NewUserServiceWithoutCache(db *Database) *UserServiceWithoutCache {
	return &UserServiceWithoutCache{db: db}
}

func (s *UserServiceWithoutCache) GetUser(userID string) (string, error) {
	// Every single request goes to the database
	// No caching = always slow
	return s.db.GetUser(userID)
}

// ============================================
// WITH CACHING - Cache-Aside Pattern
// ============================================

// UserServiceWithCache uses Redis to cache database results
type UserServiceWithCache struct {
	db          *Database
	redisClient *redis.Client
	cacheTTL    time.Duration // How long data stays in cache
}

func NewUserServiceWithCache(db *Database, redisClient *redis.Client) *UserServiceWithCache {
	return &UserServiceWithCache{
		db:          db,
		redisClient: redisClient,
		cacheTTL:    5 * time.Minute, // Cache expires after 5 minutes
	}
}

func (s *UserServiceWithCache) GetUser(userID string) (string, error) {
	// context.Background() is like a "request context" - don't worry about it for now
	ctx := context.Background()

	// STEP 1: Try to get from cache first (fast path)
	cachedUser, err := s.redisClient.Get(ctx, userID).Result()

	if err == nil {
		// Cache HIT! We found it in Redis
		fmt.Printf("✅ CACHE HIT for %s\n", userID)
		return cachedUser, nil
	}

	// Cache MISS - data not in Redis (or Redis error)
	fmt.Printf("❌ CACHE MISS for %s - fetching from database...\n", userID)

	// STEP 2: Get from database (slow path)
	user, err := s.db.GetUser(userID)
	if err != nil {
		return "", err
	}

	// STEP 3: Store in cache for next time (write to Redis)
	err = s.redisClient.Set(ctx, userID, user, s.cacheTTL).Err()
	if err != nil {
		// Log error but don't fail the request
		// The user still gets their data even if caching fails
		log.Printf("Warning: Failed to cache user %s: %v\n", userID, err)
	}

	return user, nil
}

// ============================================
// DEMO & BENCHMARKING
// ============================================

func main() {
	fmt.Println("🚀 PROJECT 1: Cache-Aside Pattern Demo")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println()

	// Initialize database
	db := NewDatabase()

	// Initialize Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Use default database
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v\n", err)
	}
	fmt.Println("✅ Connected to Redis successfully!")
	fmt.Println()

	// Clear any existing cache data for clean demo
	redisClient.FlushDB(ctx)

	// Create services
	serviceWithoutCache := NewUserServiceWithoutCache(db)
	serviceWithCache := NewUserServiceWithCache(db, redisClient)

	// ============================================
	// DEMO 1: WITHOUT CACHE
	// ============================================
	fmt.Println("📊 DEMO 1: WITHOUT CACHE (Every request hits database)")
	fmt.Println("-" + string(make([]byte, 50)))

	for i := 1; i <= 3; i++ {
		start := time.Now()
		user, err := serviceWithoutCache.GetUser("user:1")
		duration := time.Since(start)

		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Request %d: Got user '%s' in %v\n", i, user, duration)
	}

	fmt.Println()

	// ============================================
	// DEMO 2: WITH CACHE
	// ============================================
	fmt.Println("📊 DEMO 2: WITH CACHE (First request slow, subsequent fast)")
	fmt.Println("-" + string(make([]byte, 50)))

	for i := 1; i <= 3; i++ {
		start := time.Now()
		user, err := serviceWithCache.GetUser("user:2")
		duration := time.Since(start)

		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Request %d: Got user '%s' in %v\n", i, user, duration)
	}

	fmt.Println()

	// ============================================
	// DEMO 3: CACHE WITH DIFFERENT USERS
	// ============================================
	fmt.Println("📊 DEMO 3: MULTIPLE USERS (Each user's first request is slow)")
	fmt.Println("-" + string(make([]byte, 50)))

	userIDs := []string{"user:3", "user:4", "user:5"}

	// First round: All cache misses
	fmt.Println("\n🔄 First Round (expect cache misses):")
	for _, userID := range userIDs {
		start := time.Now()
		user, _ := serviceWithCache.GetUser(userID)
		duration := time.Since(start)
		fmt.Printf("  %s: '%s' in %v\n", userID, user, duration)
	}

	// Second round: All cache hits
	fmt.Println("\n🔄 Second Round (expect cache hits):")
	for _, userID := range userIDs {
		start := time.Now()
		user, _ := serviceWithCache.GetUser(userID)
		duration := time.Since(start)
		fmt.Printf("  %s: '%s' in %v\n", userID, user, duration)
	}

	// ============================================
	// PERFORMANCE SUMMARY
	// ============================================
	fmt.Println()
	fmt.Println("📈 PERFORMANCE SUMMARY")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println("Without Cache: ~2000ms per request (always slow)")
	fmt.Println("With Cache:")
	fmt.Println("  - First request (miss): ~2000ms (slow)")
	fmt.Println("  - Subsequent requests (hit): ~1-10ms (200x faster!)")
	fmt.Println()
	fmt.Println("💡 Cache Hit Rate: Higher is better!")
	fmt.Println("💡 TTL: Data expires after 5 minutes")

	// Cleanup
	redisClient.Close()
	fmt.Println()
	fmt.Println("✅ Demo completed!")
}

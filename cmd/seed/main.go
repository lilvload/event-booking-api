// Command seed fills an empty development database with sample users,
// events and registrations, so the frontend has data to show right away.
//
// Run it from the project root:
//
//	go run ./cmd/seed
//
// It uses the same DB_PATH as the server (from .env, default api.db) and
// refuses to run if the database already contains users. Never run it
// against a production database.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lilvload/event-booking-api/db"
	"github.com/lilvload/event-booking-api/models"
)

// All sample users share this password.
const seedPassword = "password123"

var seedUsers = []string{
	"alice@example.com",
	"bob@example.com",
	"carol@example.com",
}

type seedEvent struct {
	owner       string
	name        string
	description string
	location    string
	daysFromNow int // days from today; must be positive (events in the past are rejected)
	hour        int // UTC
}

var seedEvents = []seedEvent{
	{"alice@example.com", "Go meetup", "Talks about concurrency and testing in Go.", "Kyiv", 7, 18},
	{"alice@example.com", "Frontend workshop", "Hands-on React session for beginners. Bring a laptop.", "Lviv", 14, 10},
	{"bob@example.com", "Startup pitch night", "Five teams pitch their ideas to investors.", "Kyiv", 3, 19},
	{"bob@example.com", "Hackathon: Smart City", "A 48-hour hackathon on urban mobility and open city data. Teams of up to four people; mentors from local companies help throughout the weekend, and the best three projects get a demo slot at the city council.", "Kharkiv", 30, 9},
	{"carol@example.com", "Design systems talk", "How to build and maintain a design system.", "Odesa", 21, 17},
	{"carol@example.com", "Book club: Clean Code", "Monthly discussion of chapters 1-5.", "Online", 10, 19},
	{"bob@example.com", "Board games evening", "Bring your favourite game.", "Dnipro", 3, 18},
	{"alice@example.com", "Database design basics", "Normalization, indexes and transactions explained with SQLite.", "Online", 45, 18},
}

// Each registration points to an event by its index in seedEvents.
var seedRegistrations = []struct {
	user  string
	event int
}{
	{"bob@example.com", 0},
	{"carol@example.com", 0},
	{"alice@example.com", 2},
	{"carol@example.com", 3},
	{"alice@example.com", 4},
	{"bob@example.com", 5},
	{"alice@example.com", 6},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "api.db"
	}

	db.InitDB()

	var userCount int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		log.Fatalf("Could not read the users table: %v", err)
	}
	if userCount > 0 {
		log.Fatalf("Database %s already has %d user(s). Seed only fills an empty database: delete the file (or set DB_PATH to a new one) and run again.", dbPath, userCount)
	}

	fmt.Printf("Seeding %s...\n", dbPath)

	userIDs := make(map[string]int64)
	for _, email := range seedUsers {
		user := models.User{Email: email, Password: seedPassword}
		if err := user.Save(); err != nil {
			log.Fatalf("Could not create user %s: %v", email, err)
		}
		userIDs[email] = user.ID
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	events := make([]models.Event, len(seedEvents))
	for i, se := range seedEvents {
		events[i] = models.Event{
			Name:        se.name,
			Description: se.description,
			Location:    se.location,
			DateTime:    today.AddDate(0, 0, se.daysFromNow).Add(time.Duration(se.hour) * time.Hour),
			UserID:      userIDs[se.owner],
		}
		if err := events[i].Save(); err != nil {
			log.Fatalf("Could not create event %q: %v", se.name, err)
		}
	}

	for _, r := range seedRegistrations {
		if err := events[r.event].Register(userIDs[r.user]); err != nil {
			log.Fatalf("Could not register %s for %q: %v", r.user, events[r.event].Name, err)
		}
	}

	fmt.Printf("Created %d users, %d events and %d registrations.\n", len(seedUsers), len(seedEvents), len(seedRegistrations))
	fmt.Println("Log in with any of these accounts, password:", seedPassword)
	for _, email := range seedUsers {
		fmt.Println("  " + email)
	}
}

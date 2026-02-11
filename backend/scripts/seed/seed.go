package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"inside-athletics/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SeedCollege represents the JSON structure for college data
type SeedCollege struct {
	Name         string `json:"name"`
	State        string `json:"state"`
	City         string `json:"city"`
	Website      string `json:"website"`
	AcademicRank *int16 `json:"academic_rank,omitempty"`
	DivisionRank int8   `json:"division_rank"`
	Logo         string `json:"logo,omitempty"`
}

// SeedSport represents the JSON structure for sport data
type SeedSport struct {
	Name       string `json:"name"`
	Popularity *int32 `json:"popularity,omitempty"`
}

func main() {
	// Parse command line flags
	collegesFile := flag.String("colleges", "scripts/seed/data/colleges.json", "Path to colleges JSON file")
	sportsFile := flag.String("sports", "scripts/seed/data/sports.json", "Path to sports JSON file")
	dbURL := flag.String("db", os.Getenv("DEV_DB_CONNECTION_STRING"), "Database connection string")
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("Database connection string is required. Set DEV_DB_CONNECTION_STRING or use -db flag")
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(*dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Seed colleges
	if _, err := os.Stat(*collegesFile); err == nil {
		if err := seedColleges(db, *collegesFile); err != nil {
			log.Fatalf("Failed to seed colleges: %v", err)
		}
		log.Println("✓ Colleges seeded successfully")
	} else {
		log.Printf("Skipping colleges - file not found: %s", *collegesFile)
	}

	// Seed sports
	if _, err := os.Stat(*sportsFile); err == nil {
		if err := seedSports(db, *sportsFile); err != nil {
			log.Fatalf("Failed to seed sports: %v", err)
		}
		log.Println("✓ Sports seeded successfully")
	} else {
		log.Printf("Skipping sports - file not found: %s", *sportsFile)
	}

	log.Println("Seeding completed successfully!")
}

func seedColleges(db *gorm.DB, filename string) error {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read colleges file: %w", err)
	}

	// Parse JSON
	var seedColleges []SeedCollege
	if err := json.Unmarshal(data, &seedColleges); err != nil {
		return fmt.Errorf("failed to parse colleges JSON: %w", err)
	}

	// Insert into database
	for _, sc := range seedColleges {
		college := models.College{
			Name:         sc.Name,
			State:        sc.State,
			City:         sc.City,
			Website:      sc.Website,
			AcademicRank: sc.AcademicRank,
			DivisionRank: sc.DivisionRank,
			Logo:         sc.Logo,
		}

		// Use FirstOrCreate to avoid duplicates
		result := db.Where("name = ? AND state = ?", college.Name, college.State).FirstOrCreate(&college)
		if result.Error != nil {
			return fmt.Errorf("failed to create college %s: %w", college.Name, result.Error)
		}
	}

	log.Printf("Processed %d colleges", len(seedColleges))
	return nil
}

func seedSports(db *gorm.DB, filename string) error {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read sports file: %w", err)
	}

	// Parse JSON
	var seedSports []SeedSport
	if err := json.Unmarshal(data, &seedSports); err != nil {
		return fmt.Errorf("failed to parse sports JSON: %w", err)
	}

	// Insert into database
	for _, ss := range seedSports {
		sport := models.Sport{
			Name:       ss.Name,
			Popularity: ss.Popularity,
		}

		// Use FirstOrCreate to avoid duplicates
		result := db.Where("name = ?", sport.Name).FirstOrCreate(&sport)
		if result.Error != nil {
			return fmt.Errorf("failed to create sport %s: %w", sport.Name, result.Error)
		}
	}

	log.Printf("Processed %d sports", len(seedSports))
	return nil
}

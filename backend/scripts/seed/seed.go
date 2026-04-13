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
	"gorm.io/gorm/logger"
)

type SeedCollege struct {
	Name         string `json:"name"`
	State        string `json:"state"`
	City         string `json:"city"`
	Website      string `json:"website"`
	AcademicRank *int16 `json:"academic_rank,omitempty"`
	DivisionRank uint   `json:"division_rank"`
	Logo         string `json:"logo,omitempty"`
}

type SeedSport struct {
	Name       string `json:"name"`
	Popularity *int32 `json:"popularity,omitempty"`
}

func main() {
	collegesFile := flag.String("colleges", "scripts/seed/data/colleges.json", "Path to colleges JSON file")
	sportsFile := flag.String("sports", "scripts/seed/data/sports.json", "Path to sports JSON file")
	dbURL := flag.String("db", os.Getenv("DEV_DB_CONNECTION_STRING"), "Database connection string (defaults to DEV_DB_CONNECTION_STRING)")
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("ERROR: Database connection string is required. Set DEV_DB_CONNECTION_STRING or use -db flag")
	}

	log.Printf("Connecting to database: %s", maskPassword(*dbURL))

	db, err := gorm.Open(postgres.Open(*dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("ERROR: Failed to get underlying DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("ERROR: Failed to ping database: %v", err)
	}

	var currentDB, currentSchema string
	db.Raw("SELECT current_database()").Scan(&currentDB)
	db.Raw("SELECT current_schema()").Scan(&currentSchema)
	log.Printf("✓ Connected — database: %s, schema: %s", currentDB, currentSchema)

	if _, err := os.Stat(*collegesFile); err == nil {
		if err := seedColleges(db, *collegesFile); err != nil {
			log.Fatalf("ERROR: Failed to seed colleges: %v", err)
		}
		log.Println("✓ Colleges seeded successfully")
	} else {
		log.Printf("SKIP: Colleges file not found: %s", *collegesFile)
	}

	if _, err := os.Stat(*sportsFile); err == nil {
		if err := seedSports(db, *sportsFile); err != nil {
			log.Fatalf("ERROR: Failed to seed sports: %v", err)
		}
		log.Println("✓ Sports seeded successfully")
	} else {
		log.Printf("SKIP: Sports file not found: %s", *sportsFile)
	}

	log.Println("✓ Seeding completed successfully!")
}

func seedColleges(db *gorm.DB, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read colleges file: %w", err)
	}

	var seedColleges []SeedCollege
	if err := json.Unmarshal(data, &seedColleges); err != nil {
		return fmt.Errorf("failed to parse colleges JSON: %w", err)
	}

	log.Printf("Processing %d colleges from %s", len(seedColleges), filename)
	created, updated, failed := 0, 0, 0

	for _, sc := range seedColleges {
		college := models.College{
			Name:         sc.Name,
			State:        sc.State,
			City:         sc.City,
			Website:      sc.Website,
			AcademicRank: sc.AcademicRank,
			DivisionRank: models.Division(sc.DivisionRank),
			Logo:         sc.Logo,
		}

		result := db.Where("name = ? AND state = ?", college.Name, college.State).FirstOrCreate(&college)
		if result.Error != nil {
			log.Printf("  ERROR: FirstOrCreate failed for %s: %v", sc.Name, result.Error)
			failed++
			continue
		}

		if result.RowsAffected == 1 {
			log.Printf("  ✓ Created: %s (id=%s)", sc.Name, college.ID)
			created++
		} else {
			updateResult := db.Model(&college).Updates(map[string]interface{}{
				"logo":    sc.Logo,
				"city":    sc.City,
				"website": sc.Website,
			})
			if updateResult.Error != nil {
				log.Printf("  ERROR: Update failed for %s: %v", sc.Name, updateResult.Error)
				failed++
				continue
			}
			log.Printf("  → Updated: %s (id=%s)", sc.Name, college.ID)
			updated++
		}
	}

	log.Printf("Summary — created: %d, updated: %d, failed: %d", created, updated, failed)
	return nil
}

func seedSports(db *gorm.DB, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read sports file: %w", err)
	}

	var seedSports []SeedSport
	if err := json.Unmarshal(data, &seedSports); err != nil {
		return fmt.Errorf("failed to parse sports JSON: %w", err)
	}

	log.Printf("Processing %d sports from %s", len(seedSports), filename)

	for _, ss := range seedSports {
		sport := models.Sport{
			Name:       ss.Name,
			Popularity: ss.Popularity,
		}
		result := db.Where("name = ?", sport.Name).FirstOrCreate(&sport)
		if result.Error != nil {
			log.Printf("  ERROR: Failed to seed sport %s: %v", ss.Name, result.Error)
		}
	}

	return nil
}

func maskPassword(url string) string {
	result := []rune(url)
	inPassword := false
	colonCount := 0
	for i, c := range result {
		if c == ':' {
			colonCount++
			if colonCount == 2 {
				inPassword = true
			}
		}
		if c == '@' {
			inPassword = false
		}
		if inPassword && c != ':' {
			result[i] = '*'
		}
	}
	return string(result)
}

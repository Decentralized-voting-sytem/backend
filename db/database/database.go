package database

import (
	"log"
	"sync"

	"github.com/Decentralized-voting-sytem/backend/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DBLock sync.Mutex

func Init() *gorm.DB {
	DBLock.Lock()
	defer DBLock.Unlock()

	dsn := "host=localhost user=postgres password=samyak dbname=voting port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.Vote{}, &models.Voter{}, &models.Candidate{}, &models.Admin{})

	triggerFunction := `
	CREATE OR REPLACE FUNCTION prevent_duplicate_votes()
	RETURNS TRIGGER AS $$
	BEGIN
		IF EXISTS (SELECT 1 FROM votes WHERE voter_id = NEW.voter_id) THEN
			RAISE EXCEPTION 'Voter has already voted';
		END IF;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;`

	if err := db.Exec(triggerFunction).Error; err != nil {
		log.Printf("Error creating trigger function: %v", err)
	}

	trigger := `
	CREATE TRIGGER check_duplicate_vote
	BEFORE INSERT ON votes
	FOR EACH ROW
	EXECUTE FUNCTION prevent_duplicate_votes();`

	// Execute the trigger creation
	if err := db.Exec(trigger).Error; err != nil {
		log.Printf("Error creating trigger: %v", err)
	}

	// Assign the database instance to the global variable
	DB = db
	return db
}

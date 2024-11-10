package database

import (
	"log"
	"sync"
	"fmt"

	"github.com/Decentralized-voting-sytem/backend/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DBLock sync.Mutex

func Init() *gorm.DB {
	DBLock.Lock()
	defer DBLock.Unlock()

	// Database connection string
	dsn := "host=localhost user=postgres password=samyak dbname=voting port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate models
	db.AutoMigrate(&models.Vote{}, &models.Voter{}, &models.Candidate{}, &models.Admin{})

	// Create or replace the trigger function
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

	// Drop the existing trigger if it exists
	dropTrigger := `
	DO $$ 
	BEGIN
		IF EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'check_duplicate_vote') THEN
			EXECUTE 'DROP TRIGGER check_duplicate_vote ON votes';
		END IF;
	END $$;`

	if err := db.Exec(dropTrigger).Error; err != nil {
		log.Printf("Error dropping existing trigger: %v", err)
	}

	// Create the new trigger
	trigger := `
	CREATE TRIGGER check_duplicate_vote
	BEFORE INSERT ON votes
	FOR EACH ROW
	EXECUTE FUNCTION prevent_duplicate_votes();`

	if err := db.Exec(trigger).Error; err != nil {
		log.Printf("Error creating trigger: %v", err)
	}

	// Define the stored procedure for granting privileges
	grantPrivilegesProcedure := `
	CREATE OR REPLACE PROCEDURE grant_privileges_to_user(username TEXT) AS $$
	BEGIN
		EXECUTE format('GRANT INSERT, UPDATE, DELETE ON TABLE votes, voters, candidates TO %I;', username);
	END;
	$$ LANGUAGE plpgsql;`

	// Create or replace the stored procedure
	if err := db.Exec(grantPrivilegesProcedure).Error; err != nil {
		log.Printf("Error creating grant_privileges_to_user procedure: %v", err)
	}

	// Assign the database instance to the global variable
	DB = db
	return db
}

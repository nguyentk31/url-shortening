package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/nguyentk31/url-shortening/internal/config"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(config config.Database) (*Postgres, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Try connecting to the database
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(config.Timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("failed to connect to database: timeout exceeded")

		case <-ticker.C:
			err := db.Ping()
			if err == nil {
				log.Println("connected to PostgreSQL database")
				return &Postgres{DB: db}, nil
			}
			log.Printf("waiting for PostgreSQL database connection: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (p *Postgres) Close() error {
	err := p.DB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	log.Println("PostgreSQL database connection closed")
	return nil
}

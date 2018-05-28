package env

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	// Select postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ConnectPostgres connects to a Postgres database using the given settings
// and returns a *gorm.DB handle.
func ConnectPostgres(s DatabaseSettings) (*gorm.DB, error) {

	if s.Host == "" {
		s.Host = "127.0.0.1"
	}
	if s.Name == "" {
		return nil, errors.New("Must specify database name")
	}
	if s.User == "" {
		return nil, errors.New("Must specify database user")
	}
	if s.Password == "" {
		return nil, errors.New("Must specify database password")
	}

	sslMode := "require"
	if s.DisableSSL {
		sslMode = "disable"
	}

	args := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		s.User, s.Password, s.Name, sslMode)

	db, err := gorm.Open("postgres", args)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %s", err.Error())
	}

	db.LogMode(false)
	return db, nil
}

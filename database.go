package base

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	// Select postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// SSLMode defines SSL settings used to connect to Postgres
type SSLMode string

const (
	// SSLModeDisabled disables SSL
	SSLModeDisabled SSLMode = "disable"

	// SSLModeRequired makes SSL required
	SSLModeRequired SSLMode = "require"

	// SSLModeVerifyCA enables SSL with server and client certificates
	SSLModeVerifyCA SSLMode = "verify-ca"
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

	args := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s",
		s.User, s.Password, s.Name, s.SSLMode, s.Host)

	if s.SSLRootCert != "" {
		args += fmt.Sprintf(" sslrootcert=%s", s.SSLRootCert)
	}
	if s.SSLCert != "" {
		args += fmt.Sprintf(" sslcert=%s", s.SSLCert)
	}
	if s.SSLKey != "" {
		args += fmt.Sprintf(" sslkey=%s", s.SSLKey)
	}

	db, err := gorm.Open("postgres", args)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %s", err.Error())
	}

	db.LogMode(false)
	return db, nil
}

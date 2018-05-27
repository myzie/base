package env

import (
	"github.com/namsral/flag"
)

// Settings for the application
type Settings struct {
	Database    DatabaseSettings
	ObjectStore ObjectStoreSettings
	Auth        AuthSettings
	HTTP        HTTPSettings
}

// HTTPSettings configures the application HTTP interface.
type HTTPSettings struct {
	ListenAddress string
}

// AuthSettings configures application authentication.
type AuthSettings struct {
	PublicKey  string
	PrivateKey string
}

// ObjectStoreSettings configures the application object storage endpoint.
type ObjectStoreSettings struct {
	URL        string
	Region     string
	Bucket     string
	DisableSSL bool
}

// DatabaseSettings specifies location and authentication information
// needed to connect to the application database.
type DatabaseSettings struct {
	Host       string
	Port       int
	User       string
	Password   string
	Name       string
	DisableSSL bool
}

// GetSettings returns application configuration derived from command line
// options and environment variables.
func GetSettings() Settings {

	var s Settings

	flag.StringVar(&s.Auth.PublicKey, "auth-public-key", "", "Auth public key")
	flag.StringVar(&s.Auth.PrivateKey, "auth-private-key", "", "Auth private key")

	flag.StringVar(&s.ObjectStore.Bucket, "storage-bucket", "", "Storage bucket")
	flag.StringVar(&s.ObjectStore.URL, "storage-url", "s3.amazonaws.com", "Storage URL")
	flag.StringVar(&s.ObjectStore.Region, "storage-region", "us-east-1", "Storage region")
	flag.BoolVar(&s.ObjectStore.DisableSSL, "storage-disable-ssl", false, "Storage disable SSL")

	flag.StringVar(&s.Database.Name, "db-name", "", "DB name")
	flag.StringVar(&s.Database.User, "db-user", "", "DB user")
	flag.StringVar(&s.Database.Password, "db-password", "", "DB password")
	flag.StringVar(&s.Database.Host, "db-host", "localhost", "DB host address")
	flag.IntVar(&s.Database.Port, "db-port", 5432, "DB port")
	flag.BoolVar(&s.Database.DisableSSL, "db-disable-ssl", false, "DB disable SSL")

	flag.StringVar(&s.HTTP.ListenAddress, "http", "127.0.0.1:8080", "HTTP listen address")

	flag.Parse()
	return s
}

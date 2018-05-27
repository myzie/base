package env

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

// Env contains resource handles needed for this web service.
// This struct is made available to HTTP handlers.
type Env struct {
	Settings       Settings
	AuthPublicKey  *rsa.PublicKey
	AuthPrivateKey *rsa.PrivateKey
	DB             *gorm.DB
	ObjectStore    *minio.Client
	Echo           *echo.Echo
}

// Destroy resource handles
func (env *Env) Destroy() {
	if env.DB != nil {
		env.DB.Close()
		env.DB = nil
	}
	// Nothing else needs cleanup...
}

// Run the HTTP server
func (env *Env) Run() error {
	return env.Echo.Start(env.Settings.HTTP.ListenAddress)
}

// New environment using settings provided via command line options and
// environment variables. You can also implement your own version of this
// if you need to since Env fields are all public.
func New() (*Env, error) {

	var err error
	s := GetSettings()

	env := &Env{
		Settings: s,
		Echo:     echo.New(),
	}

	if s.ObjectStore.Bucket != "" {
		accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		objStoreSSL := true
		if s.ObjectStore.DisableSSL {
			objStoreSSL = false
		}
		env.ObjectStore, err = minio.New(s.ObjectStore.URL, accessKey, secretKey, objStoreSSL)
		if err != nil {
			return nil, fmt.Errorf("ObjectStore error: %s", err.Error())
		}
		err = env.ObjectStore.MakeBucket(s.ObjectStore.Bucket, s.ObjectStore.Region)
		if err != nil {
			return nil, fmt.Errorf("ObjectStore bucket error: %s", err.Error())
		}
	}

	env.DB, err = ConnectPostgres(s.Database)
	if err != nil {
		return nil, err
	}

	if s.Auth.PrivateKey != "" {
		env.AuthPrivateKey, err = LoadRSAPrivateKey(s.Auth.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("Private key error: %s", err.Error())
		}
	}

	if s.Auth.PublicKey != "" {
		env.AuthPublicKey, err = LoadRSAPublicKey(s.Auth.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("Public key error: %s", err.Error())
		}
	}

	if s.HTTP.ListenAddress != "" {
		env.Echo, err = SetupEcho()
		if err != nil {
			return nil, fmt.Errorf("Echo error: %s", err.Error())
		}
	}

	return env, nil
}

// Must is a helper that wrap New and exits the process on failure
func Must() *Env {
	e, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return e
}

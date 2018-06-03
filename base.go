package base

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

// Base contains a core set of dependencies needed for web services.
// This struct can be made available to HTTP handlers.
type Base struct {
	Settings       Settings
	AuthPublicKey  *rsa.PublicKey
	AuthPrivateKey *rsa.PrivateKey
	DB             *gorm.DB
	ObjectStore    *minio.Client
	Echo           *echo.Echo
}

// Destroy resource handles
func (base *Base) Destroy() {
	if base.DB != nil {
		base.DB.Close()
		base.DB = nil
	}
	// Nothing else needs cleanup...
}

// Run the HTTP server
func (base *Base) Run() error {
	log.WithField("address", base.Settings.HTTP.ListenAddress).Infof("Listening")
	return base.Echo.Start(base.Settings.HTTP.ListenAddress)
}

// JWTMiddleware returns JWT authentication middleware
func (base *Base) JWTMiddleware() echo.MiddlewareFunc {
	if base.AuthPublicKey == nil {
		panic("JWTMiddleware unavailable: no public key set")
	}
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:        &JWTClaims{},
		SigningKey:    base.AuthPublicKey,
		SigningMethod: "RS512",
		ContextKey:    "user",
	})
}

// New environment using settings provided via command line options and
// environment variables. You can also implement your own version of this
// if you need to since Env fields are all public.
func New() (*Base, error) {

	var err error
	s := GetSettings()

	base := &Base{
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
		base.ObjectStore, err = minio.New(s.ObjectStore.URL, accessKey, secretKey, objStoreSSL)
		if err != nil {
			return nil, fmt.Errorf("ObjectStore error: %s", err.Error())
		}
		err = base.ObjectStore.MakeBucket(s.ObjectStore.Bucket, s.ObjectStore.Region)
		if err != nil {
			return nil, fmt.Errorf("ObjectStore bucket error: %s", err.Error())
		}
	}

	if s.Database.Host != "" && s.Database.Port != 0 {
		base.DB, err = ConnectPostgres(s.Database)
		if err != nil {
			return nil, err
		}
	}

	if s.Auth.PrivateKey != "" {
		base.AuthPrivateKey, err = LoadRSAPrivateKey(s.Auth.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("Private key error: %s", err.Error())
		}
	}

	if s.Auth.PublicKey != "" {
		base.AuthPublicKey, err = LoadRSAPublicKey(s.Auth.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("Public key error: %s", err.Error())
		}
	}

	if s.HTTP.ListenAddress != "" {
		base.Echo, err = SetupEcho()
		if err != nil {
			return nil, fmt.Errorf("Echo error: %s", err.Error())
		}
	}

	return base, nil
}

// Must is a helper that wrap New and exits the process on failure
func Must() *Base {
	e, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return e
}

# Base

Standardized foundation for Go microservices.

## What is Included

This foundation includes a handful of common dependencies:

 * Flags to configure database, HTTP port, etc.
 * [Gorm](http://gorm.io) database handle
 * RSA keypair for signing and verifying [JWTs](https://jwt.io)
 * [Minio Client](https://github.com/minio/minio-go) to access S3 compatible storage
 * [Echo](https://echo.labstack.com) web server

## Usage

I suggest embeddeding `Base` into a service struct that has your HTTP handlers.

```
type myService struct {
	*base.Base
}
```

You then have easy access to the database handle wherever you have access to this service.

See the [simple example](./examples/simple/main.go).

## Limitations

This reusability works for me but it's not that flexible currently, so I don't expect it to work for you!

You can use the Base struct but customize your own setup. This would allow for example using MySQL instead of
the default Postgres.

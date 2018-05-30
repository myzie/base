package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/myzie/base"
	log "github.com/sirupsen/logrus"
)

type myService struct {
	*base.Base
}

func (s *myService) hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {

	s := &myService{Base: base.Must()}

	s.Echo.GET("/", s.hello)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	oapimiddleware "openapi-gen-auth/middleware"

	"openapi-gen-auth/api"
	"openapi-gen-auth/server"
)

const (
	jwtSecret = "secret"
)

func main() {
	si := server.NewServerImplementation(jwtSecret)

	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authMiddleware, err := oapimiddleware.CreateAuthMiddleware(jwtSecret)
	if err != nil {
		panic(err)
	}

	r.Use(authMiddleware)
	h := api.HandlerFromMux(si, r)

	s := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	fmt.Println("Starting server on port 8080")
	if err = s.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

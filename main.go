package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"net/http"
	"test_project/config"
	"test_project/handler"
)

func main() {

	ctx := context.Background()

	configInstance, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("config creation")
	}

	dbUsername := configInstance.DBUser
	dbPassword := configInstance.DBPass
	dbName := configInstance.DBName
	dbHost := configInstance.DBHost
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUsername, dbName, dbPassword)

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Error().Err(err).Msg("error opening db")
	}
	defer db.Close()

	router := chi.NewRouter()
	router.Route("/registration", func (router chi.Router) {
		router.Post("/", handler.RegPost(ctx, db))
	})
	router.Route("/login", func (router chi.Router) {
		router.Post("/", handler.LoginPost(ctx, db))
	})

	log.Log().Msg("server serve on port " + configInstance.ServerPort)
	err = http.ListenAndServe(configInstance.ServerPort, router)
	if err != nil {
		log.Error().Err(err).Msg("error listening server")
	}
}
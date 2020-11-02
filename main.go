package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"net/http"
	"test_project/config"
	"test_project/handler"
)

func main() {

	configInstance, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("config creation")
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Post("/registration", handler.RegPost())
	router.Post("/login", handler.LoginPost())

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

	log.Log().Msg("server serve on port " + configInstance.ServerPort)
	http.ListenAndServe(configInstance.ServerPort, router)
}

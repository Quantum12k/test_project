package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"test_project/helper"
	"test_project/models"
)

type response struct {
	Id string `json:"id"`
}

func RegPost(ctx context.Context, db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		dbConn, err := db.Conn(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error connecting to DB")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer dbConn.Close()

		var user models.User

		//request decode
		err = helper.DecodeJSONBody(w, r, &user)
		if err != nil {
			var mr *helper.MalformedRequest
			if errors.As(err, &mr) {
				http.Error(w, mr.GetJSONString(), mr.Status)
			} else {
				log.Error().Err(err).Msg("error decoding JSON at user registration")
				http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		//email regex check
		regex, err := regexp.Compile("[A-Z0-9a-z._-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}")
		if err != nil {
			log.Error().Err(err).Msg("error compile regex for email")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if !regex.MatchString(user.Email) {
			log.Error().Err(err).Msg("invalid e-mail format was sent at user registration")
			http.Error(w, helper.FormJSONError("invalid e-mail format was sent", http.StatusOK), http.StatusOK)
			return
		}

		//checking for email in the database
		queryStr := "select users.id from users where users.email = '" + user.Email + "'"
		rows, err := dbConn.QueryContext(ctx, queryStr)
		if err != nil {
			log.Error().Err(err).Msg("error sending query to DB to check user")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		if rows.Next() {
			//there is a user with the same email, can't register
			log.Error().Msg("error trying to register a user with existing email")
			http.Error(w, helper.FormJSONError("user with the same email already exists", http.StatusOK), http.StatusOK)
			return
		}

		if user.Password == "" {
			log.Error().Msg("error trying to register a user with empty password")
			http.Error(w, helper.FormJSONError("password can't be empty", http.StatusOK), http.StatusOK)
			return
		}

		//generating user id
		userID := response{Id: uuid.New().String()}

		queryStr = "insert into users(id, email, pass) values('" + userID.Id + "','" + user.Email + "','" + user.Password + "')"
		_, err = dbConn.ExecContext(ctx, queryStr)
		if err != nil {
			log.Error().Err(err).Msg("error executing query to register user in DB")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		idJSON, err := json.Marshal(userID)

		if err != nil {
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(idJSON)
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
		}
	}
}
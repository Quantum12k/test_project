package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"test_project/helper"
	"test_project/models"
)

func LoginPost(ctx context.Context, db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		dbConn, err := db.Conn(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error connecting to DB")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer dbConn.Close()

		var user models.RequestUserInfo

		//request decode
		err = helper.DecodeJSONBody(w, r, &user)
		if err != nil {
			var mr *helper.MalformedRequest
			if errors.As(err, &mr) {
				http.Error(w, mr.GetJSONString(), mr.Status)
			} else {
				log.Error().Err(err).Msg("error decoding JSON at user login")
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
			log.Error().Err(err).Msg("invalid e-mail format was sent at user login")
			http.Error(w, helper.FormJSONError("invalid e-mail format was sent", http.StatusOK), http.StatusOK)
			return
		}

		//checking for email in the database
		queryStr := "select * from users where users.email = '" + user.Email + "'"
		rows, err := dbConn.QueryContext(ctx, queryStr)
		if err != nil {
			log.Error().Err(err).Msg("error sending query to DB to check user")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			log.Error().Msg("error trying to find not existing user")
			http.Error(w, helper.FormJSONError("user not found", http.StatusOK), http.StatusOK)
			return
		}

		var dbUser models.User
		err = rows.Scan(&dbUser.Id, &dbUser.Email, &dbUser.Password)
		if err != nil {
			log.Error().Err(err).Msg("error scanning rows from DB query request")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if user.Password != dbUser.Password {
			log.Error().Msg("error trying to login with invalid password")
			http.Error(w, helper.FormJSONError("wrong password", http.StatusOK), http.StatusOK)
			return
		}

		idJSON, err := json.Marshal(models.ResponseUserID{Id: dbUser.Id})

		if err != nil {
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(idJSON)
		if err != nil {
			log.Error().Err(err).Msg("error writing response")
			http.Error(w, helper.FormJSONErrorByStatus(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
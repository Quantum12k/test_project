package handler

import (
	"context"
	"database/sql"
	"net/http"
)

func LoginPost(ctx context.Context, db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {


	}
}
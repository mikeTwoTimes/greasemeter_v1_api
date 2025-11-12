package utility

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

var constraintMessages = map[string]string{
	"users_email_key":                          "Email already in use",
	"users_name_key":                           "Username taken",
	"reviews_user_id_place_id_key":             "Place already reviewed",
	"bookmarks_user_id_place_id_key":           "Place already bookmarked",
	"recommendations_user_id_name_address_key": "Place already recommended",
	"reports_user_id_place_id_key":             "Place already reported",
}

func MapError(err error) (int, gin.H) {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}
	} else if pgErr.Code == "P0001" {
		return http.StatusConflict, gin.H{"error": pgErr.Message}
	} else if pgErr.Code == "23505" {
		if msg, ok := constraintMessages[pgErr.ConstraintName]; ok {
			return http.StatusConflict, gin.H{"error": msg}
		}
	}

	return http.StatusInternalServerError, gin.H{"error": pgErr.Message}
}

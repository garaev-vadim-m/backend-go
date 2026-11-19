package controller

import "database/sql"

type Handler struct {
	db           *sql.DB
	cookieSecure bool
}

func NewHandler(db *sql.DB, cookieSecure bool) *Handler {
	return &Handler{db: db, cookieSecure: cookieSecure}
}

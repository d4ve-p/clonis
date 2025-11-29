package ui

import "github.com/d4ve-p/clonis/internal/database"

type Handler struct {
	Store *database.Store
}

func New(store *database.Store) *Handler {
	return &Handler {
		Store: store,
	}
}

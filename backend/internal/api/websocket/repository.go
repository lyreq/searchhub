package websocket

import (
	"lytemp/pkg/database"
)

type repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) Repository {
	return &repository{db}
}

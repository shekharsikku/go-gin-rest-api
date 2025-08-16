package database

import "database/sql"

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id          int    `json:"id"`
	Owner       int    `json:"owner" binding:"required"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Location    string `json:"location" binding:"required,min=3"`
	Date        string `json:"date" binding:"required, datetime=2005-01-02"`
}

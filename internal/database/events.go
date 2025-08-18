package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/shekharsikku/go-gin-rest-api/internal/utils"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id          int    `json:"id"`
	Owner       int    `json:"owner" binding:"required"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Location    string `json:"location" binding:"required,min=3"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
}

func (em *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "INSERT INTO events (id, owner, name, description, location, date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	id := utils.GenerateUniqueID()

	return em.DB.QueryRowContext(ctx, query, id, event.Owner, event.Name, event.Description, event.Date, event.Location).Scan(&event.Id)
}

func (em *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "SELECT * FROM events"

	rows, err := em.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event

		err := rows.Scan(&event.Id, &event.Owner, &event.Name, &event.Description, &event.Location, &event.Date)

		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (em *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "SELECT * FROM events WHERE id = $1"

	var event Event

	err := em.DB.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.Owner, &event.Name, &event.Description, &event.Location, &event.Date)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

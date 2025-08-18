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
	Id          int       `json:"id"`
	Owner       int       `json:"owner" binding:"required"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Location    string    `json:"location" binding:"required,min=3"`
	DateTime    time.Time `json:"datetime" binding:"required" time_format:"2006-01-02 15:04:05"`
}

func (em *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "INSERT INTO events (id, owner, name, description, location, date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	id := utils.GenerateUniqueID()

	return em.DB.QueryRowContext(ctx, query, id, event.Owner, event.Name, event.Description, event.Location, event.DateTime).Scan(&event.Id)
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

		err := rows.Scan(&event.Id, &event.Owner, &event.Name, &event.Description, &event.Location, &event.DateTime)

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

	err := em.DB.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.Owner, &event.Name, &event.Description, &event.Location, &event.DateTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func (em *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "UPDATE events SET name = $1, description = $2, location = $3, date = $4 WHERE id = $5"

	_, err := em.DB.ExecContext(ctx, query, event.Name, event.Description, event.Location, event.DateTime, event.Id)

	if err != nil {
		return err
	}

	return nil
}

func (em *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "DELETE FROM events WHERE id = $1"

	_, err := em.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

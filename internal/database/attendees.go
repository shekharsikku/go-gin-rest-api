package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"user_id"`
	EventId int `json:"event_id"`
}

func (em *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "INSERT INTO attendees (event_id, user_id) VALUES ($1, $2) RETURNING id"

	err := em.DB.QueryRowContext(ctx, query, attendee.EventId, attendee.UserId).Scan(&attendee.Id)

	if err != nil {
		return nil, err
	}

	return attendee, nil
}

func (em *AttendeeModel) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "DELETE FROM attendees WHERE user_id = $1 AND event_id = $2"

	_, err := em.DB.ExecContext(ctx, query, userId, eventId)

	if err != nil {
		return err
	}

	return nil
}

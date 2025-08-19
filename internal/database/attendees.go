package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/shekharsikku/go-gin-rest-api/internal/utils"
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

	query := "INSERT INTO attendees (id, event_id, user_id) VALUES ($1, $2, $3) RETURNING id"

	id := utils.GenerateUniqueID()

	err := em.DB.QueryRowContext(ctx, query, id, attendee.EventId, attendee.UserId).Scan(&attendee.Id)

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

func (em *AttendeeModel) GetByEventAndAttendee(eventId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "SELECT * FROM attendees where event_id = $1 AND user_id = $2"

	var attendee Attendee

	err := em.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (em *AttendeeModel) GetAttendeesByEvent(eventId int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	 SELECT u.id, u.name, u.email
	 FROM users u
	 JOIN attendees a ON u.id = a.user_id
	 where a.event_id = $1
	`

	rows, err := em.DB.QueryContext(ctx, query, eventId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Name, &user.Email)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (em *AttendeeModel) GetEventsByAttendee(attendeeId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		SELECT e.id, e.owner, e.name, e.description, e.location, e.datetime
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = $1
	`

	rows, err := em.DB.QueryContext(ctx, query, attendeeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var event Event

		err := rows.Scan(&event.Id, &event.Owner, &event.Name, &event.Description, &event.Location, &event.DateTime)

		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}

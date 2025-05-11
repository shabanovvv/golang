package sqlstorage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(ctx context.Context, dsn string) *Storage {
	storage := &Storage{}
	err := storage.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	return storage
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	s.db = db

	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	query := `
		SELECT id, title, date_from, date_to, description, user_id, notification_time
		FROM event 
		WHERE id = :id`

	args := map[string]interface{}{
		"id": id,
	}
	row, err := s.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	result, err := convertRowsToEvents(row)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	query := `
		INSERT INTO event (
			   title, 
			   date_from, 
			   date_to, 
			   description, 
			   user_id, 
			   notification_time
		) 
		VALUES (
				:title, 
				:date_from, 
				:date_to, 
				:description, 
				:user_id, 
				:notification_time
		)`
	_, err := s.db.NamedExecContext(ctx, query, event)

	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := `
		UPDATE event
		SET title = :title, 
			date_from = :date_from, 
			date_to = :date_to, 
			description = :description, 
			user_id = :user_id, 
			notification_time = :notification_time
		WHERE id = :id`
	_, err := s.db.NamedExecContext(ctx, query, event)

	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM event WHERE id = :id`
	_, err := s.db.NamedExecContext(ctx, query, id)

	return err
}

func (s *Storage) GetEventsForDay(ctx context.Context, dateFrom time.Time) ([]storage.Event, error) {
	return s.GetEventsByDateRange(ctx, dateFrom, dateFrom.AddDate(0, 0, 1))
}

func (s *Storage) GetEventsForWeek(ctx context.Context, dateFrom time.Time) ([]storage.Event, error) {
	return s.GetEventsByDateRange(ctx, dateFrom, dateFrom.AddDate(0, 0, 7))
}

func (s *Storage) GetEventsForMonth(ctx context.Context, dateFrom time.Time) ([]storage.Event, error) {
	return s.GetEventsByDateRange(ctx, dateFrom, dateFrom.AddDate(0, 0, 7))
}

func (s *Storage) GetEventsByDateRange(
	ctx context.Context,
	dateFrom time.Time,
	dateTo time.Time,
) ([]storage.Event, error) {
	args := map[string]interface{}{
		"date_from": dateFrom,
		"date_to":   dateTo,
	}
	query := `
		SELECT id, title, date_from, date_to, description, user_id, notification_time
		FROM event 
		WHERE date_from >= :date_from AND date_to >= :date_to
		ORDER BY date_from DESC`
	rows, err := s.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := convertRowsToEvents(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func convertRowsToEvents(rows *sqlx.Rows) ([]storage.Event, error) {
	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

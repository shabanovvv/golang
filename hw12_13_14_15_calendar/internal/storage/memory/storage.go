package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events sync.Map
}

func New(_ context.Context) *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	_, loaded := s.events.LoadOrStore(event.ID, event)
	if loaded {
		return fmt.Errorf("event already exists")
	}

	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event storage.Event) error {
	_, loaded := s.events.Load(event.ID)
	if !loaded {
		return fmt.Errorf("event does not exist")
	}
	s.events.Store(event.ID, event)

	return nil
}

func (s *Storage) GetEvent(_ context.Context, id string) (*storage.Event, error) {
	event, loaded := s.events.Load(id)
	if loaded {
		event, ok := event.(storage.Event)
		if !ok {
			return nil, fmt.Errorf("stored value is not an event")
		}

		return &event, nil
	}

	return nil, fmt.Errorf("event not found")
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	_, loaded := s.events.Load(id)
	if !loaded {
		return fmt.Errorf("event does not exist")
	}
	s.events.Delete(id)

	return nil
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

func (s *Storage) GetEventsByDateRange(_ context.Context, dateFrom, dateTo time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)
	s.events.Range(func(_, value interface{}) bool {
		event, ok := value.(storage.Event)
		if !ok {
			return true
		}
		if (event.DateTo.After(dateFrom) || event.DateTo.Equal(dateFrom)) &&
			(event.DateTo.Before(dateTo) || event.DateTo.Equal(dateTo)) ||
			(event.DateFrom.After(dateFrom) || event.DateFrom.Equal(dateFrom)) &&
				(event.DateFrom.Before(dateTo) || event.DateFrom.Equal(dateTo)) ||
			(event.DateFrom.Before(dateFrom) && event.DateTo.After(dateTo)) {
			events = append(events, event)
		}

		return true
	})

	if len(events) == 0 {
		return nil, fmt.Errorf("no events found")
	}

	return events, nil
}

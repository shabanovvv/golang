package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/storage"
)

var ctx = context.Background()

func TestCreateEvent(t *testing.T) {
	event := storage.Event{
		ID:       "1",
		Title:    "Test Event",
		DateFrom: time.Now(),
		DateTo:   time.Now().Add(1 * time.Hour),
	}

	memorystorage := New(ctx)
	err := memorystorage.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = memorystorage.CreateEvent(ctx, event) // Should return an error because the event already exists
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestUpdateEvent(t *testing.T) {
	event := storage.Event{
		ID:       "1",
		Title:    "Test Event",
		DateFrom: time.Now(),
		DateTo:   time.Now().Add(1 * time.Hour),
	}

	memorystorage := New(ctx)
	err := memorystorage.CreateEvent(ctx, event)
	if err != nil {
		return
	}

	event.Title = "Updated Event"
	err = memorystorage.UpdateEvent(ctx, event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedEvent, err := memorystorage.GetEvent(ctx, event.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedEvent.Title != "Updated Event" {
		t.Fatalf("expected event name to be 'Updated Event', got %s", updatedEvent.Title)
	}

	err = memorystorage.UpdateEvent(ctx, storage.Event{ID: "nonexistent"})
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetEvent(t *testing.T) {
	event := storage.Event{
		ID:       "1",
		Title:    "Test Event",
		DateFrom: time.Now(),
		DateTo:   time.Now().Add(1 * time.Hour),
	}

	memorystorage := New(ctx)
	err := memorystorage.CreateEvent(ctx, event)
	if err != nil {
		return
	}

	retrievedEvent, err := memorystorage.GetEvent(ctx, event.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedEvent.ID != event.ID {
		t.Fatalf("expected event ID to be '%s', got '%s'", event.ID, retrievedEvent.ID)
	}

	_, err = memorystorage.GetEvent(ctx, "nonexistent")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestDeleteEvent(t *testing.T) {
	event := storage.Event{
		ID:       "1",
		Title:    "Test Event",
		DateFrom: time.Now(),
		DateTo:   time.Now().Add(1 * time.Hour),
	}

	memorystorage := New(ctx)
	err := memorystorage.CreateEvent(ctx, event)
	if err != nil {
		return
	}

	err = memorystorage.DeleteEvent(ctx, event.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = memorystorage.DeleteEvent(ctx, event.ID) // Should return an error because the event was deleted
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetEventsForDay(t *testing.T) {
	now := time.Now()

	event1 := storage.Event{
		ID:       "1",
		Title:    "Event 1",
		DateFrom: now.Add(-1 * time.Hour),
		DateTo:   now.Add(1 * time.Hour),
	}

	event2 := storage.Event{
		ID:       "2",
		Title:    "Event 2",
		DateFrom: now.Add(2 * time.Hour),
		DateTo:   now.Add(3 * time.Hour),
	}

	memorystorage := New(ctx)

	err := memorystorage.CreateEvent(ctx, event1)
	if err != nil {
		return
	}
	err = memorystorage.CreateEvent(ctx, event2)
	if err != nil {
		return
	}

	events, err := memorystorage.GetEventsForDay(ctx, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Ожидаем, что оба события попадают в день
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestGetEventsForWeek(t *testing.T) {
	now := time.Now()

	// Изменим event1 так, чтобы он происходил в текущую неделю
	event1 := storage.Event{
		ID:       "1",
		Title:    "Event 1",
		DateFrom: now.Add(-1 * 24 * time.Hour), // Вчера
		DateTo:   now.Add(0 * time.Hour),       // Сегодня
	}

	event2 := storage.Event{
		ID:       "2",
		Title:    "Event 2",
		DateFrom: now.Add(1 * 24 * time.Hour), // Завтра
		DateTo:   now.Add(2 * 24 * time.Hour), // Послезавтра
	}

	memorystorage := New(ctx)

	err := memorystorage.CreateEvent(ctx, event1)
	if err != nil {
		return
	}
	err = memorystorage.CreateEvent(ctx, event2)
	if err != nil {
		return
	}

	events, err := memorystorage.GetEventsForWeek(ctx, now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Ожидаем, что оба события попадают в неделю
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

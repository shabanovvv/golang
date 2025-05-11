package app

import (
	"context"
	"time"

	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	ctx     context.Context
	storage Storage
	logger  Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventsForDay(ctx context.Context, dateFrom time.Time) ([]storage.Event, error)
	GetEventsForWeek(ctx context.Context, dateFrom time.Time) ([]storage.Event, error)
	GetEventsForMonth(ctx context.Context, dateFrom time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		ctx:     context.Background(),
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(
	id string,
	title string,
	dateFrom, dateTo time.Time,
	description string,
	userID int,
	notificationTime time.Duration,
) error {
	event := storage.Event{
		ID:               id,
		Title:            title,
		DateFrom:         dateFrom,
		DateTo:           dateTo,
		Description:      description,
		UserID:           userID,
		NotificationTime: notificationTime,
	}

	return a.storage.CreateEvent(a.ctx, event)
}

func (a *App) GetEvent(id string) (*storage.Event, error) {
	return a.storage.GetEvent(a.ctx, id)
}

func (a *App) UpdateEvent(event storage.Event) error {
	return a.storage.UpdateEvent(a.ctx, event)
}

func (a *App) DeleteEvent(id string) error {
	return a.storage.DeleteEvent(a.ctx, id)
}

func (a *App) GetEventsForDay(dateFrom time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsForDay(a.ctx, dateFrom)
}

func (a *App) GetEventsForWeek(dateFrom time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsForWeek(a.ctx, dateFrom)
}

func (a *App) GetEventsForMonth(dateFrom time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsForMonth(a.ctx, dateFrom)
}

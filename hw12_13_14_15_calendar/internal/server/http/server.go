package internalhttp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/app"
	"github.com/shabanovvv/golang/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger Logger
	app    *app.App
	server *http.Server
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func NewServer(logger Logger, app *app.App) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context, host string, port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.Hello)
	mux.HandleFunc("/events", s.GetEvent)
	mux.HandleFunc("/events/create", s.CreateEvent)
	mux.HandleFunc("/events/update", s.UpdateEvent)
	mux.HandleFunc("/events/delete", s.DeleteEvent)
	mux.HandleFunc("/events/list/day", s.GetEventsForDay)
	mux.HandleFunc("/events/list/week", s.GetEventsForWeek)
	mux.HandleFunc("/events/list/month", s.GetEventsForMonth)

	server := &http.Server{
		Addr:              host + ":" + strconv.Itoa(port),
		Handler:           loggingMiddleware(mux, s.logger),
		ReadHeaderTimeout: 3 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	s.logger.Info("Starting HTTP server at ")
	err := server.ListenAndServe()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) Hello(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello world"))
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.app.CreateEvent(
		event.ID,
		event.Title,
		event.DateFrom,
		event.DateTo,
		event.Description,
		event.UserID,
		event.NotificationTime,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.app.UpdateEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
	}

	event, err := s.app.GetEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
	}

	_, err := s.app.GetEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.app.DeleteEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getEvents(
	w http.ResponseWriter,
	r *http.Request,
	eventFetcher func(time.Time) ([]storage.Event, error),
) {
	dateFromStr := r.URL.Query().Get("dateFrom")
	if dateFromStr == "" {
		http.Error(w, "dateFrom is required", http.StatusBadRequest)
		return
	}

	// Парсим дату
	dateFrom, err := time.Parse("2006-01-02", dateFromStr)
	if err != nil {
		http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Получаем события
	events, err := eventFetcher(dateFrom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок и код ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	s.getEvents(w, r, s.app.GetEventsForDay)
}

func (s *Server) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	s.getEvents(w, r, s.app.GetEventsForWeek)
}

func (s *Server) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	s.getEvents(w, r, s.app.GetEventsForMonth)
}

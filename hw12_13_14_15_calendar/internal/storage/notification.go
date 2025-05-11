package storage

import "time"

type Notification struct {
	ID        string
	Title     string
	Timestamp time.Time
	UserID    int
}

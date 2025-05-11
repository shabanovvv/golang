package storage

import "time"

type Event struct {
	ID               string        `db:"id"`
	Title            string        `db:"title"`
	DateFrom         time.Time     `db:"date_from"`
	DateTo           time.Time     `db:"date_to"`
	Description      string        `db:"description"`
	UserID           int           `db:"user_id"`
	NotificationTime time.Duration `db:"notification_time"`
}

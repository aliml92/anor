package anor

import "time"

type Store struct {
	ID          int32
	Handle      string
	UserID      int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

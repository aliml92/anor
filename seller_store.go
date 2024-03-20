package anor

import "time"

type SellerStore struct {
	ID          int32
	PublicID    string
	SellerID    int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

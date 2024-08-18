package anor

import (
	"context"
	"time"
)

type AddressDefaultType string

const (
	AddressDefaultTypeShipping AddressDefaultType = "Shipping"
	AddressDefaultTypeBilling  AddressDefaultType = "Billing"
)

type Address struct {
	ID            int64
	UserID        int64
	DefaultFor    AddressDefaultType
	Name          string
	AddressLine1  string
	AddressLine2  string
	City          string
	StateProvince string
	PostalCode    string
	Country       string
	Phone         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a Address) Equals(b Address) bool {
	eq := a.ID == b.ID &&
		a.UserID == b.UserID &&
		a.Name == b.Name &&
		a.AddressLine1 == b.AddressLine1 &&
		a.AddressLine2 == b.AddressLine2 &&
		a.City == b.City &&
		a.StateProvince == b.StateProvince &&
		a.PostalCode == b.PostalCode &&
		a.Country == b.Country

	return eq
}

func (a Address) IsZero() bool {
	if a.ID == 0 {
		return true
	}

	return false
}

type AddressCreateParams struct {
	UserID        int64
	DefaultFor    AddressDefaultType
	Name          string
	AddressLine1  string
	AddressLine2  string
	City          string
	StateProvince string
	PostalCode    string
	Country       string
}

type AddressListParams struct {
	UserID   int64
	Page     int
	PageSize int
}

type AddressService interface {
	Create(ctx context.Context, params AddressCreateParams) (Address, error)
	Get(ctx context.Context, id int64) (Address, error)
	GetDefault(ctx context.Context, userID int64, defaultFor AddressDefaultType) (Address, error)
	List(ctx context.Context, params AddressListParams) ([]Address, error)
}

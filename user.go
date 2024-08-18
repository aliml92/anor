package anor

import (
	"context"
	"strings"
	"time"
)

type UserStatus string

const (
	UserStatusBlocked             UserStatus = "Blocked"
	UserStatusRegistrationPending UserStatus = "RegistrationPending"
	UserStatusActive              UserStatus = "Active"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleSeller   UserRole = "seller"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID          int64
	Email       string
	Password    string
	PhoneNumber string
	FullName    string
	Status      UserStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Roles []UserRole
	Cart  *Cart
	//Wishlist     Wishlist
	Orders         []Order
	ActivityCounts UserActivityCounts
}

type UserActivityCounts struct {
	CartItemsCount     int
	WishlistItemsCount int
	ActiveOrdersCount  int
}

func (u User) GetFirstname() string {
	return strings.Fields(u.FullName)[0]
}

type AuthProvider string

const (
	ProviderBuiltin AuthProvider = "builtin"
	ProviderGoogle  AuthProvider = "google"
)

type UserCreateParams struct {
	Name        string
	Email       string
	Password    string
	Provider    AuthProvider
	PhoneNumber string
	Status      UserStatus
}

type UserService interface {
	Create(ctx context.Context, params UserCreateParams) (User, error)
	Get(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	UpdateStatusByEmail(ctx context.Context, status UserStatus, email string) error // switch status and email
	UpdatePassword(ctx context.Context, id int64, password string) error
	GetActivityCounts(ctx context.Context, id int64) (UserActivityCounts, error)
}

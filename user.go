package anor

import (
	"context"
	"strings"
	"time"
)

const (
	userKey = "userKey"
)

// NewContextWithUser NewContext returns a new Context that carries value u.
func NewContextWithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserFromContext FromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(userKey).(*User)
	if !ok {
		return nil
	}
	return u
}

type UserStatus string

const (
	UserStatusBlocked             UserStatus = "Blocked"
	UserStatusRegistrationPending UserStatus = "RegistrationPending"
	UserStatusActive              UserStatus = "Active"
)

type Role string

const (
	RoleCustomer Role = "customer"
	RoleSeller   Role = "seller"
	RoleAdmin    Role = "admin"
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

	Roles []Role
	Cart  *Cart
	//Wishlist    Wishlist
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

type UserService interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUserStatusByEmail(ctx context.Context, status UserStatus, email string) error
	UpdateUserPassword(ctx context.Context, id int64, password string) error

	GetUserActivityCounts(ctx context.Context, id int64) (UserActivityCounts, error)
}

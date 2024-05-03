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
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

type UserStatus string

const (
	UserStatusBlocked             UserStatus = "Blocked"
	UserStatusRegistrationPending UserStatus = "RegistrationPending"
	UserStatusActive              UserStatus = "Active"
	UserStatusInactive            UserStatus = "Inactive"
)

type Role string

const (
	RoleCustomer Role = "user"
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
	OTP         string
	OTPExpiry   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Roles []Role
}

func (u User) GetFirstname() string {
	return strings.Fields(u.FullName)[0]
}

type UserService interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUserStatus(ctx context.Context, status UserStatus, id int64) error
	UpdateUserOTP(ctx context.Context, id int64, otp string, otpExpiry int64) error
	UpdateUserPassword(ctx context.Context, id int64, password string) error
}

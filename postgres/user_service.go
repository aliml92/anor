package postgres

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/cart"
	"github.com/aliml92/anor/postgres/repository/order"
	"github.com/aliml92/anor/postgres/repository/user"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
	"log/slog"
	mrand "math/rand"
)

// Ensure service implements interface.
var _ anor.UserService = (*UserService)(nil)

type UserService struct {
	userRepository  user.Repository
	cartRepository  cart.Repository
	orderRepository order.Repository
}

func NewUserService(us user.Repository, cs cart.Repository, os order.Repository) *UserService {
	return &UserService{
		userRepository:  us,
		cartRepository:  cs,
		orderRepository: os,
	}
}

func (s UserService) Create(ctx context.Context, params anor.UserCreateParams) (anor.User, error) {
	if params.Provider == anor.ProviderBuiltin && params.Password == "" {
		return anor.User{}, oops.Errorf("password is required")
	}

	if params.Provider == anor.ProviderGoogle {
		params.Password = "google_" + generateRandomPassword(12)
		params.Status = anor.UserStatusActive
	}

	var (
		newUser *user.User
		err     error
	)
	err = s.userRepository.ExecTx(ctx, func(tx pgx.Tx) error {
		if params.Status != "" && params.PhoneNumber != "" {
			newUser, err = s.userRepository.WithTx(tx).CreateUserWithStatusAndPhone(ctx,
				params.Email,
				params.Password,
				params.Name,
				user.UserStatus(params.Status),
				&params.PhoneNumber,
			)
		} else if params.Status != "" {
			newUser, err = s.userRepository.WithTx(tx).CreateUserWithStatus(ctx,
				params.Email,
				params.Password,
				params.Name,
				user.UserStatus(params.Status),
			)
		} else if params.PhoneNumber != "" {
			newUser, err = s.userRepository.WithTx(tx).CreateUserWithPhone(ctx,
				params.Email,
				params.Password,
				params.Name,
				&params.PhoneNumber,
			)
		} else {
			newUser, err = s.userRepository.WithTx(tx).CreateUser(ctx, params.Email, params.Password, params.Name)
		}
		if err != nil {
			return oops.Wrap(err)
		}

		return nil
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return anor.User{}, anor.ErrUserExists
			}
		}

		return anor.User{}, oops.Errorf("failed to create user: %v", err)
	}

	u := anor.User{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Password:  newUser.Password,
		FullName:  newUser.FullName,
		Status:    anor.UserStatus(newUser.Status),
		CreatedAt: newUser.CreatedAt.Time,
		UpdatedAt: newUser.UpdatedAt.Time,
	}

	return u, nil
}

func generateRandomPassword(length int) string {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		slog.Warn("Error generating token", "error", err)
		for i := range randomBytes {
			randomBytes[i] = byte(mrand.Int())
		}
	}
	token := base64.URLEncoding.EncodeToString(randomBytes)
	return token
}

func (s UserService) Get(ctx context.Context, id int64) (anor.User, error) {
	ru, err := s.userRepository.GetUser(ctx, id) // retrieved user
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.User{}, anor.ErrNotFound
		}
		return anor.User{}, oops.Errorf("failed to get user: %v", err)
	}

	u := anor.User{
		ID:        ru.ID,
		Email:     ru.Email,
		Password:  ru.Password,
		FullName:  ru.FullName,
		Status:    anor.UserStatus(ru.Status),
		CreatedAt: ru.CreatedAt.Time,
		UpdatedAt: ru.UpdatedAt.Time,
	}

	if ru.PhoneNumber != nil {
		u.PhoneNumber = *ru.PhoneNumber
	}

	return u, nil
}

func (s UserService) GetByEmail(ctx context.Context, email string) (anor.User, error) {
	ru, err := s.userRepository.GetUserByEmail(ctx, email) // retrieved user
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.User{}, anor.ErrUserNotFound
		}
		return anor.User{}, oops.Wrap(err)
	}

	u := anor.User{
		ID:        ru.ID,
		Email:     ru.Email,
		Password:  ru.Password,
		FullName:  ru.FullName,
		Status:    anor.UserStatus(ru.Status),
		CreatedAt: ru.CreatedAt.Time,
		UpdatedAt: ru.UpdatedAt.Time,
	}

	if ru.PhoneNumber != nil {
		u.PhoneNumber = *ru.PhoneNumber
	}

	return u, nil
}

func (s UserService) UpdateStatusByEmail(ctx context.Context, status anor.UserStatus, email string) error {
	err := s.userRepository.UpdateUserStatusByEmail(ctx, user.UserStatus(status), email)
	if err != nil {
		return oops.Errorf("failed to update user status: %v", err)
	}

	return nil
}

func (s UserService) UpdatePassword(ctx context.Context, id int64, password string) error {
	err := s.userRepository.UpdateUserPassword(ctx, password, id)
	if err != nil {
		return oops.Errorf("failed to update user password: %v", err)
	}
	return nil
}

func (s UserService) GetActivityCounts(ctx context.Context, id int64) (anor.UserActivityCounts, error) {
	cartItemCount, err := s.cartRepository.CountCartItemsByUserIdAndCartStatus(ctx, &id, cart.CartStatusOpen)
	if err != nil {
		return anor.UserActivityCounts{}, oops.Errorf("failed to get user cart item counts: %v", err)
	}

	//orderCount, err := s.orderRepository.CountActiveOrders(ctx, &id)
	//if err != nil {
	//	return anor.UserActivityCounts{}, oops.Errorf("failed to get user order counts: %v", err)
	//}

	// TODO: count wishlist items here

	return anor.UserActivityCounts{
		CartItemsCount:     int(cartItemCount),
		WishlistItemsCount: 0, // FIX: Remove hardcoded value
		//ActiveOrdersCount:  int(orderCount),
	}, nil
}

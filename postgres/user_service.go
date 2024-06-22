package postgres

import (
	"context"
	"errors"
	"github.com/aliml92/anor/postgres/repository/cart"
	"github.com/aliml92/anor/postgres/repository/order"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/user"
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

func (s UserService) CreateUser(ctx context.Context, u anor.User) error {
	err := s.userRepository.ExecTx(ctx, func(tx pgx.Tx) error {
		err := s.userRepository.WithTx(tx).CreateUser(ctx, u.Email, u.Password, u.FullName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return anor.ErrUserExists
			}
		}

		return err
	}

	return nil
}

func (s UserService) GetUser(ctx context.Context, id int64) (anor.User, error) {
	ru, err := s.userRepository.GetUser(ctx, id) // retrieved user
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.User{}, anor.ErrNotFound
		}
		return anor.User{}, err
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

func (s UserService) GetUserByEmail(ctx context.Context, email string) (anor.User, error) {
	ru, err := s.userRepository.GetUserByEmail(ctx, email) // retrieved user
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.User{}, anor.ErrNotFound
		}
		return anor.User{}, err
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

func (s UserService) UpdateUserStatusByEmail(ctx context.Context, status anor.UserStatus, email string) error {
	err := s.userRepository.UpdateUserStatusByEmail(ctx, user.UserStatus(status), email)
	if err != nil {
		return err
	}

	return nil
}

func (s UserService) UpdateUserPassword(ctx context.Context, id int64, password string) error {
	err := s.userRepository.UpdateUserPassword(ctx, password, id)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) GetUserActivityCounts(ctx context.Context, id int64) (anor.UserActivityCounts, error) {
	cartItemCount, err := s.cartRepository.CountCartItems(ctx, &id)
	if err != nil {
		return anor.UserActivityCounts{}, err
	}

	orderCount, err := s.orderRepository.CountActiveOrders(ctx, &id)
	if err != nil {
		return anor.UserActivityCounts{}, err
	}

	// TODO: count wishlist items here

	return anor.UserActivityCounts{
		CartItemsCount:     int(cartItemCount),
		WishlistItemsCount: 0, // FIX: Remove hardcoded value
		ActiveOrdersCount:  int(orderCount),
	}, nil
}

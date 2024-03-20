package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/pkg/emailer"
	"github.com/aliml92/anor/pkg/utils"
)

type authService struct {
	userService anor.UserService
	emailer     emailer.Emailer
	session     *scs.SessionManager
}

func NewAuthService(us anor.UserService, e emailer.Emailer) anor.AuthService {
	return &authService{
		userService: us,
		emailer:     e,
	}
}

func (s *authService) Signup(ctx context.Context, name, email, password string) error {
	// Prepare data to save
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	email = strings.ToLower(email)
	otp := utils.GeneraterRandomCode(6)
	otpExp := time.Now().UTC().Add(2 * time.Minute).Unix()

	u := anor.User{
		Email:     email,
		Password:  hashedPassword,
		FullName:  name,
		OTP:       otp,
		OTPExpiry: otpExp,
	}

	if err := s.userService.CreateUser(ctx, u); err != nil {
		if errors.Is(err, anor.ErrUserExists) {
			return ErrEmailAlreadyTaken
		}

		return err
	}

	if err := s.emailer.SendVerificationMessageWithOTP(ctx, otp, email); err != nil {
		return err
	}

	return nil
}

func (s *authService) SignupConfirm(ctx context.Context, otp, email string) error {
	u, err := s.userService.GetUserByEmail(ctx, email) // retrieved user
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			return ErrInvalidOTP
		}
		return err
	}

	if otp != u.OTP {
		return ErrInvalidOTP
	}

	// send bad request if otp code has expired
	currentTime := time.Now().UTC().Unix()
	if u.OTPExpiry < currentTime {
		return ErrExpiredOTP
	}

	if err := s.userService.UpdateUserStatus(ctx, anor.UserStatusActive, u.ID); err != nil {
		return err
	}

	return nil
}

func (s *authService) Signin(ctx context.Context, email, password string) (int64, error) {
	var id int64
	u, err := s.userService.GetUserByEmail(ctx, email) // retrieved user
	if err != nil {
		if errors.Is(err, anor.ErrNotFound) {
			return id, ErrInvalidCredentials
		}
		return id, err
	}

	switch u.Status {
	case anor.UserStatusRegistrationPending:
		return id, ErrEmailNotConfirmed
	case anor.UserStatusBlocked:
		return id, ErrAccountBlocked
	case anor.UserStatusInactice:
		return id, ErrAccountInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return id, ErrInvalidCredentials
		}
		return id, err
	}

	return u.ID, nil
}

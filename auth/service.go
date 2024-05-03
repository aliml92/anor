package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aliml92/anor/cache"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/pkg/emailer"
	"github.com/aliml92/anor/pkg/utils"
)

type authService struct {
	userService anor.UserService
	emailer     emailer.Emailer
	cache       cache.ResetPasswordTokenCacher
}

func NewAuthService(us anor.UserService, e emailer.Emailer, cacher cache.ResetPasswordTokenCacher) anor.AuthService {
	return &authService{
		userService: us,
		emailer:     e,
		cache:       cacher,
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
	otpExp := time.Now().UTC().Add(5 * time.Minute).Unix()

	u := anor.User{
		Email:     email,
		Password:  hashedPassword,
		FullName:  name,
		OTP:       otp,
		OTPExpiry: otpExp,
	}

	if err := s.userService.CreateUser(ctx, u); err != nil {
		if errors.Is(err, anor.ErrUserExists) {
			// get user by email and check its status
			user, err := s.userService.GetUserByEmail(ctx, email)
			if err != nil {
				return err
			}
			switch user.Status {
			case anor.UserStatusActive:
				return ErrEmailAlreadyTaken
			case anor.UserStatusRegistrationPending:
				return s.saveAndResendOTP(ctx, user)
			}
		}
		return err
	}

	m, err := s.emailer.NewMessage("Anor account details confirmation", email, "email-verification-otp.gohtml", otp)
	if err != nil {
		return err
	}

	if err := s.emailer.Send(ctx, m); err != nil {
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
	case anor.UserStatusInactive:
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

func (s *authService) ResendOTP(ctx context.Context, email string) error {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	otp := utils.GeneraterRandomCode(6)
	otpExp := time.Now().UTC().Add(2 * time.Minute).Unix()

	user.OTP = otp
	user.OTPExpiry = otpExp
	return s.saveAndResendOTP(ctx, user)
}

func (s *authService) SendResetPasswordLink(ctx context.Context, email string) error {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	token, err := generateToken()
	if err != nil {
		return err
	}
	fmt.Printf("token: %s\n", token)
	fmt.Printf("user: %v\n", user)
	err = s.cache.CacheToken(ctx, token, user.ID, 10*time.Minute)
	if err != nil {
		return err
	}

	resetURL := "http://localhost:8008/auth/reset-password?token=" + token
	m, err := s.emailer.NewMessage("Anor | Reset Password", user.Email, "reset-password.gohtml", resetURL)
	if err != nil {
		return err
	}

	return s.emailer.Send(ctx, m)
}

func (s *authService) VerifyResetPasswordToken(ctx context.Context, token string) (bool, error) {
	ok, err := s.cache.TokenExists(ctx, token)
	return ok, err
}

func (s *authService) ResetPassword(ctx context.Context, token string, password string) error {
	userID, err := s.cache.GetUserIDByToken(ctx, token)
	if err != nil {
		return err
	}
	if userID == 0 {
		return ErrInvalidOrExpiredResetURL
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	if err := s.userService.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (s *authService) saveAndResendOTP(ctx context.Context, user anor.User) error {
	if err := s.userService.UpdateUserOTP(ctx, user.ID, user.OTP, user.OTPExpiry); err != nil {
		return err
	}

	m, err := s.emailer.NewMessage("Anor | Signup Confirmation", user.Email, "email-verification-otp.gohtml", user.OTP)
	if err != nil {
		return err
	}

	err = s.emailer.Send(ctx, m)

	return err
}

func generateToken() (string, error) {
	tokenLength := 16
	randomBytes := make([]byte, tokenLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(randomBytes)
	return token, nil
}

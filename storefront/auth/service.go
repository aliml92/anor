package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/aliml92/anor/cache/auth"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/email"
	"github.com/aliml92/anor/pkg/utils"
)

type authService struct {
	userService                 anor.UserService
	emailer                     email.Emailer
	signupConfirmationOTPCacher auth.SignupConfirmationOTPCache
	resetPasswordTokenCacher    auth.ResetPasswordTokenCache
}

func NewAuthService(us anor.UserService, e email.Emailer, otpCacher auth.SignupConfirmationOTPCache, tokenCacher auth.ResetPasswordTokenCache) anor.AuthService {
	return &authService{
		userService:                 us,
		emailer:                     e,
		signupConfirmationOTPCacher: otpCacher,
		resetPasswordTokenCacher:    tokenCacher,
	}
}

func (s *authService) Signup(ctx context.Context, name, email, password string) error {
	// Prepare data to save
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	email = strings.ToLower(email)
	u := anor.User{
		Email:    email,
		Password: hashedPassword,
		FullName: name,
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
				otp := utils.GeneraterRandomCode(6)
				return s.cacheAndResendOTP(ctx, otp, email)
			}
		}
		return err
	}

	otp := utils.GeneraterRandomCode(6)
	if err := s.signupConfirmationOTPCacher.CacheOTP(ctx, otp, email, 5*time.Minute); err != nil {
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
	e, err := s.signupConfirmationOTPCacher.GetEmailByOTP(ctx, otp)
	if err != nil {
		return err
	}
	if email != e {
		return ErrInvalidOTP
	}

	if err := s.userService.UpdateUserStatusByEmail(ctx, anor.UserStatusActive, email); err != nil {
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
	_, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	otp := utils.GeneraterRandomCode(6)

	return s.cacheAndResendOTP(ctx, otp, email)
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
	err = s.resetPasswordTokenCacher.CacheToken(ctx, token, user.ID, 10*time.Minute)
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
	ok, err := s.resetPasswordTokenCacher.TokenExists(ctx, token)
	return ok, err
}

func (s *authService) ResetPassword(ctx context.Context, token string, password string) error {
	userID, err := s.resetPasswordTokenCacher.GetUserIDByToken(ctx, token)
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

func (s *authService) GetUser(ctx context.Context, id int64) (*anor.User, error) {
	user, err := s.userService.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *authService) cacheAndResendOTP(ctx context.Context, otp, email string) error {
	if err := s.signupConfirmationOTPCacher.CacheOTP(ctx, otp, email, 5*time.Minute); err != nil {
		return err
	}

	m, err := s.emailer.NewMessage("Anor | Signup Confirmation", email, "email-verification-otp.gohtml", otp)
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

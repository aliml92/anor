package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aliml92/anor/cache/auth"
	"log/slog"
	mrand "math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/email"
	"github.com/aliml92/anor/pkg/utils"
)

type OTPGenerator func() string
type TokenGenerator func() string

type ServiceConfig struct {
	UserService                 anor.UserService
	Emailer                     email.Emailer
	SignupConfirmationOTPCacher auth.SignupConfirmationOTPCache
	ResetPasswordTokenCacher    auth.ResetPasswordTokenCache
	OTPGenerator                OTPGenerator
	OTPExpiration               time.Duration
	TokenGenerator              TokenGenerator
	TokenExpiration             time.Duration
	ServerURL                   string
}

type authService struct {
	userService                 anor.UserService
	emailer                     email.Emailer
	signupConfirmationOTPCacher auth.SignupConfirmationOTPCache
	resetPasswordTokenCacher    auth.ResetPasswordTokenCache
	otpGenerator                OTPGenerator
	otpExpiration               time.Duration
	tokenGenerator              TokenGenerator
	tokenExpiration             time.Duration

	serverURL string
}

func NewAuthService(config ServiceConfig) anor.AuthService {
	if config.OTPGenerator == nil {
		config.OTPGenerator = func() string {
			return generateOTP(6)
		}
	}

	if config.OTPExpiration == 0 {
		config.OTPExpiration = 5 * time.Minute
	}

	if config.TokenGenerator == nil {
		config.TokenGenerator = func() string {
			return generateToken(16)
		}
	}

	if config.TokenExpiration == 0 {
		config.TokenExpiration = 10 * time.Minute
	}

	if config.ServerURL == "" {
		config.ServerURL = "http://localhost:8008"
	}

	return &authService{
		userService:                 config.UserService,
		emailer:                     config.Emailer,
		signupConfirmationOTPCacher: config.SignupConfirmationOTPCacher,
		resetPasswordTokenCacher:    config.ResetPasswordTokenCacher,
		otpGenerator:                config.OTPGenerator,
		tokenGenerator:              config.TokenGenerator,
	}
}

func (s *authService) Signup(ctx context.Context, name, email, password string) (anor.User, error) {
	// Prepare data to save
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return anor.User{}, err
	}
	email = strings.ToLower(email)

	u, err := s.userService.Create(ctx, anor.UserCreateParams{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Provider: anor.ProviderBuiltin,
		Status:   anor.UserStatusRegistrationPending,
	})
	if err != nil {
		if errors.Is(err, anor.ErrUserExists) {
			// get user by email and check its status
			user, err := s.userService.GetByEmail(ctx, email)
			if err != nil {
				return anor.User{}, err
			}

			// TODO: temporary solution
			// check if the user is created via oauth2
			if strings.HasPrefix(u.Password, "google_") {
				return anor.User{}, ErrOAuth2RegisteredAccount
			}

			switch user.Status {
			case anor.UserStatusActive:
				return anor.User{}, ErrEmailAlreadyTaken
			case anor.UserStatusRegistrationPending:
				otp := utils.GeneraterRandomCode(6)
				err = s.signupConfirmationOTPCacher.CacheOTP(ctx, otp, email, 5*time.Minute)
				if err != nil {
					return anor.User{}, err
				}

				err = s.sendEmail(ctx, "Anor | Signup Confirmation", email, "email_verification_otp.gohtml", otp)
				return anor.User{}, err
			}
		}
		return anor.User{}, err
	}

	otp := s.otpGenerator()
	if err := s.signupConfirmationOTPCacher.CacheOTP(ctx, otp, email, s.otpExpiration); err != nil {
		return u, err
	}

	err = s.sendEmail(ctx, "Anor account details confirmation", email, "email_verification_otp.gohtml", otp)
	if err != nil {
		return u, err
	}

	return u, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *authService) SignupConfirm(ctx context.Context, otp, email string) error {
	e, err := s.signupConfirmationOTPCacher.GetEmailByOTP(ctx, otp)
	if err != nil {
		return err
	}
	if email != e {
		return ErrInvalidOTP
	}

	if err := s.userService.UpdateStatusByEmail(ctx, anor.UserStatusActive, email); err != nil {
		return err
	}

	return nil
}

func (s *authService) Signin(ctx context.Context, email, password string) (anor.User, error) {
	var u anor.User
	u, err := s.userService.GetByEmail(ctx, email) // retrieved user
	if err != nil {
		if errors.Is(err, anor.ErrUserNotFound) {
			return u, ErrInvalidCredentials
		}
		return u, err
	}

	switch u.Status {
	case anor.UserStatusRegistrationPending:
		return u, ErrEmailNotConfirmed
	case anor.UserStatusBlocked:
		return u, ErrAccountBlocked
	}

	// TODO: temporary solution
	// check if the user is created via oauth2
	if strings.HasPrefix(u.Password, "google_") {
		return u, ErrOAuth2RegisteredAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return u, ErrInvalidCredentials
		}
		return u, err
	}

	return u, nil
}

func (s *authService) ResendOTP(ctx context.Context, email string) error {
	_, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	otp := utils.GeneraterRandomCode(6)
	err = s.signupConfirmationOTPCacher.CacheOTP(ctx, otp, email, s.otpExpiration)
	if err != nil {
		return err
	}

	err = s.sendEmail(ctx, "Anor | Signup Confirmation", email, "email_verification_otp.gohtml", otp)
	return err

}

func (s *authService) SendResetPasswordLink(ctx context.Context, email string) error {
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	token := s.tokenGenerator()
	err = s.resetPasswordTokenCacher.CacheToken(ctx, token, user.ID, s.tokenExpiration)
	if err != nil {
		return err
	}

	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", s.serverURL, token)
	err = s.sendEmail(ctx, "Anor | Reset Password", user.Email, "reset_password.gohtml", resetURL)
	return err
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

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	if err := s.userService.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (s *authService) SigninWithGoogle(ctx context.Context, email, name string) (anor.User, error) {
	u, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, anor.ErrUserNotFound) {
			// User doesn't exist in our system, so create a new one
			newUser, err := s.userService.Create(ctx, anor.UserCreateParams{
				Name:     name,
				Email:    email,
				Provider: anor.ProviderGoogle,
				Status:   anor.UserStatusActive,
			})
			if err != nil {
				return u, err
			}

			return newUser, nil
		}
		return u, err
	}

	if u.Status == anor.UserStatusBlocked {
		return u, ErrAccountBlocked
	}

	return u, nil

}

func (s *authService) GetUser(ctx context.Context, id int64) (*anor.User, error) {
	u, err := s.userService.Get(ctx, id)
	return &u, err
}

func (s *authService) sendEmail(ctx context.Context, subject, recipient, template string, data interface{}) error {
	m, err := s.emailer.NewMessage(subject, recipient, template, data)
	if err != nil {
		return err
	}
	return s.emailer.Send(ctx, m)
}

func generateToken(length int) string {
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

func generateOTP(n int) string {
	const digits = "0123456789"
	otp := make([]byte, n)

	for i := 0; i < n; i++ {
		randomDigit := digits[mrand.Intn(len(digits))]
		otp[i] = randomDigit
	}

	return string(otp)
}

package auth

import "errors"

var (
	ErrEmailAlreadyTaken        = errors.New("email already taken")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrInvalidOTP               = errors.New("invalid OTP")
	ErrEmailNotConfirmed        = errors.New("email not confirmed")
	ErrAccountBlocked           = errors.New("account blocked")
	ErrInvalidOrExpiredResetURL = errors.New("invalid or expired reset URL")
	ErrOAuth2RegisteredAccount  = errors.New("OAuth2 registered account")
)

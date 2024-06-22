package auth

import "errors"

var (
	ErrEmailAlreadyTaken        = errors.New("Email already taken")
	ErrInvalidCredentials       = errors.New("Invalid credentials")
	ErrInvalidOTP               = errors.New("Invalid OTP")
	ErrExpiredOTP               = errors.New("Expired OTP")
	ErrEmailNotConfirmed        = errors.New("Email not confirmed")
	ErrAccountBlocked           = errors.New("Account blocked")
	ErrAccountInactive          = errors.New("Account inactive")
	ErrInvalidOrExpiredResetURL = errors.New("Invalid or expired reset URL")
)

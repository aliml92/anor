package emailer

import "context"

type Emailer interface {
	SendVerificationMessageWithOTP(ctx context.Context, otp string, email string) error
}

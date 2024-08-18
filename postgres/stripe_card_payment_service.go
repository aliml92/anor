package postgres

import (
	"context"
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/payment"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/oops"
)

var _ anor.StripePaymentService = (*StripePaymentService)(nil)

type StripePaymentService struct {
	paymentRepository payment.Repository
}

func NewStipePaymentService(pr payment.Repository) *StripePaymentService {
	return &StripePaymentService{
		paymentRepository: pr,
	}
}

func (s *StripePaymentService) GetByOrderID(ctx context.Context, orderID int64) (anor.StripeCardPayment, error) {
	p, err := s.paymentRepository.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.StripeCardPayment{}, anor.ErrPaymentNotFound
		}
		return anor.StripeCardPayment{}, oops.Wrap(err)
	}

	res := anor.StripeCardPayment{
		ID:               p.ID,
		OrderID:          p.OrderID,
		UserID:           anor.Int64Value(p.UserID),
		BillingAddressID: p.BillingAddressID,
		PaymentIntentID:  p.PaymentIntentID,
		PaymentMethodID:  anor.StringValue(p.PaymentMethodID),
		Amount:           p.Amount,
		Currency:         p.Currency,
		Status:           p.Status,
		ClientSecret:     anor.StringValue(p.ClientSecret),
		LastError:        anor.StringValue(p.LastError),
		CreatedAt:        p.CreatedAt.Time,
		UpdatedAt:        p.UpdatedAt.Time,
	}

	return res, nil
}

func (s *StripePaymentService) Create(ctx context.Context, params anor.StripePaymentCreateParams) error {
	err := s.paymentRepository.Create(ctx,
		params.OrderID,
		anor.Int64(params.UserID),
		params.BillingAddressID,
		params.PaymentIntentID,
		anor.String(params.PaymentMethodID),
		params.Amount,
		params.Currency,
		params.Status,
		anor.String(params.ClientSecret),
		anor.String(params.LastError),
		params.CardLast4,
		params.CardBrand,
		pgtype.Timestamptz{
			Time:  params.CreatedAt,
			Valid: true,
		},
	)

	return oops.Wrap(err)
}

package postgres

import (
	"context"
	"errors"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/postgres/repository/address"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/samber/oops"
)

var _ anor.AddressService = (*AddressService)(nil)

type AddressService struct {
	addressRepository address.Repository
}

func NewAddressService(ar address.Repository) *AddressService {
	return &AddressService{
		addressRepository: ar,
	}
}

func (s *AddressService) Create(ctx context.Context, params anor.AddressCreateParams) (anor.Address, error) {
	defaultFor := address.NullAddressDefaultType{}
	switch params.DefaultFor {
	case anor.AddressDefaultTypeShipping:
		defaultFor.AddressDefaultType = address.AddressDefaultTypeShipping
	case anor.AddressDefaultTypeBilling:
		defaultFor.AddressDefaultType = address.AddressDefaultTypeBilling
	}

	a, err := s.addressRepository.CreateAddress(ctx,
		params.UserID,
		defaultFor,
		params.Name,
		params.AddressLine1,
		anor.String(params.AddressLine2),
		params.City,
		anor.String(params.StateProvince),
		anor.String(params.PostalCode),
		anor.String(params.Country),
	)
	if err != nil {
		return anor.Address{}, oops.Wrap(err)
	}

	return anor.Address{
		ID:            a.ID,
		UserID:        a.UserID,
		DefaultFor:    anor.AddressDefaultType(a.DefaultFor.AddressDefaultType),
		Name:          a.Name,
		AddressLine1:  a.AddressLine1,
		AddressLine2:  anor.StringValue(a.AddressLine2),
		City:          a.City,
		StateProvince: anor.StringValue(a.StateProvince),
		PostalCode:    anor.StringValue(a.PostalCode),
		Country:       anor.StringValue(a.Country),
		Phone:         anor.StringValue(a.Phone),
		CreatedAt:     a.CreatedAt.Time,
		UpdatedAt:     a.UpdatedAt.Time,
	}, nil
}

func (s *AddressService) Get(ctx context.Context, id int64) (anor.Address, error) {
	a, err := s.addressRepository.GetAddressByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.Address{}, anor.ErrAddressNotFound
		}

		return anor.Address{}, oops.Wrap(err)
	}

	return anor.Address{
		ID:            a.ID,
		UserID:        a.UserID,
		DefaultFor:    anor.AddressDefaultType(a.DefaultFor.AddressDefaultType),
		Name:          a.Name,
		AddressLine1:  a.AddressLine1,
		AddressLine2:  anor.StringValue(a.AddressLine2),
		City:          a.City,
		StateProvince: anor.StringValue(a.StateProvince),
		PostalCode:    anor.StringValue(a.PostalCode),
		Country:       anor.StringValue(a.Country),
		Phone:         anor.StringValue(a.Phone),
		CreatedAt:     a.CreatedAt.Time,
		UpdatedAt:     a.UpdatedAt.Time,
	}, nil
}

func (s *AddressService) GetDefault(ctx context.Context, userID int64, defaultFor anor.AddressDefaultType) (anor.Address, error) {
	d := address.NullAddressDefaultType{
		AddressDefaultType: address.AddressDefaultType(defaultFor),
		Valid:              true,
	}
	a, err := s.addressRepository.GetUserDefaultAddress(ctx, userID, d)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anor.Address{}, anor.ErrAddressNotFound
		}

		return anor.Address{}, oops.Wrap(err)
	}

	return anor.Address{
		ID:            a.ID,
		UserID:        a.UserID,
		DefaultFor:    anor.AddressDefaultType(a.DefaultFor.AddressDefaultType),
		Name:          a.Name,
		AddressLine1:  a.AddressLine1,
		AddressLine2:  anor.StringValue(a.AddressLine2),
		City:          a.City,
		StateProvince: anor.StringValue(a.StateProvince),
		PostalCode:    anor.StringValue(a.PostalCode),
		Country:       anor.StringValue(a.Country),
		Phone:         anor.StringValue(a.Phone),
		CreatedAt:     a.CreatedAt.Time,
		UpdatedAt:     a.UpdatedAt.Time,
	}, nil
}

func (s *AddressService) List(ctx context.Context, params anor.AddressListParams) ([]anor.Address, error) {
	limit := lo.Ternary[int32](params.PageSize > 0, int32(params.PageSize), 5)
	offset := limit * (int32(params.Page) - 1)

	addresses, err := s.addressRepository.ListAddressesByUserID(ctx, params.UserID, offset, limit)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	result := make([]anor.Address, len(addresses))
	for i, a := range addresses {
		result[i] = anor.Address{
			ID:            a.ID,
			UserID:        a.UserID,
			DefaultFor:    anor.AddressDefaultType(a.DefaultFor.AddressDefaultType),
			Name:          a.Name,
			AddressLine1:  a.AddressLine1,
			AddressLine2:  anor.StringValue(a.AddressLine2),
			City:          a.City,
			StateProvince: anor.StringValue(a.StateProvince),
			PostalCode:    anor.StringValue(a.PostalCode),
			Country:       anor.StringValue(a.Country),
			Phone:         anor.StringValue(a.Phone),
			CreatedAt:     a.CreatedAt.Time,
			UpdatedAt:     a.UpdatedAt.Time,
		}
	}

	return result, nil
}

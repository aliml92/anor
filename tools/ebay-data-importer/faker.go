package main

import (
	"context"
	"math/rand"

	"github.com/brianvoe/gofakeit"
	"github.com/gosimple/slug"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/pkg/utils"
	"github.com/aliml92/anor/postgres/store/sellerstore"
	"github.com/aliml92/anor/postgres/store/user"
	ts2 "github.com/aliml92/anor/typesense"
)

const (
	fakeUserPassword = "Password1@"
)

var discounts = []float32{
	0, 0, 0.02, 0, 0, 0, 0.03, 0, 0.04, 0, 0, 0.05, 0, 0, 0.07, 0, 0, 0, 0.08, 0, 0, 0, 0.09, 0,
	0.10, 0, 0.11, 0, 0.12, 0, 0.13, 0, 0.14, 0, 0.15, 0, 0.20, 0, 0, 0.25, 0, 0, 0, 0, 0, 0, 0,
	0.30, 0, 0, 0.35, 0, 0.40, 0, 0.45, 0, 0.50, 0.55, 0, 0.60, 0, 0.65, 0, 0.70, 0, 0, 0, 0, 0,
	0.15, 0, 0.12, 0, 0, 0, 0.05, 0.07, 0.20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

func generateRandomDiscount() float32 {
	idx := rand.Intn(len(discounts))
	return discounts[idx]
}

func createSellerUsers(ctx context.Context, us user.Store, n int) ([]int64, error) {
	userIDs := make([]int64, n)
	for i := 0; i < n; i++ {
		// save a default user and get its id
		hashedPassword, _ := utils.HashPassword(fakeUserPassword)

		email := gofakeit.Email()
		fname := gofakeit.Name()
		status := user.UserStatusActive

		userID, err := us.CreateSeller(ctx, email, hashedPassword, fname, status)
		if err != nil {
			return nil, err
		}

		userIDs[i] = userID
	}

	return userIDs, nil
}

func createSellerStores(ctx context.Context, ss sellerstore.Store, userIDs []int64, searcher *ts2.Searcher) ([]int32, error) {
	sellerStoreIDs := make([]int32, len(userIDs))
	for index, userID := range userIDs {
		// save a default store and get its id
		name := gofakeit.Company()
		publicID := slug.Make(name)
		description := gofakeit.Sentence(20)

		storeID, err := ss.CreateStore(ctx, publicID, userID, name, description)
		if err != nil {
			return nil, err
		}

		err = searcher.IndexSellerStore(ctx, anor.SellerStore{
			ID:       storeID,
			Name:     name,
			PublicID: publicID,
		})
		if err != nil {
			return nil, err
		}

		sellerStoreIDs[index] = storeID
	}

	return sellerStoreIDs, nil
}

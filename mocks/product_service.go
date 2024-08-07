// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/aliml92/anor"
	"github.com/shopspring/decimal"
	"sync"
)

// Ensure, that ProductServiceMock does implement anor.ProductService.
// If this is not the case, regenerate this file with moq.
var _ anor.ProductService = &ProductServiceMock{}

// ProductServiceMock is a mock implementation of anor.ProductService.
//
//	func TestSomethingThatUsesProductService(t *testing.T) {
//
//		// make and configure a mocked anor.ProductService
//		mockedProductService := &ProductServiceMock{
//			GetMinMaxPricesByCategoryFunc: func(ctx context.Context, category anor.Category) ([2]decimal.Decimal, error) {
//				panic("mock out the GetMinMaxPricesByCategory method")
//			},
//			GetNewArrivalsFunc: func(ctx context.Context, limit int) ([]anor.Product, error) {
//				panic("mock out the GetNewArrivals method")
//			},
//			GetProductFunc: func(ctx context.Context, id int64) (*anor.Product, error) {
//				panic("mock out the GetProduct method")
//			},
//			GetProductBrandsByCategoryFunc: func(ctx context.Context, category anor.Category) ([]string, error) {
//				panic("mock out the GetProductBrandsByCategory method")
//			},
//			GetProductsByCategoryFunc: func(ctx context.Context, category anor.Category, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error) {
//				panic("mock out the GetProductsByCategory method")
//			},
//			GetProductsByLeafCategoryIDFunc: func(ctx context.Context, categoryID int32, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error) {
//				panic("mock out the GetProductsByLeafCategoryID method")
//			},
//			GetProductsByNonLeafCategoryIDFunc: func(ctx context.Context, categoryID int32, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error) {
//				panic("mock out the GetProductsByNonLeafCategoryID method")
//			},
//		}
//
//		// use mockedProductService in code that requires anor.ProductService
//		// and then make assertions.
//
//	}
type ProductServiceMock struct {
	// GetMinMaxPricesByCategoryFunc mocks the GetMinMaxPricesByCategory method.
	GetMinMaxPricesByCategoryFunc func(ctx context.Context, category anor.Category) ([2]decimal.Decimal, error)

	// GetNewArrivalsFunc mocks the GetNewArrivals method.
	GetNewArrivalsFunc func(ctx context.Context, limit int) ([]anor.Product, error)

	// GetProductFunc mocks the GetProduct method.
	GetProductFunc func(ctx context.Context, id int64) (*anor.Product, error)

	// GetProductBrandsByCategoryFunc mocks the GetProductBrandsByCategory method.
	GetProductBrandsByCategoryFunc func(ctx context.Context, category anor.Category) ([]string, error)

	// GetProductsByCategoryFunc mocks the GetProductsByCategory method.
	GetProductsByCategoryFunc func(ctx context.Context, category anor.Category, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error)

	// GetProductsByLeafCategoryIDFunc mocks the GetProductsByLeafCategoryID method.
	GetProductsByLeafCategoryIDFunc func(ctx context.Context, categoryID int32, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error)

	// GetProductsByNonLeafCategoryIDFunc mocks the GetProductsByNonLeafCategoryID method.
	GetProductsByNonLeafCategoryIDFunc func(ctx context.Context, categoryID int32, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetMinMaxPricesByCategory holds details about calls to the GetMinMaxPricesByCategory method.
		GetMinMaxPricesByCategory []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Category is the category argument value.
			Category anor.Category
		}
		// GetNewArrivals holds details about calls to the GetNewArrivals method.
		GetNewArrivals []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Limit is the limit argument value.
			Limit int
		}
		// GetProduct holds details about calls to the GetProduct method.
		GetProduct []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
		}
		// GetProductBrandsByCategory holds details about calls to the GetProductBrandsByCategory method.
		GetProductBrandsByCategory []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Category is the category argument value.
			Category anor.Category
		}
		// GetProductsByCategory holds details about calls to the GetProductsByCategory method.
		GetProductsByCategory []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Category is the category argument value.
			Category anor.Category
			// Params is the params argument value.
			Params anor.GetProductsByCategoryParams
		}
		// GetProductsByLeafCategoryID holds details about calls to the GetProductsByLeafCategoryID method.
		GetProductsByLeafCategoryID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// CategoryID is the categoryID argument value.
			CategoryID int32
			// Params is the params argument value.
			Params anor.GetProductsByCategoryParams
		}
		// GetProductsByNonLeafCategoryID holds details about calls to the GetProductsByNonLeafCategoryID method.
		GetProductsByNonLeafCategoryID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// CategoryID is the categoryID argument value.
			CategoryID int32
			// Params is the params argument value.
			Params anor.GetProductsByCategoryParams
		}
	}
	lockGetMinMaxPricesByCategory      sync.RWMutex
	lockGetNewArrivals                 sync.RWMutex
	lockGetProduct                     sync.RWMutex
	lockGetProductBrandsByCategory     sync.RWMutex
	lockGetProductsByCategory          sync.RWMutex
	lockGetProductsByLeafCategoryID    sync.RWMutex
	lockGetProductsByNonLeafCategoryID sync.RWMutex
}

// GetMinMaxPricesByCategory calls GetMinMaxPricesByCategoryFunc.
func (mock *ProductServiceMock) GetMinMaxPricesByCategory(ctx context.Context, category anor.Category) ([2]decimal.Decimal, error) {
	if mock.GetMinMaxPricesByCategoryFunc == nil {
		panic("ProductServiceMock.GetMinMaxPricesByCategoryFunc: method is nil but ProductService.GetMinMaxPricesByCategory was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Category anor.Category
	}{
		Ctx:      ctx,
		Category: category,
	}
	mock.lockGetMinMaxPricesByCategory.Lock()
	mock.calls.GetMinMaxPricesByCategory = append(mock.calls.GetMinMaxPricesByCategory, callInfo)
	mock.lockGetMinMaxPricesByCategory.Unlock()
	return mock.GetMinMaxPricesByCategoryFunc(ctx, category)
}

// GetMinMaxPricesByCategoryCalls gets all the calls that were made to GetMinMaxPricesByCategory.
// Check the length with:
//
//	len(mockedProductService.GetMinMaxPricesByCategoryCalls())
func (mock *ProductServiceMock) GetMinMaxPricesByCategoryCalls() []struct {
	Ctx      context.Context
	Category anor.Category
} {
	var calls []struct {
		Ctx      context.Context
		Category anor.Category
	}
	mock.lockGetMinMaxPricesByCategory.RLock()
	calls = mock.calls.GetMinMaxPricesByCategory
	mock.lockGetMinMaxPricesByCategory.RUnlock()
	return calls
}

// GetNewArrivals calls GetNewArrivalsFunc.
func (mock *ProductServiceMock) GetNewArrivals(ctx context.Context, limit int) ([]anor.Product, error) {
	if mock.GetNewArrivalsFunc == nil {
		panic("ProductServiceMock.GetNewArrivalsFunc: method is nil but ProductService.GetNewArrivals was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Limit int
	}{
		Ctx:   ctx,
		Limit: limit,
	}
	mock.lockGetNewArrivals.Lock()
	mock.calls.GetNewArrivals = append(mock.calls.GetNewArrivals, callInfo)
	mock.lockGetNewArrivals.Unlock()
	return mock.GetNewArrivalsFunc(ctx, limit)
}

// GetNewArrivalsCalls gets all the calls that were made to GetNewArrivals.
// Check the length with:
//
//	len(mockedProductService.GetNewArrivalsCalls())
func (mock *ProductServiceMock) GetNewArrivalsCalls() []struct {
	Ctx   context.Context
	Limit int
} {
	var calls []struct {
		Ctx   context.Context
		Limit int
	}
	mock.lockGetNewArrivals.RLock()
	calls = mock.calls.GetNewArrivals
	mock.lockGetNewArrivals.RUnlock()
	return calls
}

// GetProduct calls GetProductFunc.
func (mock *ProductServiceMock) GetProduct(ctx context.Context, id int64) (*anor.Product, error) {
	if mock.GetProductFunc == nil {
		panic("ProductServiceMock.GetProductFunc: method is nil but ProductService.GetProduct was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  int64
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetProduct.Lock()
	mock.calls.GetProduct = append(mock.calls.GetProduct, callInfo)
	mock.lockGetProduct.Unlock()
	return mock.GetProductFunc(ctx, id)
}

// GetProductCalls gets all the calls that were made to GetProduct.
// Check the length with:
//
//	len(mockedProductService.GetProductCalls())
func (mock *ProductServiceMock) GetProductCalls() []struct {
	Ctx context.Context
	ID  int64
} {
	var calls []struct {
		Ctx context.Context
		ID  int64
	}
	mock.lockGetProduct.RLock()
	calls = mock.calls.GetProduct
	mock.lockGetProduct.RUnlock()
	return calls
}

// GetProductBrandsByCategory calls GetProductBrandsByCategoryFunc.
func (mock *ProductServiceMock) GetProductBrandsByCategory(ctx context.Context, category anor.Category) ([]string, error) {
	if mock.GetProductBrandsByCategoryFunc == nil {
		panic("ProductServiceMock.GetProductBrandsByCategoryFunc: method is nil but ProductService.GetProductBrandsByCategory was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Category anor.Category
	}{
		Ctx:      ctx,
		Category: category,
	}
	mock.lockGetProductBrandsByCategory.Lock()
	mock.calls.GetProductBrandsByCategory = append(mock.calls.GetProductBrandsByCategory, callInfo)
	mock.lockGetProductBrandsByCategory.Unlock()
	return mock.GetProductBrandsByCategoryFunc(ctx, category)
}

// GetProductBrandsByCategoryCalls gets all the calls that were made to GetProductBrandsByCategory.
// Check the length with:
//
//	len(mockedProductService.GetProductBrandsByCategoryCalls())
func (mock *ProductServiceMock) GetProductBrandsByCategoryCalls() []struct {
	Ctx      context.Context
	Category anor.Category
} {
	var calls []struct {
		Ctx      context.Context
		Category anor.Category
	}
	mock.lockGetProductBrandsByCategory.RLock()
	calls = mock.calls.GetProductBrandsByCategory
	mock.lockGetProductBrandsByCategory.RUnlock()
	return calls
}

// GetProductsByCategory calls GetProductsByCategoryFunc.
func (mock *ProductServiceMock) GetProductsByCategory(ctx context.Context, category anor.Category, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error) {
	if mock.GetProductsByCategoryFunc == nil {
		panic("ProductServiceMock.GetProductsByCategoryFunc: method is nil but ProductService.GetProductsByCategory was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Category anor.Category
		Params   anor.GetProductsByCategoryParams
	}{
		Ctx:      ctx,
		Category: category,
		Params:   params,
	}
	mock.lockGetProductsByCategory.Lock()
	mock.calls.GetProductsByCategory = append(mock.calls.GetProductsByCategory, callInfo)
	mock.lockGetProductsByCategory.Unlock()
	return mock.GetProductsByCategoryFunc(ctx, category, params)
}

// GetProductsByCategoryCalls gets all the calls that were made to GetProductsByCategory.
// Check the length with:
//
//	len(mockedProductService.GetProductsByCategoryCalls())
func (mock *ProductServiceMock) GetProductsByCategoryCalls() []struct {
	Ctx      context.Context
	Category anor.Category
	Params   anor.GetProductsByCategoryParams
} {
	var calls []struct {
		Ctx      context.Context
		Category anor.Category
		Params   anor.GetProductsByCategoryParams
	}
	mock.lockGetProductsByCategory.RLock()
	calls = mock.calls.GetProductsByCategory
	mock.lockGetProductsByCategory.RUnlock()
	return calls
}

// GetProductsByLeafCategoryID calls GetProductsByLeafCategoryIDFunc.
func (mock *ProductServiceMock) GetProductsByLeafCategoryID(ctx context.Context, categoryID int32, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error) {
	if mock.GetProductsByLeafCategoryIDFunc == nil {
		panic("ProductServiceMock.GetProductsByLeafCategoryIDFunc: method is nil but ProductService.GetProductsByLeafCategoryID was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		CategoryID int32
		Params     anor.GetProductsByCategoryParams
	}{
		Ctx:        ctx,
		CategoryID: categoryID,
		Params:     params,
	}
	mock.lockGetProductsByLeafCategoryID.Lock()
	mock.calls.GetProductsByLeafCategoryID = append(mock.calls.GetProductsByLeafCategoryID, callInfo)
	mock.lockGetProductsByLeafCategoryID.Unlock()
	return mock.GetProductsByLeafCategoryIDFunc(ctx, categoryID, params)
}

// GetProductsByLeafCategoryIDCalls gets all the calls that were made to GetProductsByLeafCategoryID.
// Check the length with:
//
//	len(mockedProductService.GetProductsByLeafCategoryIDCalls())
func (mock *ProductServiceMock) GetProductsByLeafCategoryIDCalls() []struct {
	Ctx        context.Context
	CategoryID int32
	Params     anor.GetProductsByCategoryParams
} {
	var calls []struct {
		Ctx        context.Context
		CategoryID int32
		Params     anor.GetProductsByCategoryParams
	}
	mock.lockGetProductsByLeafCategoryID.RLock()
	calls = mock.calls.GetProductsByLeafCategoryID
	mock.lockGetProductsByLeafCategoryID.RUnlock()
	return calls
}

// GetProductsByNonLeafCategoryID calls GetProductsByNonLeafCategoryIDFunc.
func (mock *ProductServiceMock) GetProductsByNonLeafCategoryID(ctx context.Context, categoryID int32, params anor.GetProductsByCategoryParams) ([]anor.Product, int64, error) {
	if mock.GetProductsByNonLeafCategoryIDFunc == nil {
		panic("ProductServiceMock.GetProductsByNonLeafCategoryIDFunc: method is nil but ProductService.GetProductsByNonLeafCategoryID was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		CategoryID int32
		Params     anor.GetProductsByCategoryParams
	}{
		Ctx:        ctx,
		CategoryID: categoryID,
		Params:     params,
	}
	mock.lockGetProductsByNonLeafCategoryID.Lock()
	mock.calls.GetProductsByNonLeafCategoryID = append(mock.calls.GetProductsByNonLeafCategoryID, callInfo)
	mock.lockGetProductsByNonLeafCategoryID.Unlock()
	return mock.GetProductsByNonLeafCategoryIDFunc(ctx, categoryID, params)
}

// GetProductsByNonLeafCategoryIDCalls gets all the calls that were made to GetProductsByNonLeafCategoryID.
// Check the length with:
//
//	len(mockedProductService.GetProductsByNonLeafCategoryIDCalls())
func (mock *ProductServiceMock) GetProductsByNonLeafCategoryIDCalls() []struct {
	Ctx        context.Context
	CategoryID int32
	Params     anor.GetProductsByCategoryParams
} {
	var calls []struct {
		Ctx        context.Context
		CategoryID int32
		Params     anor.GetProductsByCategoryParams
	}
	mock.lockGetProductsByNonLeafCategoryID.RLock()
	calls = mock.calls.GetProductsByNonLeafCategoryID
	mock.lockGetProductsByNonLeafCategoryID.RUnlock()
	return calls
}

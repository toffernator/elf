package service_test

import (
	"context"
	"elf/internal/core"
	"elf/internal/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock WishlistStore
type MockWishlistStore struct {
	mock.Mock
}

func (s *MockWishlistStore) Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error) {
	return core.Wishlist{}, errors.New("Not yet implemented")
}

func (s *MockWishlistStore) Read(ctx context.Context, id int64) (core.Wishlist, error) {
	return core.Wishlist{}, errors.New("Not yet implemented")
}

func (s *MockWishlistStore) ReadBy(ctx context.Context, p core.WishlistReadByParams) ([]core.Wishlist, error) {
	return make([]core.Wishlist, 0), errors.New("Not yet implemented")
}

func (s *MockWishlistStore) IsOwnedBy(ctx context.Context, wishlistId int64, userId int64) (bool, error) {
	args := s.Called(ctx, wishlistId, userId)
	return args.Bool(0), args.Error(1)
}

func (s *MockWishlistStore) Update(ctx context.Context, p core.WishlistUpdateParams) (core.Wishlist, error) {
	return core.Wishlist{}, errors.New("Not yet implemented")
}

// Mock ProductStore
type MockProductStore struct {
	mock.Mock
}

func (s *MockProductStore) Create(ctx context.Context, p core.ProductCreateParams) (int64, error) {
	args := s.Called(ctx, p)
	return int64(args.Int(0)), args.Error(1)
}

// Write a table drive test for AddProduct (Adding a product to an owned wishlist, adding a product to an unowned wishlist)
func TestWishlistAddProductGivenUserOwnsTheWishlist(t *testing.T) {
	tests := map[string]struct {
		wishlistId          int64
		userId              int64
		productCreateParams core.ProductCreateParams
		expectedErr         error
	}{
		"AddProduct a product to a wishlist the user owns": {
			wishlistId: 1,
			userId:     1,
			productCreateParams: core.ProductCreateParams{
				BelongsToId: 1,
				Name:        "a test product",
				Url:         "https://www.example.com",
				Price:       100,
				Currency:    core.CurrencyEur,
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wishlistStore := &MockWishlistStore{}
			wishlistStore.On("IsOwnedBy", context.TODO(), tt.userId, tt.wishlistId).Return(true, nil)

			productStore := &MockProductStore{}
			productStore.On("Create", context.TODO(), tt.productCreateParams).Return(1, nil)

			wishlists := service.NewWishlistService(wishlistStore, productStore)

			err := wishlists.AddProduct(context.TODO(), tt.wishlistId, tt.userId, tt.productCreateParams)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWishlistAddProductGivenUserDoesNotOwnTheWishlist(t *testing.T) {
	tests := map[string]struct {
		wishlistId          int64
		userId              int64
		productCreateParams core.ProductCreateParams
		expectedErr         error
	}{
		"AddProduct a product to a wishlist the user owns": {
			wishlistId: 1,
			userId:     1,
			productCreateParams: core.ProductCreateParams{
				BelongsToId: 1,
				Name:        "a test product",
				Url:         "https://www.example.com",
				Price:       100,
				Currency:    core.CurrencyEur,
			},
			expectedErr: core.UnauthorizedError{
				Resource: "wishlist",
				Action:   "add a product",
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wishlistStore := &MockWishlistStore{}
			wishlistStore.On("IsOwnedBy", context.TODO(), tt.userId, tt.wishlistId).Return(false, nil)

			productStore := &MockProductStore{}
			productStore.On("Create", context.TODO(), tt.productCreateParams).Return(1, nil)

			wishlists := service.NewWishlistService(wishlistStore, productStore)

			err := wishlists.AddProduct(context.TODO(), tt.wishlistId, tt.userId, tt.productCreateParams)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

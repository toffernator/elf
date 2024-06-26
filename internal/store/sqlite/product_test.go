package sqlite_test

import (
	"context"
	"elf/internal/core"
	"elf/internal/store/sqlite"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: Tests that are expected to fail

func TestProductCreate(t *testing.T) {
	tests := map[string]struct {
		// The "expected" wishlist is derived from the params
		params      core.ProductCreateParams
		expectedErr error
	}{
		"Create a product": {
			params: core.ProductCreateParams{
				BelongsToId: 1,
				Name:        "iPad",
				Url:         "https://www.apple.com/ipad-10.9/",
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := setupSqlite(name)
			products := sqlite.NewProductStore(db)
			t.Cleanup(func() { db.Close() })

			_, err := products.Create(context.TODO(), tt.params)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, err)
		})
	}
}

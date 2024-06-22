package sqlite_test

import (
	"context"
	"elf/internal/sqlite"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var users = sqlite.NewUserStore(db)

var userReadByIdTest = []struct {
	input              int64
	expectedName       string
	expectedOwnerId    int64
	expectedProductLen int
}{
	{
		input:        1,
		expectedName: "test user 1",
	},
}

func TestUserReadById(t *testing.T) {
	seed()
	for _, tt := range userReadByIdTest {
		t.Run(fmt.Sprintf("Read User %d", tt.input), func(t *testing.T) {
			actual, err := users.Read(context.Background(), tt.input)
			if err != nil {
				t.Errorf("%s failed with error: %v", t.Name(), err)
				t.FailNow()
			}

			assert.Equal(t, tt.input, actual.Id)
			assert.Equal(t, tt.expectedName, actual.Name)
		})
	}
}

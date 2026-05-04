package ynab

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetUser(t *testing.T) {
	t.Run("returns user on success", func(t *testing.T) {
		fixture := `{"data":{"user":{"id":"` + testID4 + `"}}}`
		client, _ := newTestClient(fixture, 200)

		user, err := client.GetUser(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := uuid.MustParse(testID4)
		if user.ID != want {
			t.Errorf("got ID %v, want %v", user.ID, want)
		}
	})
}

package ynab

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetUser(t *testing.T) {
	t.Run("returns user on success", func(t *testing.T) {
		fixture := `{"data":{"user":{"id":"123e4567-e89b-12d3-a456-426614174000"}}}`
		client, _ := newTestClient(fixture, 200)

		user, err := client.GetUser(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if user.ID != want {
			t.Errorf("got ID %v, want %v", user.ID, want)
		}
	})
}

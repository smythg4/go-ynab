package ynab

import (
	"context"

	"github.com/google/uuid"
)

type userData struct {
	Data struct {
		User `json:"user"`
	} `json:"data"`
}

type User struct {
	ID uuid.UUID `json:"id"`
}

// GET Methods using user
func (c *Client) GetUser(ctx context.Context) (*User, error) {
	var result userData
	if err := c.get(ctx, "user", nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.User, nil
}

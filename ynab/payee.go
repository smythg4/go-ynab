package ynab

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type payeeData struct {
	Data struct {
		Payee           Payee `json:"payee"`
		ServerKnowledge int64 `json:"server_knowledge"`
	} `json:"data"`
}

type payeesData struct {
	Data struct {
		Payees          []Payee `json:"payees"`
		ServerKnowledge int64   `json:"server_knowledge"`
	} `json:"data"`
}

// Payee represents a payee that can be associated with transactions.
type Payee struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	TransferAccountID *uuid.UUID `json:"transfer_account_id"`
	Deleted           bool       `json:"deleted"`
}

type payeeLocationData struct {
	Data struct {
		PayeeLocation PayeeLocation `json:"payee_location"`
	} `json:"data"`
}

type payeeLocationsData struct {
	Data struct {
		PayeeLocations []PayeeLocation `json:"payee_locations"`
	} `json:"data"`
}

// PayeeLocation represents a geographic location associated with a payee.
type PayeeLocation struct {
	ID        uuid.UUID `json:"id"`
	PayeeID   uuid.UUID `json:"payee_id"`
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
	Deleted   bool      `json:"deleted"`
}

// GET Methods using payees

func (c *Client) GetPayees(ctx context.Context, planID uuid.UUID, params *ListParams) ([]Payee, int64, error) {
	var result payeesData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/payees", planID), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Payees, result.Data.ServerKnowledge, nil
}

func (c *Client) GetPayee(ctx context.Context, planID, payeeID uuid.UUID) (*Payee, error) {
	var result payeeData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/payees/%s", planID, payeeID), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Payee, nil
}

func (c *Client) GetPayeeLocations(ctx context.Context, planID uuid.UUID) ([]PayeeLocation, error) {
	var result payeeLocationsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/payee_locations", planID), nil, &result); err != nil {
		return nil, err
	}
	return result.Data.PayeeLocations, nil
}

func (c *Client) GetPayeeLocationsByPayee(ctx context.Context, planID, payeeID uuid.UUID) ([]PayeeLocation, error) {
	var result payeeLocationsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/payees/%s/payee_locations", planID, payeeID), nil, &result); err != nil {
		return nil, err
	}
	return result.Data.PayeeLocations, nil
}

func (c *Client) GetPayeeLocation(ctx context.Context, planID, locationID uuid.UUID) (*PayeeLocation, error) {
	var result payeeLocationData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/payee_locations/%s", planID, locationID), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.PayeeLocation, nil
}

// POST Methods and infrastructure using payees

// PostPayee is the request body for creating or updating a payee.
type PostPayee struct {
	Name string `json:"name"`
}

type postPayeeWrapper struct {
	Payee PostPayee `json:"payee"`
}

func (c *Client) CreatePayee(ctx context.Context, planID uuid.UUID, pp PostPayee) (*Payee, int64, error) {
	var result payeeData
	err := c.post(ctx, fmt.Sprintf("plans/%s/payees", planID), postPayeeWrapper{pp}, &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.Payee, result.Data.ServerKnowledge, nil
}

// PATCH Methods and infrastructure using payees

func (c *Client) UpdatePayee(ctx context.Context, planID, payeeID uuid.UUID, pp PostPayee) (*Payee, int64, error) {
	var result payeeData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/payees/%s", planID, payeeID), postPayeeWrapper{pp}, &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.Payee, result.Data.ServerKnowledge, nil
}

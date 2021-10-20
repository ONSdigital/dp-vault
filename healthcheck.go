package vault

import (
	"context"
	"errors"

	health "github.com/ONSdigital/dp-healthcheck/v2/healthcheck"
)

// ServiceName vault
const ServiceName = "vault"

// MsgHealthy Check message returned when vault is healthy
const MsgHealthy = "vault is healthy"

// Error definitions
var (
	ErrNotInitialised = errors.New("vault not initialised")
)

// Healthcheck determines the state of vault
func (c *Client) Healthcheck() (string, error) {
	resp, err := c.client.Health()
	if err != nil {
		return "vault", err
	}

	if !resp.Initialized {
		return "vault", ErrNotInitialised
	}

	return "", nil
}

// Checker performs a check health of Vault and updates the provided CheckState accordingly
func (c *Client) Checker(ctx context.Context, state *health.CheckState) error {
	_, err := c.Healthcheck()
	if err != nil {
		state.Update(health.StatusCritical, err.Error(), 0)
		return nil
	}
	state.Update(health.StatusOK, MsgHealthy, 0)
	return nil
}

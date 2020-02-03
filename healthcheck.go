package vault

import (
	"context"
	"errors"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out ./mock/check_state.go -pkg mock . CheckState

// CheckState interface corresponds to the healthcheck CheckState structure
type CheckState interface {
	Update(status, message string, statusCode int) error
}

// ServiceName vault
const ServiceName = "vault"

// MsgHealthy Check message returned when vault is healthy
const MsgHealthy = "vault is healthy"

// Error definitions
var (
	ErrNotInitialised = errors.New("vault not initialised")
)

// minTime is the oldest time for Check structure.
var minTime = time.Unix(0, 0)

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
func (c *Client) Checker(ctx context.Context, state CheckState) error {
	_, err := c.Healthcheck()
	if err != nil {
		state.Update(health.StatusCritical, err.Error(), 0)
		return nil
	}
	state.Update(health.StatusOK, MsgHealthy, 0)
	return nil
}

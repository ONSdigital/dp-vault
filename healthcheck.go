package vault

import (
	"context"
	"errors"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

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

// Checker performs a Vault health check and return it inside a Check structure
func (c *Client) Checker(ctx *context.Context) (*health.Check, error) {
	_, err := c.Healthcheck()
	if err != nil {
		return getCheck(ctx, health.StatusCritical, err.Error()), err
	}
	return getCheck(ctx, health.StatusOK, MsgHealthy), nil
}

// getCheck reates a Check structure and populates it according to status and message provided
func getCheck(ctx *context.Context, status, message string) *health.Check {

	currentTime := time.Now().UTC()

	check := &health.Check{
		Name:        ServiceName,
		Status:      status,
		Message:     message,
		LastChecked: currentTime,
		LastSuccess: minTime,
		LastFailure: minTime,
	}

	switch status {
	case health.StatusOK:
		check.StatusCode = 200
		check.LastSuccess = currentTime
	case health.StatusWarning:
		check.StatusCode = 429
		check.LastFailure = currentTime
	default:
		check.StatusCode = 500
		check.Status = health.StatusCritical
		check.LastFailure = currentTime
	}

	return check
}

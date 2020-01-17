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
	currentTime := time.Now().UTC()
	c.Check.LastChecked = &currentTime
	if err != nil {
		c.Check.LastFailure = &currentTime
		c.Check.Status = health.StatusCritical
		c.Check.Message = err.Error()
		return c.Check, err
	}
	c.Check.LastSuccess = &currentTime
	c.Check.Status = health.StatusOK
	c.Check.Message = MsgHealthy
	return c.Check, nil
}

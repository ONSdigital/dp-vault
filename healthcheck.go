package vault

import (
	"context"
	"errors"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

// ServiceName vault
const ServiceName = "vault"

// StatusDescription : Map of descriptions by status
var StatusDescription = map[string]string{
	health.StatusOK:       "Everything is ok",
	health.StatusWarning:  "Things are degraded, but at least partially functioning",
	health.StatusCritical: "The checked functionality is unavailable or non-functioning",
}

// Error definitions
var (
	ErrNotInitialised = errors.New("vault not initialised")
)

// minTime : Oldest time for Check structure.
var minTime = time.Unix(0, 0)

// Healthcheck determines the state of vault
func (c *VaultClient) Healthcheck() (string, error) {
	resp, err := c.client.Health()
	if err != nil {
		return "vault", err
	}

	if !resp.Initialized {
		return "vault", ErrNotInitialised
	}

	return "", nil
}

// Checker : Check health of Vault and return it inside a Check structure
func (c *VaultClient) Checker(ctx *context.Context) (*health.Check, error) {
	_, err := c.Healthcheck()
	if err != nil {
		return getCheck(ctx, 500), err
	}
	return getCheck(ctx, 200), nil
}

// getCheck : Create a Check structure and populate it according to the code
func getCheck(ctx *context.Context, code int) *health.Check {

	currentTime := time.Now().UTC()

	check := &health.Check{
		Name:        ServiceName,
		StatusCode:  code,
		LastChecked: currentTime,
		LastSuccess: minTime,
		LastFailure: minTime,
	}

	switch code {
	case 200:
		check.Message = StatusDescription[health.StatusOK]
		check.Status = health.StatusOK
		check.LastSuccess = currentTime
	case 429:
		check.Message = StatusDescription[health.StatusWarning]
		check.Status = health.StatusWarning
		check.LastFailure = currentTime
	default:
		check.Message = StatusDescription[health.StatusCritical]
		check.Status = health.StatusCritical
		check.LastFailure = currentTime
	}

	return check
}

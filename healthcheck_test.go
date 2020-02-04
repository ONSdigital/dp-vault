package vault_test

import (
	"context"
	"errors"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/dp-vault/mock"
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/smartystreets/goconvey/convey"
)

// Error definitions for testing
var (
	ErrNilRequest = errors.New("nil request created")
)

var respNotInitialised = &vaultapi.HealthResponse{
	Initialized: false,
}

var healthSuccess = func() (*vaultapi.HealthResponse, error) {
	return &vaultapi.HealthResponse{
		Initialized: true,
	}, nil
}

var healthNotInitialised = func() (*vaultapi.HealthResponse, error) {
	return &vaultapi.HealthResponse{
		Initialized: false,
	}, nil
}

var healthError = func() (*vaultapi.HealthResponse, error) {
	return nil, ErrNilRequest
}

func TestVaultHealthy(t *testing.T) {
	Convey("Given that Vault is healthy", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthSuccess,
		}
		cli := vault.CreateClientWithAPIClient(apiCli)

		// CheckState for test validation
		checkState := health.NewCheckState(vault.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, vault.MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestVaultNotInitialised(t *testing.T) {
	Convey("Given that Vault has not been initialised", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthNotInitialised,
		}
		cli := vault.CreateClientWithAPIClient(apiCli)

		// CheckState for test validation
		checkState := health.NewCheckState(vault.ServiceName)

		Convey("Checker updates the CheckState to a Critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, vault.ErrNotInitialised.Error())
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestVaultAPIError(t *testing.T) {
	Convey("Given that Vault API Health returns an error", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthError,
		}
		cli := vault.CreateClientWithAPIClient(apiCli)

		// CheckState for test validation
		checkState := health.NewCheckState(vault.ServiceName)

		Convey("Checker updates the CheckState to a Critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, ErrNilRequest.Error())
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

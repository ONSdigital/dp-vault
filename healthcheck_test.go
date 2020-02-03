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

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusOK)
			So(updateCalls[0].Message, ShouldEqual, vault.MsgHealthy)
			So(updateCalls[0].StatusCode, ShouldEqual, 0)
		})
	})
}

func TestVaultNotInitialised(t *testing.T) {
	Convey("Given that Vault has not been initialised", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthNotInitialised,
		}
		cli := vault.CreateClientWithAPIClient(apiCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a Critical state with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, vault.ErrNotInitialised.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 0)
		})
	})
}

func TestVaultAPIError(t *testing.T) {
	Convey("Given that Vault API Health returns an error", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthError,
		}
		cli := vault.CreateClientWithAPIClient(apiCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a Critical state with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, ErrNilRequest.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 0)
		})
	})
}

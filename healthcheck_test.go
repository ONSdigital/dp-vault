package vault_test

import (
	"errors"
	"testing"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/dp-vault/mock"
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/smartystreets/goconvey/convey"
)

// initial check that should be created by client constructor
var expectedInitialCheck = &health.Check{
	Name: vault.ServiceName,
}

// create a successful check without lastFailed value
func createSuccessfulCheck(t time.Time, msg string) health.Check {
	return health.Check{
		Name:        vault.ServiceName,
		LastSuccess: &t,
		LastChecked: &t,
		Status:      health.StatusOK,
		Message:     msg,
	}
}

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
		So(cli.Check, ShouldResemble, expectedInitialCheck)

		Convey("Checker returns a successful Check struct", func() {
			validateSuccessfulCheck(cli)
			So(cli.Check.LastFailure, ShouldBeNil)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
		})
	})
}

func TestVaultNotInitialised(t *testing.T) {
	Convey("Given that Vault has not been initilised", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthNotInitialised,
		}
		cli := vault.CreateClientWithAPIClient(apiCli)
		So(cli.Check, ShouldResemble, expectedInitialCheck)

		Convey("Checker returns a Critical Check struct", func() {
			_, err := validateCriticalCheck(cli, vault.ErrNotInitialised.Error())
			So(err, ShouldEqual, vault.ErrNotInitialised)
			So(cli.Check.LastSuccess, ShouldBeNil)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
		})
	})
}

func TestVaultAPIError(t *testing.T) {
	Convey("Given that Vault API Health returns an error", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthError,
		}
		cli := vault.CreateClientWithAPIClient(apiCli)
		So(cli.Check, ShouldResemble, expectedInitialCheck)

		Convey("Checker returns a Critical Check struct", func() {
			_, err := validateCriticalCheck(cli, ErrNilRequest.Error())
			So(err, ShouldResemble, ErrNilRequest)
			So(cli.Check.LastSuccess, ShouldBeNil)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
		})
	})
}

func validateSuccessfulCheck(cli *vault.Client) (check *health.Check) {
	t0 := time.Now().UTC()
	check, err := cli.Checker(nil)
	t1 := time.Now().UTC()
	So(err, ShouldBeNil)
	So(check.Name, ShouldEqual, vault.ServiceName)
	So(check.Status, ShouldEqual, health.StatusOK)
	So(check.Message, ShouldEqual, vault.MsgHealthy)
	So(*check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(*check.LastSuccess, ShouldHappenOnOrBetween, t0, t1)
	return check
}

func validateWarningCheck(cli *vault.Client, expectedMessage string) (check *health.Check, err error) {
	t0 := time.Now().UTC()
	check, err = cli.Checker(nil)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, vault.ServiceName)
	So(check.Status, ShouldEqual, health.StatusWarning)
	So(check.Message, ShouldEqual, expectedMessage)
	So(*check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(*check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}

func validateCriticalCheck(cli *vault.Client, expectedMessage string) (check *health.Check, err error) {
	t0 := time.Now().UTC()
	check, err = cli.Checker(nil)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, vault.ServiceName)
	So(check.Status, ShouldEqual, health.StatusCritical)
	So(check.Message, ShouldEqual, expectedMessage)
	So(*check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(*check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}

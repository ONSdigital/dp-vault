package vault_test

import (
	"fmt"
	"testing"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/dp-vault/mock"
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/smartystreets/goconvey/convey"
)

var respNotInitialised = &vaultapi.HealthResponse{
	Initialized: false,
}

// 	Initialized                bool   `json:"initialized"`
// 	Sealed                     bool   `json:"sealed"`
// 	Standby                    bool   `json:"standby"`
// 	PerformanceStandby         bool   `json:"performance_standby"`
// 	ReplicationPerformanceMode string `json:"replication_performance_mode"`
// 	ReplicationDRMode          string `json:"replication_dr_mode"`
// 	ServerTimeUTC              int64  `json:"server_time_utc"`
// 	Version                    string `json:"version"`
// 	ClusterName                string `json:"cluster_name,omitempty"`
// 	ClusterID                  string `json:"cluster_id,omitempty"`
// 	LastWAL                    uint64 `json:"last_wal,omitempty"`
// }

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
	return nil, fmt.Errorf("nil request created")
}

// var doOkMissingStatus = func(ctx context.Context, request *http.Request) (*http.Response, error) {
// 	return resp(clusterMissingStatus, 200), nil
// }

// var doUnexpectedCode = func(ctx context.Context, request *http.Request) (*http.Response, error) {
// 	return resp(clusterHealthy, 300), nil
// }

// var doUnreacheable = func(ctx context.Context, request *http.Request) (*http.Response, error) {
// 	return nil, errors.New("unreacheable")
// }

// func resp(body string, code int) *http.Response {
// 	return &http.Response{
// 		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
// 		StatusCode: code,
// 	}
// }

func TestVaultHealthy(t *testing.T) {
	Convey("Given that Vault is healthy", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthSuccess,
		}
		cli := vault.CreateVaultClientWithAPIClient(apiCli)

		Convey("Checker returns a successful Check struct", func() {
			validateSuccessfulCheck(cli)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
		})
	})
}

func TestVaultNotInitialised(t *testing.T) {
	Convey("Given that Vault has not been initilised", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthNotInitialised,
		}
		cli := vault.CreateVaultClientWithAPIClient(apiCli)

		Convey("Checker returns a Critical Check struct", func() {
			_, err := validateCriticalCheck(cli)
			So(err, ShouldEqual, vault.ErrNotInitialised)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
		})
	})
}

func TestVaultAPIError(t *testing.T) {
	Convey("Given that Vault API Health returns an error", t, func() {

		var apiCli = &mock.APIClientMock{
			HealthFunc: healthError,
		}
		cli := vault.CreateVaultClientWithAPIClient(apiCli)

		Convey("Checker returns a Critical Check struct", func() {
			_, err := validateCriticalCheck(cli)
			So(err, ShouldNotBeNil)
			So(len(apiCli.HealthCalls()), ShouldEqual, 1)
		})
	})
}

func validateSuccessfulCheck(cli *vault.VaultClient) (check *health.Check) {
	t0 := time.Now().UTC()
	check, err := cli.Checker(nil)
	t1 := time.Now().UTC()
	So(err, ShouldBeNil)
	So(check.Name, ShouldEqual, vault.ServiceName)
	So(check.Status, ShouldEqual, health.StatusOK)
	So(check.StatusCode, ShouldEqual, 200)
	So(check.Message, ShouldEqual, vault.StatusDescription[health.StatusOK])
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastFailure, ShouldHappenBefore, t0)
	return check
}

func validateWarningCheck(cli *vault.VaultClient) (check *health.Check, err error) {
	t0 := time.Now().UTC()
	check, err = cli.Checker(nil)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, vault.ServiceName)
	So(check.Status, ShouldEqual, health.StatusWarning)
	So(check.StatusCode, ShouldEqual, 429)
	So(check.Message, ShouldEqual, vault.StatusDescription[health.StatusWarning])
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenBefore, t0)
	So(check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}

func validateCriticalCheck(cli *vault.VaultClient) (check *health.Check, err error) {
	t0 := time.Now().UTC()
	check, err = cli.Checker(nil)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, vault.ServiceName)
	So(check.Status, ShouldEqual, health.StatusCritical)
	So(check.StatusCode, ShouldEqual, 500)
	So(check.Message, ShouldEqual, vault.StatusDescription[health.StatusCritical])
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenBefore, t0)
	So(check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}

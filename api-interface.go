package vault

import vaultapi "github.com/hashicorp/vault/api"

//go:generate moq -out ./mock/api-client.go -pkg mock . APIClient

// APIClient is an interface to wrap Vault API Client, which is used by the dp-vault Client in order to interact with Vault, and can be easily mocked for testing
type APIClient interface {
	SetToken(v string)
	Read(path string) (*vaultapi.Secret, error)
	Write(path string, data map[string]interface{}) (*vaultapi.Secret, error)
	Health() (*vaultapi.HealthResponse, error)
}

// APIClientImpl implements the APIClient interface wrapping the real vault API client calls to nested clients (e.g Logical or Sys)
type APIClientImpl struct {
	client *vaultapi.Client
}

// SetToken calls SetToken directly to vault API client
func (api *APIClientImpl) SetToken(v string) {
	api.client.SetToken(v)
}

// Read calls Read(path) from the Logical client in vault API client
func (api *APIClientImpl) Read(path string) (*vaultapi.Secret, error) {
	return api.client.Logical().Read(path)
}

// Write calls Write(path, data) from the Logical client in vault API client
func (api *APIClientImpl) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	return api.client.Logical().Write(path, data)
}

// Health calls Health() from the Sys client in vault API client
func (api *APIClientImpl) Health() (*vaultapi.HealthResponse, error) {
	return api.client.Sys().Health()
}

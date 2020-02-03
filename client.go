package vault

import (
	"encoding/json"
	"errors"

	vaultapi "github.com/hashicorp/vault/api"
)

// Error definitions
var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrVersionNotFound  = errors.New("version not found")
	ErrMetadataNotFound = errors.New("metadata not found")
	ErrDataNotFound     = errors.New("data not found")
	ErrVersionInvalid   = errors.New("version failed to convert to number")
)

// Client Used to read and write secrets from vault using a vault API client wrapper.
type Client struct {
	client APIClient
}

// CreateClient by providing an auth token, vault address and the maximum number of retries for a request
func CreateClient(token, vaultAddress string, retries int) (*Client, error) {
	config := vaultapi.Config{Address: vaultAddress, MaxRetries: retries}
	return CreateClientWithConfig(&config, token)
}

// CreateClientTLS is like the CreateClient function but wraps the HTTP client with TLS
func CreateClientTLS(token, vaultAddress string, retries int, cacert, cert, key string) (*Client, error) {
	config := vaultapi.Config{Address: vaultAddress, MaxRetries: retries}
	config.ConfigureTLS(&vaultapi.TLSConfig{CACert: cacert, ClientCert: cert, ClientKey: key})
	return CreateClientWithConfig(&config, token)
}

// CreateClientWithConfig creates a Client with provided config and token as inputs
func CreateClientWithConfig(config *vaultapi.Config, token string) (*Client, error) {
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	apiClient := &APIClientImpl{client}
	apiClient.SetToken(token)
	return CreateClientWithAPIClient(apiClient), nil
}

// CreateClientWithAPIClient creates a Client with a provided Vault API client as input
func CreateClientWithAPIClient(apiClient APIClient) *Client {
	return &Client{
		client: apiClient,
	}
}

// Read reads a secret from vault. If the token does not have the correct policy this returns an error;
// if the vault server is not reachable, return all the information stored about the secret.
func (c *Client) Read(path string) (map[string]interface{}, error) {
	secret, err := c.client.Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		// If there is no secret and no error return a empty map.
		return make(map[string]interface{}), nil
	}
	return secret.Data, err
}

// ReadKey from vault. Like read but only return a single value from the secret
func (c *Client) ReadKey(path, key string) (string, error) {
	data, err := c.Read(path)
	if err != nil {
		return "", err
	}
	val, ok := data[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return val.(string), nil
}

// Write writes a secret to vault. Returns an error if the token does not have the correct policy or the
// vault server is not reachable. Returns nil when the operation was successful.
func (c *Client) Write(path string, data map[string]interface{}) error {
	_, err := c.client.Write(path, data)
	return err
}

// WriteKey writes a secret value for a specific key to vault.
func (c *Client) WriteKey(path, key, value string) error {
	data := make(map[string]interface{})
	data[key] = value
	return c.Write(path, data)
}

// VRead reads a versioned secret from vault - cf Read, above -
// returns the secret (map) and the version
func (c *Client) VRead(path string) (map[string]interface{}, int64, error) {
	secret, err := c.Read(path)
	if err != nil {
		return nil, -1, err
	}
	if len(secret) == 0 {
		// if there is no secret and no error return a empty map
		return secret, -1, nil
	}
	metadata, ok := secret["metadata"]
	if !ok {
		return nil, -1, ErrMetadataNotFound
	}
	verObj, ok := metadata.(map[string]interface{})["version"]
	if !ok {
		return nil, -1, ErrVersionNotFound
	}
	ver, err := verObj.(json.Number).Int64()
	if err != nil {
		return nil, -1, ErrVersionInvalid
	}
	data, ok := secret["data"]
	if !ok {
		return nil, -1, ErrDataNotFound
	}
	return data.(map[string]interface{}), ver, err
}

// VReadKey - cf Read but for versioned secret - return the value of the key and the version
func (c *Client) VReadKey(path, key string) (string, int64, error) {
	data, ver, err := c.VRead(path)
	if err != nil {
		return "", -1, err
	}
	val, ok := data[key]
	if !ok {
		return "", -1, ErrKeyNotFound
	}
	return val.(string), ver, nil
}

// VWriteKey creates a data map, with a data field containing the key-value, and writes it to vault
func (c *Client) VWriteKey(path, key, value string) error {
	data := map[string]interface{}{
		"data": map[string]interface{}{
			key: value,
		},
	}
	return c.Write(path, data)
}

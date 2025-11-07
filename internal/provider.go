package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type APIClient struct {
	Endpoint   string
	User       string
	Password   string
	AuthToken  string
	HTTPClient http.Client
}

// Provider returns the Terraform provider schema and resources.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAZUH_ENDPOINT", nil),
				Description: "Full URL to Wazuh API endpoint (e.g. https://wazuh.example.com:55000).",
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("WAZUH_USER", nil),
				Description: "Wazuh API username.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("WAZUH_PASSWORD", nil),
				Description: "Wazuh API password.",
			},
			"skip_ssl_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WAZUH_SKIP_SSL_VERIFY", false),
				Description: "Skip SSL certificate verification.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wazuh_group":               resourceGroup(),
			"wazuh_group_configuration": resourceGroupConfiguration(),
			"wazuh_active_response":     resourceActiveResponse(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	endpoint := d.Get("endpoint").(string)
	user := d.Get("user").(string)
	password := d.Get("password").(string)
	skipSSL := d.Get("skip_ssl_verify").(bool)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSL,
		},
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	client := &APIClient{
		Endpoint:   endpoint,
		User:       user,
		Password:   password,
		HTTPClient: *httpClient,
	}

	token, err := client.authenticate()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	client.AuthToken = token

	return client, diags
}

// Authenticate with basicAuth and obtain JWT token
func (c *APIClient) authenticate() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/security/user/authenticate", c.Endpoint), nil)
	if err != nil {
		return "", fmt.Errorf("failed to build auth request: %w", err)
	}
	req.SetBasicAuth(c.User, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform authentication request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("authentication failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
		Error int `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %w", err)
	}
	if result.Error != 0 || result.Data.Token == "" {
		return "", fmt.Errorf("authentication failed: %s", string(body))
	}

	return result.Data.Token, nil
}

// doRequest sends authenticated requests to Wazuh API
func (c *APIClient) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.Endpoint, path), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

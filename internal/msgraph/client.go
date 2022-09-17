package msgraph

import (
	"fmt"
	"net/http"
)

// ClientConfiguration represents the vinyldns client configuration.
type ClientConfiguration struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	Scope        string
	GrantType    string
	AuthHost     string
	GraphHost    string
	UserAgent    string
}

func defaultUA() string {
	return fmt.Sprintf("go-msgraph-application/%s", "0.0.1")
}

func NewClient(config ClientConfiguration) *Client {
	if config.UserAgent == "" {
		config.UserAgent = defaultUA()
	}

	return &Client{
		HTTPClient:   &http.Client{},
		UserAgent:    config.UserAgent,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TenantID:     config.TenantID,
		Scope:        config.Scope,
		GrantType:    config.GrantType,
		AuthHost:     config.AuthHost,
		GraphHost:    config.GraphHost,
	}
}

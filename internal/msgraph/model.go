package msgraph

import "net/http"

type Auth struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	Scope        string
	GrantType    string
	AuthHost     string
	GraphHost    string
	AppID        string

	AuthResult
}

type AuthResult struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	ExtExpiresIn int64  `json:"ext_expires_in"`
	TokenType    string `json:"token_type"`
}

type Applications struct {
	Odata_context string        `json:"@odata.context"`
	Value         []Application `json:"value"`
}

type Application struct {
	AppID       string         `json:"appId"`
	DisplayName string         `json:"displayName"`
	ID          string         `json:"id"`
	Web         ApplicationWeb `json:"web"`
}

type ApplicationWeb struct {
	ImplicitGrantSettings struct {
		EnableAccessTokenIssuance bool `json:"enableAccessTokenIssuance"`
		EnableIDTokenIssuance     bool `json:"enableIdTokenIssuance"`
	} `json:"implicitGrantSettings"`
	// RedirectURISettings []RedirectURISettings `json:"redirectUriSettings"`
	RedirectUris []string `json:"redirectUris"`
}

type RedirectURISettings struct {
	Index interface{} `json:"index"`
	URI   string      `json:"uri"`
}

type Client struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	HTTPClient   *http.Client
	UserAgent    string `json:"user_agent"`
	ClientID     string
	ClientSecret string
	TenantID     string
	Scope        string
	GrantType    string
	AuthHost     string
	GraphHost    string
}

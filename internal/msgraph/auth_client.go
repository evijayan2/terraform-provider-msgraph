package msgraph

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GraphAccess() {

	url := fmt.Sprintf("%s/%s/oauth2/v2.0/token?", c.AuthHost, c.TenantID)
	method := "POST"

	requestBody := fmt.Sprintf("client_id=%s&scope=%s&client_secret=%s&grant_type=%s", c.ClientID, c.Scope, c.ClientSecret, c.GrantType)
	payload := strings.NewReader(requestBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Host", "login.microsoftonline.com")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	authResult := &AuthResult{}
	json.NewDecoder(res.Body).Decode(authResult)
	c.AccessToken = authResult.AccessToken
	c.TokenType = authResult.TokenType
}

package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (c *Client) GetApplication(appId string) *Application {

	url := fmt.Sprintf("%s/v1.0/applications?$count=true&$select=id,appId,displayName,web&$filter=appId%%20eq%%20'%s'", c.GraphHost, appId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Printf("got error %v", err)
		return nil
	}

	auth_header := fmt.Sprintf("%s %s", c.TokenType, c.AccessToken)
	req.Header.Add("Authorization", auth_header)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("got http error %v", err)
		return nil
	}
	defer res.Body.Close()

	applications := &Applications{}
	err = json.NewDecoder(res.Body).Decode(applications)
	if err != nil {
		fmt.Printf("got json marshal error %v", err)
		return nil
	}

	if len(applications.Value) <= 0 {
		tflog.Error(context.Background(), "no application found.")
	}

	return &applications.Value[0]
}

func (c *Client) PatchWebAddRedirectURI(application Application, redirectUri string) error {
	application.Web.RedirectUris = append(application.Web.RedirectUris, redirectUri)

	url := fmt.Sprintf("%s/v1.0/applications/%s", c.GraphHost, application.ID)
	method := "PATCH"

	payload, err := json.Marshal(application)
	if err != nil {
		return fmt.Errorf("got http error %v", err)
	}
	tflog.Trace(context.Background(), fmt.Sprintf("payload %s\r\n", string(payload)))

	client := c.HTTPClient
	req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))

	if err != nil {
		return fmt.Errorf("got http error %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	auth_header := fmt.Sprintf("%s %s", c.TokenType, c.AccessToken)
	req.Header.Add("Authorization", auth_header)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("got http error %v", err)
	}
	defer res.Body.Close()

	tflog.Trace(context.Background(), fmt.Sprintf("%d: %s for %s\r\n", res.StatusCode, res.Status, redirectUri))
	if res.StatusCode != 204 {
		tflog.Trace(context.Background(), "patch add request not processed")
		return fmt.Errorf("patch add request not processed")
	} else {
		tflog.Info(context.Background(), fmt.Sprintf("%s added successfully\r\n", redirectUri))
	}

	return nil
}

func (c *Client) CheckRedirectURI(application Application, redirectUri string) bool {

	for i := 0; i < len(application.Web.RedirectUris); i++ {
		if application.Web.RedirectUris[i] == redirectUri {
			return true
		}
	}
	return false
}

func (c *Client) PatchWebRemoveRedirectURI(application Application, redirectUri string) error {

	newRedirectUris := make([]string, 0)
	for i := 0; i < len(application.Web.RedirectUris); i++ {
		if application.Web.RedirectUris[i] != redirectUri {
			newRedirectUris = append(newRedirectUris, application.Web.RedirectUris[i])
		}
	}

	if len(newRedirectUris) == len(application.Web.RedirectUris) {
		return fmt.Errorf("sorry nothing to remove")
	}

	application.Web.RedirectUris = newRedirectUris

	url := fmt.Sprintf("%s/v1.0/applications/%s", c.GraphHost, application.ID)
	method := "PATCH"

	payload, err := json.Marshal(application)
	if err != nil {
		return fmt.Errorf("got json marshal error %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))

	if err != nil {
		return fmt.Errorf("got http error %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	auth_header := fmt.Sprintf("%s %s", c.TokenType, c.AccessToken)
	req.Header.Add("Authorization", auth_header)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("got http error %v", err)
	}
	defer res.Body.Close()

	tflog.Trace(context.Background(), fmt.Sprintf("%d: %s for %s\r\n", res.StatusCode, res.Status, redirectUri))
	if res.StatusCode != 204 {
		tflog.Trace(context.Background(), "patch delete request not processed")
		return fmt.Errorf("patch delete request not processed")
	} else {
		tflog.Info(context.Background(), fmt.Sprintf("%s removed successfully\r\n", redirectUri))
	}

	return nil
}

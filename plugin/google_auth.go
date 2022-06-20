package gaccauth

import (
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleAuth struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	RedirectURL    string `json:"redirect_url"`
	FetchGroups    bool   `json:"fetch_groups"`
	ServiceAccount string `json:"service_acc_key"`
	DelegationUser string `json:"delegation_user"`
}

func (c *googleAuth) oauth2Config() *oauth2.Config {
	oauthRedirectURL := c.RedirectURL
	if len(strings.TrimSpace(oauthRedirectURL)) == 0 {
		oauthRedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	}

	config := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  oauthRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}

	return config
}

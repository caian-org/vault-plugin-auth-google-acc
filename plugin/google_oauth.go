package gaccauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleOAuth struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	RedirectURL    string `json:"redirect_url"`
	FetchGroups    bool   `json:"fetch_groups"`
	ServiceAccount string `json:"service_acc_key"`
	DelegationUser string `json:"delegation_user"`
}

func (c *googleOAuth) build() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}
}

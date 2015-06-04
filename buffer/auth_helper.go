package buffer

import (
	"net/http"

	"golang.org/x/oauth2"
)

func GetOauth2Client(accessToken string) *http.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken})

	return oauth2.NewClient(oauth2.NoContext, tokenSource)
}

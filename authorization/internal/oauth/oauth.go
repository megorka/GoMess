package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

var GoogleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("OAUTH_CLIENT"),
	ClientSecret: os.Getenv("OAUTH_SECRET"),
	Endpoint:     google.Endpoint,
	RedirectURL:  os.Getenv("OAUTH_URL"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
}

//var YandexOauthConfig = &oauth2.Config{
//	ClientID:     "",
//	ClientSecret: "",
//	Endpoint:     yandex.Endpoint,
//	RedirectURL:  "",
//	Scopes:       nil,
//}

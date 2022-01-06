package module

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
)

type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var conf *oauth2.Config
var c Credentials
var state string
var scopes []string

func Setup() string {
	state = RandToken()
	GetDetailsFromJson()
	return getLoginURL(state)
}

func GetDetailsFromJson() {
	file, err := ioutil.ReadFile("./google.json")
	if err != nil {
		glog.Fatalf("[Gin-OAuth] File error: %v\n", err)
	}
	json.Unmarshal(file, &c)

	var scopes = append(scopes, "openid", "profile", "email")
	conf = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  "http://sso.scd.edu.om:1234/sso/v1/callback",
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}
func GetStateCode() string {
	return state
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func RandToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)

}
func getCode(code string) {

}

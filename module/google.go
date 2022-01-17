package module

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"sso.scd.edu.om/structs"
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
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	file, err := ioutil.ReadFile(exPath + "/google.json")
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

func AuthHandler(c *gin.Context, code string, state string) (structs.UserSession, structs.UserLogin, error) {
	session := sessions.Default(c)
	token, err := conf.Exchange(c, code)
	if err != nil {
		return structs.UserSession{}, structs.UserLogin{}, err
	}
	tokenID, _ := token.Extra("id_token").(string)
	scope, _ := token.Extra("scope").(string)
	session.Set("TokenID", tokenID)

	client := conf.Client(c, token)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return structs.UserSession{}, structs.UserLogin{}, err
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)

	users := structs.UserLogin{}
	userSessions := structs.UserSession{}
	if err = json.Unmarshal(data, &users); err != nil {
		return structs.UserSession{}, structs.UserLogin{}, err
	}
	session.Set("UserEmail", users.ClaimUserEmail)
	userSessions.UserAgent = c.Request.UserAgent()
	userSessions.UserIP = c.ClientIP()
	userSessions.UserState = state
	userSessions.GoogleCode = code
	userSessions.GoogleAccessToken = token.AccessToken
	userSessions.GoogleExpireIn = token.Expiry.Unix()
	userSessions.GoogleTokenType = token.Type()
	userSessions.GoogleIDToken = tokenID
	userSessions.GoogleScope = scope
	//userSessions.GoogleRefreshToken = token.RefreshToken
	fmt.Println(token.RefreshToken)
	session.Save()
	return userSessions, users, nil
}

package route

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"sso.scd.edu.om/handler"
	"sso.scd.edu.om/module"
)

//SetupRoutes : all the routes are defined here
func SetupRoutes(db *gorm.DB) {
	httpRouter := gin.Default()
	httpRouter.Use(module.SessionConnection())
	httpRouter.Use(handler.LoggerToFile("sso.scd.edu.om"))
	gin.SetMode(gin.ReleaseMode)
	httpRouter.Use(cors.Default())
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	httpRouter.LoadHTMLGlob(exPath + "/web/pages/*")
	httpRouter.Static("/static", exPath+"/web/assets")
	httpRouter.GET("/sso/v1/login", func(c *gin.Context) {
		urlRedirect, _ := c.GetQuery("redirect")
		urlC, urlErr := url.ParseRequestURI(urlRedirect)
		if urlErr != nil {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-106",
				"ErrorMessage": "Malformed or missing header.",
			})
		}
		if !strings.Contains(urlC.Host, "scd.edu.om") {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-106",
				"ErrorMessage": "Malformed or missing header.",
			})
		}
		if len(urlRedirect) > 0 {
			isActive := module.CreateSession(c, urlRedirect)
			if isActive == false {
				c.HTML(http.StatusOK, "login.html", gin.H{
					"GoogleLink": module.Setup(),
				})
			} else {
				location := url.URL{Path: urlRedirect}
				c.Redirect(http.StatusFound, location.RequestURI())
			}
		}

	})
	//session extension
	httpRouter.GET("/sso/v1/refresh", func(c *gin.Context) {
		urlFrom := c.Request.Header.Get("Referer")
		fmt.Println(urlFrom)
		extend, redirectUrl, extendErr := module.ExtendUserSession(c, db, urlFrom)
		if extend == false && extendErr != nil {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-107",
				"ErrorMessage": "Cant extend for user, due to " + extendErr.Error(),
			})
			return
		} else {
			c.HTML(http.StatusOK, "refresh.html", gin.H{
				"ErrorTitle":   "Session Extended",
				"ErrorMessage": "Your session has been extended successfully",
				"RedirectUrl":  redirectUrl,
			})
			return
		}

	})

	//google callback path
	httpRouter.GET("/sso/v1/callback", func(c *gin.Context) {
		//Fetch url parameters
		UrlState, stateErr := c.GetQuery("state")
		if stateErr == false {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-101",
				"ErrorMessage": "Cant verify user, please re-login. If the issue presist contact IT office.",
			})
			return
		}
		UrlCode, codeErr := c.GetQuery("code")
		if codeErr == false {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-102",
				"ErrorMessage": "Cant verify user, please re-login. If the issue presist contact IT office.",
			})
			return
		}
		//check state with google module
		if strings.Compare(UrlState, module.GetStateCode()) != 0 {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-103",
				"ErrorMessage": "Cant verify user, please re-login. If the issue presist contact IT office.",
			})
			return
		}
		UserSession, UserData, authErr := module.AuthHandler(c, UrlCode, UrlState)
		if authErr != nil {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-104",
				"ErrorMessage": "Cant login user, please re-login. If the issue presist contact IT office.",
			})
			return
		}
		//Create full Session and add to DB
		urlRedirect, dberr, _ := module.LoginUserIntoDB(c, UserData, UserSession, db)
		if dberr == false {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-105",
				"ErrorMessage": "DB Error check logs, please re-login. If the issue presist contact IT office.",
			})
			return
		}
		location := url.URL{Path: urlRedirect}
		c.Redirect(http.StatusFound, location.RequestURI())
	})
	httpRouter.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "pagenotfound.html", gin.H{})
	})
	//httpRouter.RunTLS(":1234", "certs/server.crt", "certs/server.key")
	httpRouter.Run(":80")

}

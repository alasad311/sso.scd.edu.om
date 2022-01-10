package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sso.scd.edu.om/handler"
	"sso.scd.edu.om/module"
	"strings"
)

//SetupRoutes : all the routes are defined here
func SetupRoutes() {
	httpRouter := gin.Default()
	httpRouter.Use(module.SessionConnection())
	httpRouter.Use(handler.LoggerToFile("sso.scd.edu.om"))
	gin.SetMode(gin.ReleaseMode)
	httpRouter.Use(cors.Default())
	httpRouter.LoadHTMLGlob("web/pages/*")
	httpRouter.Static("/static", "./web/assets")
	httpRouter.GET("/sso/v1/login", func(c *gin.Context) {
		urlRedirect, _ := c.GetQuery("redirect")
		_, urlErr := url.ParseRequestURI(urlRedirect)
		if urlErr != nil {
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
	httpRouter.GET("/sso/v1/callback", func(c *gin.Context) {
		//Fetch url parameters
		UrlState, stateErr := c.GetQuery("state")
		if stateErr == false {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-101",
				"ErrorMessage": "Cant verify user, please re-login. If the issue presist contact IT office.",
			})
		}
		UrlCode, codeErr := c.GetQuery("code")
		if codeErr == false {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-102",
				"ErrorMessage": "Cant verify user, please re-login. If the issue presist contact IT office.",
			})
		}
		//check state with google module
		if strings.Compare(UrlState, module.GetStateCode()) != 0 {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-103",
				"ErrorMessage": "Cant verify user, please re-login. If the issue presist contact IT office.",
			})
		}
		UserSession, UserData, authErr := module.AuthHandler(c, UrlCode, UrlState)
		if authErr != nil {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-104",
				"ErrorMessage": "Cant login user, please re-login. If the issue presist contact IT office.",
			})
		}
		//Create full Session and add to DB
		urlRedirect, dberr, _ := module.LoginUserIntoDB(c, UserData, UserSession)
		if dberr == false {
			c.HTML(http.StatusOK, "error500.html", gin.H{
				"ErrorTitle":   "Error SSO-105",
				"ErrorMessage": "DB Error check logs, please re-login. If the issue presist contact IT office.",
			})
		}
		location := url.URL{Path: urlRedirect}
		c.Redirect(http.StatusFound, location.RequestURI())
	})
	httpRouter.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "pagenotfound.html", gin.H{})
	})
	//httpRouter.RunTLS(":1234", "certs/server.crt", "certs/server.key")
	httpRouter.Run(":1234")

}

package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sso.scd.edu.om/handler"
	"sso.scd.edu.om/module"
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
		module.CreateSession(c)
		c.HTML(http.StatusOK, "login.html", gin.H{
			"GoogleLink": module.Setup(),
		})
	})
	httpRouter.GET("/sso/v1/callback", func(c *gin.Context) {
		//being with creating a session for the user
		c.HTML(http.StatusOK, "callback.html", gin.H{})
	})
	httpRouter.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "pagenotfound.html", gin.H{})
	})
	//httpRouter.RunTLS(":1234", "certs/server.crt", "certs/server.key")
	httpRouter.Run(":1234")

}

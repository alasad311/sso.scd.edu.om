package module

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func SessionConnection() gin.HandlerFunc {
	store, _ := redis.NewStore(10, "tcp", "192.168.11.154:6379", "4a2d1148ce53a323cdf62a65d24b8e70dedd4457e7eb2fa9250f62b8b74009b5", []byte("cL]Pept6#AOM4L~"))
	store.Options(sessions.Options{MaxAge: 3600, Domain: "scd.edu.om", HttpOnly: true})
	return sessions.Sessions("SCD-SSO", store)
}

func CreateSession(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("ClientIP") == nil {
		//Build User Details and begin authentation
		session.Set("ClientIP", c.ClientIP())
		session.Set("ClientAgent", c.Request.UserAgent())
		session.Set("ClientB", c.Request.UserAgent())
		session.Save()
	} else {
		fmt.Println(session.Get("ClientIP"))
	}

}

package module

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func CreateSession() gin.HandlerFunc {
	store, _ := redis.NewStore(10, "tcp", "192.168.11.154:6379", "4a2d1148ce53a323cdf62a65d24b8e70dedd4457e7eb2fa9250f62b8b74009b5", []byte("secret"))
	store.Options(sessions.Options{MaxAge: 3600 * 12, Domain: "scd.edu.om", Secure: true, HttpOnly: true})
	return sessions.Sessions("mysession", store)
}

func GetSessionCurrentAge(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("ClientIP") == nil {
		session.Set("ClientIP", c.ClientIP())
		session.Set("ClientIP", c.ClientIP())
		session.Set("ClientIP", c.ClientIP())
		session.Set("ClientIP", c.ClientIP())
		session.Set("ClientIP", c.ClientIP())
		session.Set("ClientIP", c.ClientIP())
		session.Set("ClientIP", c.ClientIP())
	}

}

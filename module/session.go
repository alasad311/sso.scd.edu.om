package module

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"sso.scd.edu.om/structs"
	"time"
)

func SessionConnection() gin.HandlerFunc {
	store, _ := redis.NewStore(10, "tcp", "192.168.11.154:6379", "4a2d1148ce53a323cdf62a65d24b8e70dedd4457e7eb2fa9250f62b8b74009b5", []byte("cL]Pept6#AOM4L~"))
	store.Options(sessions.Options{MaxAge: 12 * 3600, Domain: "scd.edu.om", HttpOnly: true})
	return sessions.Sessions("SCD-SSO", store)
}

func CreateSession(c *gin.Context, url string) bool {
	session := sessions.Default(c)
	if session.Get("TokenID") == nil {
		fmt.Println("Session for " + c.ClientIP() + " has started!")
		session.Set("ClientIP", c.ClientIP())
		session.Set("StateCode", GetStateCode())
		session.Set("SessionStart", time.Now().Unix())
		session.Set("SessionExpire", time.Now().Add(time.Minute*1).Unix())
		session.Set("UrlRedirect", url)
		session.Save()
		return false
	} else {
		current := time.Now().Unix()
		ExpireSession := session.Get("SessionExpire").(int64)
		fmt.Println(time.Unix(ExpireSession, 0).Sub(time.Unix(current, 0)))
		if time.Unix(ExpireSession, 0).Sub(time.Unix(current, 0)) <= 0 {
			session.Clear()
			session.Save()
			return false
		} else {
			return true
		}

	}
}

func LoginUserIntoDB(c *gin.Context, UserData structs.UserLogin, userSession structs.UserSession) (string, bool, error) {
	session := sessions.Default(c)
	userSession.SessionID = session.ID()
	userSession.SessionStart = time.Now()
	session.Set("SessionStart", time.Now().Unix())
	userSession.SessionExpire = time.Now().Add(time.Hour * 1)
	session.Set("SessionExpire", time.Now().Add(time.Hour*1).Unix())
	//DB connection
	db, dbError := DBConnectionSSO()
	if dbError != nil {
		return "", false, dbError
	}
	//check if user exsist into the database
	usersResult := db.Where("claim_user_email = ?", UserData.ClaimUserEmail).FirstOrCreate(&UserData)
	//result := db.Create(&UserData)
	if usersResult.Error != nil {
		return "", false, usersResult.Error
	}
	//insert session
	userSession.Userid = UserData.Userid
	fmt.Println(UserData.Userid)
	createSession := db.Create(&userSession)
	if createSession.Error != nil {
		return "", false, createSession.Error
	}

	return session.Get("UrlRedirect").(string), true, nil
}

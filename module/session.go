package module

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"sso.scd.edu.om/structs"
)

func SessionConnection() gin.HandlerFunc {
	store, _ := redis.NewStore(10, "tcp", "192.168.11.154:6379", "4a2d1148ce53a323cdf62a65d24b8e70dedd4457e7eb2fa9250f62b8b74009b5", []byte("cL]Pept6#AOM4L~"))
	store.Options(sessions.Options{Path: "/", MaxAge: 12 * 3600, Domain: "scd.edu.om", HttpOnly: true})
	return sessions.Sessions("SCD-SSO", store)
}

func CreateSession(c *gin.Context, url string) bool {
	session := sessions.Default(c)
	if session.Get("TokenID") == nil {
		fmt.Println("Session for " + c.ClientIP() + " has started!")
		session.Set("ClientIP", c.ClientIP())
		session.Set("StateCode", GetStateCode())
		session.Set("SessionStart", time.Now().Unix())
		session.Set("SessionExpire", time.Now().Add(time.Hour*5).Unix())
		session.Set("UrlRedirect", url)
		session.Save()
		return false
	} else {
		current := time.Now().Unix()
		ExpireSession := session.Get("SessionExpire").(int64)
		if time.Unix(ExpireSession, 0).Sub(time.Unix(current, 0)).Seconds() <= 0 {
			session.Clear()
			session.Save()
			return false
		} else {
			return true
		}

	}
}

func LoginUserIntoDB(c *gin.Context, UserData structs.UserLogin, userSession structs.UserSession, db *gorm.DB) (string, bool, error) {
	session := sessions.Default(c)
	userSession.SessionID = session.ID()
	userSession.SessionStart = time.Now().Unix()
	var SessionExpire = time.Now().Add(time.Hour * 5).Unix()
	session.Set("SessionStart", time.Now().Unix())
	userSession.SessionExpire = SessionExpire
	session.Set("SessionExpire", SessionExpire)
	session.Set("UserEmail", UserData.ClaimUserEmail)
	//check if user exsist into the database
	usersResult := db.Where("claim_user_email = ?", UserData.ClaimUserEmail).FirstOrCreate(&UserData)
	//result := db.Create(&UserData)
	if usersResult.Error != nil {
		return "", false, usersResult.Error
	}
	//insert session
	userSession.Userid = UserData.Userid
	userSession.CreatedOn = time.Now().Unix()
	createSession := db.Create(&userSession)
	if createSession.Error != nil {
		return "", false, createSession.Error
	}

	session.Save()
	return session.Get("UrlRedirect").(string), true, nil
}

func ExtendUserSession(c *gin.Context, db *gorm.DB, url string) (bool, string, error) {
	session := sessions.Default(c)
	current := time.Now().Unix()
	var userSession structs.UserSession

	if session.Get("UserEmail") != nil {
		usersResult := db.Last(&userSession, "session_id = ? ", session.ID())
		if usersResult.RowsAffected == 0 {
			return false, "/sso/v1/login?redirect=" + url, errors.New("session cant be found!")
		}
		ExpireSession := session.Get("SessionExpire").(int64)
		if time.Unix(ExpireSession, 0).Sub(time.Unix(current, 0)) <= 0 {
			session.Clear()
			session.Save()
			return false, "/sso/v1/login?redirect=" + url, errors.New("session has expired!")
		} else if time.Unix(ExpireSession, 0).Sub(time.Unix(current, 0)).Seconds() <= 300 && time.Unix(ExpireSession, 0).Sub(time.Unix(current, 0)).Seconds() > 0 {
			var newExpire = time.Now().Add(time.Hour * 5).Unix()
			session.Set("SessionExpire", newExpire)
			session.Save()
			updateSession := db.Model(&userSession).Updates(map[string]interface{}{"session_expire": newExpire, "session_extended": true, "updated_on": time.Now().Unix()})
			if updateSession.Error != nil {
				return false, "/sso/v1/login?redirect=" + url, errors.New("session cant be updated!")
			}
			return true, url, nil
		} else {
			return false, "/sso/v1/login?redirect=" + url, errors.New("session is still alive!")
		}
	} else {
		return false, "/sso/v1/login?redirect=" + url, errors.New("session cant be found.")
	}
}

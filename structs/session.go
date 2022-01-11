package structs

type UserSession struct {
	ID                uint   `gorm:"column:id; PRIMARY_KEY"`
	Userid            uint   `gorm:"column:userid"`
	UserIP            string `gorm:"column:user_ip"`
	UserAgent         string `gorm:"column:user_agent"`
	UserState         string `gorm:"column:user_state"`
	GoogleCode        string `gorm:"column:google_code"`
	GoogleAccessToken string `gorm:"column:google_access_token"`
	GoogleTokenType   string `gorm:"column:google_token_type"`
	GoogleScope       string `gorm:"column:google_scope"`
	GoogleIDToken     string `gorm:"column:google_id_token"`
	GoogleExpireIn    int64  `gorm:"column:google_expire_in"`
	SessionID         string `gorm:"column:session_id"`
	SessionStart      int64  `gorm:"column:session_start"`
	SessionExpire     int64  `gorm:"column:session_expire"`
	CreatedOn         int64  `gorm:"column:created_on"`
	SessionExtended   bool   `gorm:"column:session_extended"`
}

package structs

type UserLogin struct {
	Userid              uint   `gorm:"column:userID; PRIMARY_KEY" json:"userID"`
	ClaimUserSub        string `json:"sub"`
	ClaimUserEmail      string `json:"email"`
	ClaimUserFaimlyname string `json:"family_name"`
	ClaimUserGivenname  string `json:"given_name"`
	ClaimUserLocale     string `json:"locale"`
	ClaimUserName       string `json:"name"`
	ClaimUserPicture    string `json:"picture"`
}

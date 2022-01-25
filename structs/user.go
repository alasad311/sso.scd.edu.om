package structs

type UserLogin struct {
	Userid         uint   `gorm:"column:userID; PRIMARY_KEY" json:"userID"`
	ClaimUserSub   string `json:"sub"`
	ClaimUserEmail string `json:"email"`
}

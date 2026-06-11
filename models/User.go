package models

type User struct {
	Username string `gorm:"column:USERNAME;index"`
	Password string `gorm:"column:PASSWORD"`
	Channel  string `gorm:"column:CHANNEL"`
}

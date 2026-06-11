package models

import "time"

type HistoryLog struct {
	ID              int       `gorm:"column:ID;primaryKey;autoIncrement"`
	Timestamp       time.Time `gorm:"column:TIMESTAMP"`
	User            string    `gorm:"column:USER"`
	Menu            string    `gorm:"column:MENU"`
	Action          string    `gorm:"column:ACTION"`
	NewValue        string    `gorm:"column:NEW_VALUE"`
	OldValue        string    `gorm:"column:OLD_VALUE"`
	ResponseMessage string    `gorm:"column:RESPONSE_MESSAGE"`
	IpAddress       string    `gorm:"column:IP_ADDRESS"`
}

func (HistoryLog) TableName() string {
	return "HISTORYLOG"
}

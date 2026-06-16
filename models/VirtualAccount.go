package models

import "time"

type VirtualAccount struct {
	ID           string     `gorm:"column:ID;primaryKey"`
	VANumber     string     `gorm:"column:VA_NUMBER"`
	CustomerID   string     `gorm:"column:CUSTOMER_ID"`
	CustomerName string     `gorm:"column:CUSTOMER_NAME"`
	Amount       float64    `gorm:"column:AMOUNT"`
	ExpiredAt    time.Time  `gorm:"column:EXPIRED_AT"`
	Description  string     `gorm:"column:DESCRIPTION"`
	Status       int        `gorm:"column:STATUS"`
	StatusName   string     `gorm:"column:STATUS_NAME"`
	Action       string     `gorm:"column:ACTION"`
	ReferenceID  string     `gorm:"column:REFERENCE_ID"`
	PaidAt       *time.Time `gorm:"column:PAID_AT"`
	CreatedAt    time.Time  `gorm:"column:CREATED_AT"`
	UpdatedAt    time.Time  `gorm:"column:UPDATED_AT"`
}

func (VirtualAccount) TableName() string {
	return "virtual_accounts"
}

type OCServiceLog struct {
	RowId           int       `gorm:"column:ROW_ID;primaryKey;autoIncrement"`
	Timestamp       time.Time `gorm:"column:TIMESTAMP"`
	UserID          string    `gorm:"column:USER_ID"`
	Menu            string    `gorm:"column:MENU"`
	Action          string    `gorm:"column:ACTION"`
	NewValue        string    `gorm:"column:NEW_VALUE"`
	OldValue        string    `gorm:"column:OLD_VALUE"`
	ResponseMessage string    `gorm:"column:RESPONSE_MESSAGE"`
	IpClient        string    `gorm:"column:IP_CLIENT"`
}

func (OCServiceLog) TableName() string {
	return "OCSERVICELOG"
}

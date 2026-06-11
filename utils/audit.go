package utils

import (
	"context"
	"virtual_account_api/models"

	"gorm.io/gorm"
)

func InsertHistoryLog(ctx context.Context, db *gorm.DB, history models.HistoryLog) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	query := `INSERT INTO HISTORYLOG (
		[TIMESTAMP], [USER], [MENU], [ACTION], 
		NEW_VALUE, OLD_VALUE, RESPONSE_MESSAGE, IP_ADDRESS
	) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`

	_, err = sqlDB.ExecContext(ctx, query,
		history.Timestamp, history.User, history.Menu, history.Action,
		history.NewValue, history.OldValue, history.ResponseMessage, history.IpAddress,
	)
	return err
}

func InsertOCSserviceLog(ctx context.Context, db *gorm.DB, serviceLog models.OCServiceLog) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	query := `INSERT INTO OCSERVICELOG (
		[TIMESTAMP], USER_ID, [MENU], [ACTION], 
		NEW_VALUE, OLD_VALUE, RESPONSE_MESSAGE, IP_CLIENT
	) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`

	_, err = sqlDB.ExecContext(ctx, query,
		serviceLog.Timestamp, serviceLog.UserID, serviceLog.Menu, serviceLog.Action,
		serviceLog.NewValue, serviceLog.OldValue, serviceLog.ResponseMessage, serviceLog.IpClient,
	)
	return err
}

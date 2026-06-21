package repositories

import (
	"context"
	"fmt"
	"strconv"

	modelVa "virtual_account_api/models"
	"virtual_account_api/resources"

	"gorm.io/gorm"
)

type PaymentRepository struct {
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (r *PaymentRepository) DoPaymentCallback(c context.Context, data *modelVa.Payment, db *gorm.DB) error {

	query := `INSERT INTO payments
			(id, va_number, virtual_account_id, paid_amount, payment_channel, status, raw_payload, created_at) VALUES
			(?, ?, ?, ?,?, ?, ?, ?) `

	err := db.WithContext(c).Exec(query,
		data.ID,
		data.VANumber,
		data.VirtualAccountID,
		data.PaidAmount,
		data.PaymentChannel,
		data.Status,
		data.RawPayload,
		data.CreatedAt).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentRepository) DoGetVAStatus(c context.Context, vaNumber string, db *gorm.DB) (*resources.GetVAResource, error) {
	var va resources.GetVAResource
	err := db.WithContext(c).Table("virtual_accounts va").Select("va.*,vs.status as status_name").Joins("JOIN status_va vs ON vs.id_status = va.status").
		Where("va_number = ?", vaNumber).First(&va).Error
	if err != nil {
		return nil, err
	}

	return &va, nil
}

func (r *PaymentRepository) DoGetPaymentHistory(c context.Context, vaNumber string, limit, offset int, db *gorm.DB) ([]*resources.GetPaymentListResource, int64, error) {

	query := `
		SELECT  
			ID,
			VA_NUMBER,
			VIRTUAL_ACCOUNT_ID,
			PAID_AMOUNT,
			PAYMENT_CHANNEL,
			STATUS,
			RAW_PAYLOAD,
			CREATED_AT,
			status_va.STATUS as status_name
		FROM 
			payments
			JOIN status_va ON status_va.id_status = payments.status
		WHERE 
			1=1 `

	args := []interface{}{}
	paramIdx := 1

	if vaNumber != "" {
		query = query + " AND VA_NUMBER = @p" + strconv.Itoa(paramIdx)
		args = append(args, vaNumber)
		paramIdx++
	}

	SqlDB, err := db.DB()
	if err != nil {
		return nil, 0, fmt.Errorf("get db connection failed: %w", err)
	}

	var total int64
	result := SqlDB.QueryRowContext(c, `SELECT COUNT(*) FROM (`+query+`) AS tb`, args...).Scan(&total)
	if result != nil {
		return nil, 0, fmt.Errorf("count transaction: %w", result.Error)
	}

	query += " ORDER BY ID DESC"

	if limit != -1 {
		query += fmt.Sprintf("  OFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY", paramIdx, paramIdx+1)
	}
	args = append(args, offset)
	args = append(args, limit)

	sqlRows, err := SqlDB.QueryContext(c, query, args...)

	if err != nil {
		return nil, 0, fmt.Errorf("query get all transaction failed: %w", err)
	}

	defer sqlRows.Close()

	results := []*resources.GetPaymentListResource{}
	//var tt time.Time

	for sqlRows.Next() {
		var py resources.GetPaymentListResource
		err := sqlRows.Scan(
			&py.Id,
			&py.VANumber,
			&py.PaidAmount,
			&py.PaymentChannel,
			&py.StatusName,
			&py.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transaction row failed: %w", err)
		}
		results = append(results, &py)
	}

	return results, total, nil
}

package repositories

import (
	"context"
	"fmt"
	"strconv"

	modelVa "virtual_account_api/models"
	"virtual_account_api/resources"

	"gorm.io/gorm"
)

type VirtualAccountRepository struct {
}

func NewVirtualAccountRepository() *VirtualAccountRepository {
	return &VirtualAccountRepository{}
}

func (r *VirtualAccountRepository) DoCreateVA(c context.Context, data *modelVa.VirtualAccount, db *gorm.DB) error {

	query := `INSERT INTO virtual_accounts (id, va_number, customer_id, customer_name, amount, description,
			 status, reference_id, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.WithContext(c).Exec(query,
		data.ID,
		data.VANumber,
		data.CustomerID,
		data.CustomerName,
		data.Amount,
		data.Description,
		data.Status,
		data.ReferenceID,
		data.ExpiredAt,
		data.CreatedAt,
		data.UpdatedAt).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *VirtualAccountRepository) DoGetVAStatus(c context.Context, vaNumber string, db *gorm.DB) (*resources.GetVAResource, error) {
	var va resources.GetVAResource
	err := db.WithContext(c).Table("virtual_accounts va").Select("va.*,vs.status as status_name").Joins("JOIN status_va vs ON vs.id_status = va.status").
		Where("va_number = ?", vaNumber).First(&va).Error
	if err != nil {
		return nil, err
	}

	return &va, nil
}

func (r *VirtualAccountRepository) DoGetVA(c context.Context, custId, status string, limit, offset int, db *gorm.DB) ([]*resources.GetVAListResource, int64, error) {

	query := `
		SELECT  
			ID,
			VA_NUMBER,
			CUSTOMER_ID,
			CUSTOMER_NAME,
			AMOUNT,
			DESCRIPTION,
			EXPIRED_AT,
			REFERENCE_ID,
			ISNULL(PAID_AT, '') as PAID_AT,
			CREATED_AT,
			UPDATED_AT,
			status_va.STATUS as status_name
		FROM 
			virtual_accounts 
			JOIN status_va ON status_va.id_status = virtual_accounts.status
		WHERE 
			1=1 `

	args := []interface{}{}
	paramIdx := 1

	if custId != "" {
		query = query + " AND CUSTOMER_ID = @p" + strconv.Itoa(paramIdx)
		args = append(args, custId)
		paramIdx++
	}

	if status != "" {
		query = query + " AND virtual_accounts.STATUS = @p" + strconv.Itoa(paramIdx)
		args = append(args, status)
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

	results := []*resources.GetVAListResource{}
	//var tt time.Time

	for sqlRows.Next() {
		var va resources.GetVAListResource
		err := sqlRows.Scan(
			&va.Id,
			&va.VANumber,
			&va.CustomerId,
			&va.CustomerName,
			&va.Amount,
			&va.Description,
			&va.ExpiredAt,
			&va.ReferenceId,
			&va.PaidAt,
			&va.CreatedAt,
			&va.UpdatedAt,
			&va.StatusName)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transaction row failed: %w", err)
		}
		results = append(results, &va)
	}

	return results, total, nil
}

func (r *VirtualAccountRepository) CheckRefId(c context.Context, refid string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.WithContext(c).Table("virtual_accounts").
		Where("reference_id = ?", refid).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

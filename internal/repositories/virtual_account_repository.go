package repositories

import (
	"context"

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
	err := db.WithContext(c).Where("va_number = ?", vaNumber).First(&va).Error
	if err != nil {
		return nil, err
	}

	return &va, nil
}

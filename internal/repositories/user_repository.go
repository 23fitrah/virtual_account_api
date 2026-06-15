package repositories

import (
	"context"
	"database/sql"
	"virtual_account_api/models"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUserWS(c context.Context, db *gorm.DB, username string) (*models.User, error) {
	var result models.User

	query := `SELECT TOP 1 USERNAME, PASSWORD FROM service_user WITH (NOLOCK) WHERE USERNAME = ?`

	err := db.WithContext(c).Raw(query, username).Scan(&result).Error

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

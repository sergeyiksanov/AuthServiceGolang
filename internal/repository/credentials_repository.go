package repository

import (
	"AuthService/internal/entity"
	"gorm.io/gorm"
)

type CredentialsRepository struct {
	Repository[entity.Credentials]
}

func NewCredentialsRepository() *CredentialsRepository {
	return &CredentialsRepository{}
}

func (cr *CredentialsRepository) GetCountByEmail(db *gorm.DB, email string) (int64, error) {
	var count int64
	err := db.Model(new(entity.Credentials)).Where("email = ?", email).Count(&count).Error
	return count, err
}

func (cr *CredentialsRepository) GetByEmail(db *gorm.DB, email string, entity *entity.Credentials) error {
	return db.Where("email = ?", email).Take(entity).Error
}

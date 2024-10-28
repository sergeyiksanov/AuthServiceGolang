package repository

import (
	"AuthService/internal/dto"
	"gorm.io/gorm"
)

type CredentialsRepository struct {
	Repository[dto.CredentialsDto]
}

func NewCredentialsRepository() *CredentialsRepository {
	return &CredentialsRepository{}
}

func (cr *CredentialsRepository) GetCountByEmail(db *gorm.DB, email string) (int64, error) {
	var count int64
	err := db.Model(new(dto.CredentialsDto)).Where("email = ?", email).Count(&count).Error
	return count, err
}

func (cr *CredentialsRepository) GetByEmail(db *gorm.DB, email string, dto *dto.CredentialsDto) error {
	return db.Where("email = ?", email).Take(dto).Error
}

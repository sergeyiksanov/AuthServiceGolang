package repository

import (
	"AuthService/internal/dto"
	"gorm.io/gorm"
)

type TokensRepository struct {
	Repository[dto.TokenDto]
}

func NewTokensRepository() *TokensRepository {
	return &TokensRepository{}
}

func (ts *TokensRepository) GetTokenByJTI(db *gorm.DB, jti string, dto *dto.TokenDto) error {
	return db.Where("jti = ?", jti).Take(dto).Error
}

func (ts *TokensRepository) RevokeAllTokensWithBySubjectId(db *gorm.DB, subjectId int64) error {
	return db.Model(&dto.TokenDto{}).Where("subject_id = ?", subjectId).Update("revoked", true).Error
}

func (ts *TokensRepository) RevokeTokenByJTI(db *gorm.DB, jti string) error {
	return db.Model(&dto.TokenDto{}).Where("jti = ?", jti).Update("revoked", true).Error
}

package repository

import (
	"AuthService/internal/entity"
	"gorm.io/gorm"
)

type TokensRepository struct {
	Repository[entity.Token]
}

func (*TokensRepository) GetTokenByJTI(db *gorm.DB, jti string, entity *entity.Token) error {
	return db.Where("jti = ?", jti).Take(entity).Error
}

func (*TokensRepository) RevokeAllTokensWithBySubjectId(db *gorm.DB, subjectId int64) error {
	return db.Model(&entity.Token{}).Where("subject_id = ?", subjectId).Update("revoked", true).Error
}

func (*TokensRepository) RevokeTokenByJTI(db *gorm.DB, jti string) error {
	return db.Model(&entity.Token{}).Where("jti = ?", jti).Update("revoked", true).Error
}

func NewTokensRepository() *TokensRepository {
	return &TokensRepository{}
}

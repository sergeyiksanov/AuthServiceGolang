package service

import (
	"AuthService/internal/dto"
	"gorm.io/gorm"
)

type credentialsRepository interface {
	GetById(db *gorm.DB, dto *dto.CredentialsDto, id any) error
	Create(db *gorm.DB, dto *dto.CredentialsDto) error
	Update(db *gorm.DB, dto *dto.CredentialsDto) error
	Delete(db *gorm.DB, dto *dto.CredentialsDto) error
	GetCountById(db *gorm.DB, id any) (int64, error)
	GetCountByEmail(db *gorm.DB, email string) (int64, error)
	GetByEmail(db *gorm.DB, email string, entity *dto.CredentialsDto) error
}

type tokensRepository interface {
	GetById(db *gorm.DB, dto *dto.TokenDto, id any) error
	Create(db *gorm.DB, dto *dto.TokenDto) error
	Update(db *gorm.DB, dto *dto.TokenDto) error
	Delete(db *gorm.DB, dto *dto.TokenDto) error
	GetCountById(db *gorm.DB, id any) (int64, error)
	GetTokenByJTI(db *gorm.DB, jti string, entity *dto.TokenDto) error
	RevokeAllTokensWithBySubjectId(db *gorm.DB, subjectId int64) error
	RevokeTokenByJTI(db *gorm.DB, jti string) error
}

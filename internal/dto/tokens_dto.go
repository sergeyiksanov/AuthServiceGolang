package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type TokenDto struct {
	JTI       string           `gorm:"column_id:jti,primaryKey"`
	SubjectId int64            `gorm:"column_id:subject_id"`
	TokenType string           `gorm:"column_id:token_type"`
	Revoked   bool             `gorm:"column_id:revoked"`
	IssuedAt  pgtype.Timestamp `gorm:"column_id:issued_at"`
	ExpiresAt pgtype.Timestamp `gorm:"column:expires_at"`
}

func (c TokenDto) TableName() string {
	return "issued_jwt_token"
}

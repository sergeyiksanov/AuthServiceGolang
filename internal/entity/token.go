package entity

import "github.com/jackc/pgx/v5/pgtype"

type Token struct {
	JTI       string
	SubjectId int64
	TokenType string
	Revoked   bool
	IssuedAt  pgtype.Timestamp
	ExpiresAt pgtype.Timestamp
}

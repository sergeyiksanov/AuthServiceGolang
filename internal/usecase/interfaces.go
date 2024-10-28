package usecase

import (
	"AuthService/internal/dto"
	"AuthService/internal/entity"
	"context"
)

type credentialsService interface {
	GetCredentialsById(ctx context.Context, id int64) (entity.Credentials, error)
	CheckAlreadyExistsEmail(ctx context.Context, email string) (bool, error)
	HashPassword(password string) (string, error)
	ValidatePassword(password, hash string) bool
	CreateCredentials(ctx context.Context, credentials entity.Credentials) error
	GetCredentialsByEmail(ctx context.Context, email string) (entity.Credentials, error)
}

type tokensService interface {
	RevokeTokenByJTI(ctx context.Context, jti string) error
	RevokeAllTokensWithBySubjectId(ctx context.Context, subjectId int64) error
	GetTokenByJTI(ctx context.Context, jti string) (dto.TokenDto, error)
	VerifyToken(ctx context.Context, tokenString string, expectedType string) (string, error)
	CreateAccessRefreshPairTokens(ctx context.Context, credentialsId int64, email string) (string, string, error)
}

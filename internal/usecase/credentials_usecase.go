package usecase

import (
	"AuthService/internal/entity"
	"AuthService/internal/repository"
	"AuthService/internal/utils"
	proto "AuthService/pkg/api/v1"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"log"
)

type CredentialsUseCase struct {
	DB                    *gorm.DB
	CredentialsRepository *repository.CredentialsRepository
	TokensRepository      *repository.TokensRepository
}

func NewCredentialsUseCase(db *gorm.DB, cr *repository.CredentialsRepository, tr *repository.TokensRepository) *CredentialsUseCase {
	return &CredentialsUseCase{
		DB:                    db,
		CredentialsRepository: cr,
		TokensRepository:      tr,
	}
}

func (c *CredentialsUseCase) Logout(ctx context.Context, req *proto.LogoutRequest) (*emptypb.Empty, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	jwtAccessToken, err := utils.VerifyAccessToken(req.Tokens.Access)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	jwtRefreshToken, err := utils.VerifyRefreshToken(req.Tokens.Refresh)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	jtiAccess, err := utils.GetJTIFromToken(jwtAccessToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	jtiRefresh, err := utils.GetJTIFromToken(jwtRefreshToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := c.TokensRepository.RevokeTokenByJTI(tx, jtiAccess); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := c.TokensRepository.RevokeTokenByJTI(tx, jtiRefresh); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to commit transaction: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (c *CredentialsUseCase) RefreshTokens(ctx context.Context, req *proto.RefreshTokensRequest) (*proto.RefreshTokensResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	jwtRefreshToken, err := utils.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	jti, err := utils.GetJTIFromToken(jwtRefreshToken)
	if err != nil {
		return nil, err
	}

	token := new(entity.Token)
	if err := c.TokensRepository.GetTokenByJTI(tx, jti, token); err != nil {
		return nil, utils.InvalidRefreshToken
	}

	if token.Revoked {
		return nil, utils.RevokedRefreshToken
	}

	if err := c.TokensRepository.RevokeAllTokensWithBySubjectId(tx, token.SubjectId); err != nil {
		return nil, utils.InternalServerError
	}

	credentials := new(entity.Credentials)
	if err := c.CredentialsRepository.GetById(tx, credentials, token.SubjectId); err != nil {
		return nil, utils.InternalServerError
	}

	accessToken, refreshToken, accessTokenString, refreshTokenString, err := utils.CreatePairTokens(credentials.ID, credentials.Email)

	if err := c.TokensRepository.Create(tx, accessToken); err != nil {
		log.Fatalf("Failed save access token: %v", err)
		return nil, utils.InternalServerError
	}

	if err := c.TokensRepository.Create(tx, refreshToken); err != nil {
		log.Fatalf("Failed save refresh token: %v", err)
		return nil, utils.InternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, utils.InternalServerError
	}

	return &proto.RefreshTokensResponse{
		Tokens: &proto.Tokens{
			Refresh: refreshTokenString,
			Access:  accessTokenString,
		},
	}, nil
}

func (c *CredentialsUseCase) VerifyAccessToken(ctx context.Context, req *proto.VerifyAccessTokenRequest) (*proto.VerifyAccessTokenResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	jwtToken, err := utils.VerifyAccessToken(req.Access)
	if err != nil {
		return nil, err
	}

	jti, err := utils.GetJTIFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	token := new(entity.Token)
	if err := c.TokensRepository.GetTokenByJTI(tx, jti, token); err != nil {
		log.Fatalf("GetTokenByJTI error: %v", err)
		return nil, utils.InvalidAccessToken
	}

	if token.TokenType != "access" {
		return nil, utils.InvalidAccessToken
	}

	if token.Revoked {
		return nil, utils.InvalidAccessToken
	}

	return &proto.VerifyAccessTokenResponse{
		UserId: token.SubjectId,
	}, nil
}

func (c *CredentialsUseCase) SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.SignInResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	count, err := c.CredentialsRepository.GetCountByEmail(tx, req.Credentials.Email)
	if err != nil {
		log.Fatalf("Failed to get count credentials by email: %v", err)
		return nil, err
	}

	if count == 0 {
		return nil, status.Errorf(codes.NotFound, "Credentials not found")
	}

	credentials := new(entity.Credentials)
	if err := c.CredentialsRepository.GetByEmail(tx, req.Credentials.Email, credentials); err != nil {
		log.Fatalf("Failed get user: %v", err)
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Credentials.Password, credentials.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Password")
	}

	accessToken, refreshToken, accessTokenString, refreshTokenString, err := utils.CreatePairTokens(credentials.ID, credentials.Email)
	if err != nil {
		log.Fatalf("Failed to create access token: %v", err)
		return nil, utils.InternalServerError
	}

	if err := c.TokensRepository.Create(tx, accessToken); err != nil {
		log.Fatalf("Failed save access token: %v", err)
		return nil, utils.InternalServerError
	}

	if err := c.TokensRepository.Create(tx, refreshToken); err != nil {
		log.Fatalf("Failed save refresh token: %v", err)
		return nil, utils.InternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, utils.InternalServerError
	}

	return &proto.SignInResponse{
		Tokens: &proto.Tokens{
			Refresh: refreshTokenString,
			Access:  accessTokenString,
		},
	}, nil
}

func (c *CredentialsUseCase) SignUp(ctx context.Context, req *proto.SignUpRequest) (*emptypb.Empty, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	count, err := c.CredentialsRepository.GetCountByEmail(tx, req.Credentials.Email)
	if err != nil {
		log.Fatalf("Failed to get count creaditianals with email: +%v", err)
		return nil, utils.InternalServerError
	}

	if count > 0 {
		return nil, utils.EmailAlreadyExists
	}

	hashPassword, err := utils.HashPassword(req.Credentials.Password)
	if err != nil {
		log.Fatalf("Failed to hash password: +%v", err)
		return nil, utils.InternalServerError
	}

	user := &entity.Credentials{
		Email:    req.Credentials.Email,
		Password: hashPassword,
	}

	if err := c.CredentialsRepository.Create(tx, user); err != nil {
		return nil, utils.InternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, utils.InternalServerError
	}

	return &emptypb.Empty{}, nil
}

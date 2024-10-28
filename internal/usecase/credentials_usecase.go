package usecase

import (
	"AuthService/internal/entity"
	"AuthService/internal/utils"
	proto "AuthService/pkg/api/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CredentialsUseCase struct {
	crs credentialsService
	ts  tokensService
}

func NewCredentialsUseCase(crs credentialsService, ts tokensService) *CredentialsUseCase {
	return &CredentialsUseCase{
		crs: crs,
		ts:  ts,
	}
}

func (c CredentialsUseCase) Logout(ctx context.Context, req *proto.LogoutRequest) (*emptypb.Empty, error) {
	jti, err := c.ts.VerifyToken(ctx, req.Tokens.Access, "access")
	if err != nil {
		return nil, err
	}
	jti, err = c.ts.VerifyToken(ctx, req.Tokens.Refresh, "refresh")
	if err != nil {
		return nil, err
	}

	if err := c.ts.RevokeTokenByJTI(ctx, jti); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c CredentialsUseCase) RefreshTokens(ctx context.Context, req *proto.RefreshTokensRequest) (*proto.RefreshTokensResponse, error) {
	jti, err := c.ts.VerifyToken(ctx, req.RefreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	token, err := c.ts.GetTokenByJTI(ctx, jti)
	if err != nil {
		return nil, err
	}

	if token.Revoked {
		return nil, utils.RevokedToken
	}

	if err := c.ts.RevokeAllTokensWithBySubjectId(ctx, token.SubjectId); err != nil {
		return nil, err
	}

	credentials, err := c.crs.GetCredentialsById(ctx, token.SubjectId)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := c.ts.CreateAccessRefreshPairTokens(ctx, credentials.ID, credentials.Email)
	if err != nil {
		return nil, err
	}

	return &proto.RefreshTokensResponse{
		Tokens: &proto.Tokens{
			Refresh: refreshToken,
			Access:  accessToken,
		},
	}, nil
}

func (c CredentialsUseCase) VerifyAccessToken(ctx context.Context, req *proto.VerifyAccessTokenRequest) (*proto.VerifyAccessTokenResponse, error) {
	jti, err := c.ts.VerifyToken(ctx, req.Access, "access")
	if err != nil {
		return nil, err
	}

	tokenDto, err := c.ts.GetTokenByJTI(ctx, jti)
	if err != nil {
		return nil, err
	}

	if tokenDto.TokenType != "access" {
		return nil, utils.InvalidToken
	}

	if tokenDto.Revoked {
		return nil, utils.InvalidToken
	}

	return &proto.VerifyAccessTokenResponse{
		UserId: tokenDto.SubjectId,
	}, nil
}

func (c CredentialsUseCase) SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.SignInResponse, error) {
	res, err := c.crs.CheckAlreadyExistsEmail(ctx, req.Credentials.Email)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, utils.InvalidCredentials
	}

	credentials, err := c.crs.GetCredentialsByEmail(ctx, req.Credentials.Email)
	if err != nil {
		return nil, err
	}

	if !c.crs.ValidatePassword(req.Credentials.Password, credentials.Password) {
		return nil, utils.InvalidCredentials
	}

	accessToken, refreshToken, err := c.ts.CreateAccessRefreshPairTokens(ctx, credentials.ID, credentials.Email)
	if err != nil {
		return nil, err
	}

	return &proto.SignInResponse{
		Tokens: &proto.Tokens{
			Refresh: refreshToken,
			Access:  accessToken,
		},
	}, nil
}

func (c CredentialsUseCase) SignUp(ctx context.Context, req *proto.SignUpRequest) error {
	res, err := c.crs.CheckAlreadyExistsEmail(ctx, req.Credentials.Email)
	if err != nil {
		return err
	}
	if res {
		return utils.EmailAlreadyExists
	}

	hash, err := c.crs.HashPassword(req.Credentials.Password)
	if err != nil {
		return err
	}

	credentials := entity.Credentials{
		Email:    req.Credentials.Email,
		Password: hash,
	}

	if err := c.crs.CreateCredentials(ctx, credentials); err != nil {
		return err
	}

	return nil
}

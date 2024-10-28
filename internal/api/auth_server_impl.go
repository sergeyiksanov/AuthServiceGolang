package api

import (
	"AuthService/internal/usecase"
	desc "AuthService/pkg/api/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthImplementationSever struct {
	desc.UnimplementedAuthServer
	credentialsUseCase *usecase.CredentialsUseCase
}

func NewAuthImplementationSever(useCase *usecase.CredentialsUseCase) *AuthImplementationSever {
	return &AuthImplementationSever{
		credentialsUseCase: useCase,
	}
}

func (is *AuthImplementationSever) RefreshTokens(ctx context.Context, req *desc.RefreshTokensRequest) (*desc.RefreshTokensResponse, error) {
	return is.credentialsUseCase.RefreshTokens(ctx, req)
}

func (is *AuthImplementationSever) SignUp(ctx context.Context, req *desc.SignUpRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, is.credentialsUseCase.SignUp(ctx, req)
}

func (is *AuthImplementationSever) SignIn(ctx context.Context, req *desc.SignInRequest) (*desc.SignInResponse, error) {
	return is.credentialsUseCase.SignIn(ctx, req)
}

func (is *AuthImplementationSever) VerifyAccessToken(ctx context.Context, req *desc.VerifyAccessTokenRequest) (*desc.VerifyAccessTokenResponse, error) {
	return is.credentialsUseCase.VerifyAccessToken(ctx, req)
}

func (is *AuthImplementationSever) Logout(ctx context.Context, req *desc.LogoutRequest) (*emptypb.Empty, error) {
	return is.credentialsUseCase.Logout(ctx, req)
}

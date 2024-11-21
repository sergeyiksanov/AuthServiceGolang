package api

import (
	"AuthService/internal/metrics"
	"AuthService/internal/usecase"
	desc "AuthService/pkg/api/v1"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	start := time.Now()
	resp, err := is.credentialsUseCase.RefreshTokens(ctx, req)
	defer func() {
		code := codes.OK
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				code = codes.Internal
			} else {
				code = st.Code()
			}
		}
		metrics.ObserveRefreshTokensRequest(time.Since(start), code)
	}()

	return resp, err
}

func (is *AuthImplementationSever) SignUp(ctx context.Context, req *desc.SignUpRequest) (*emptypb.Empty, error) {
	start := time.Now()
	err := is.credentialsUseCase.SignUp(ctx, req)
	defer func() {
		code := codes.OK
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				code = codes.Internal
			} else {
				code = st.Code()
			}
		}
		metrics.ObserveSignUpRequest(time.Since(start), code)
	}()
	return &emptypb.Empty{}, err
}

func (is *AuthImplementationSever) SignIn(ctx context.Context, req *desc.SignInRequest) (*desc.SignInResponse, error) {
	start := time.Now()
	resp, err := is.credentialsUseCase.SignIn(ctx, req)
	defer func() {
		code := codes.OK
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				code = codes.Internal
			} else {
				code = st.Code()
			}
		}
		metrics.ObserveSignInRequest(time.Since(start), code)
	}()

	return resp, err
}

func (is *AuthImplementationSever) VerifyAccessToken(ctx context.Context, req *desc.VerifyAccessTokenRequest) (*desc.VerifyAccessTokenResponse, error) {
	start := time.Now()
	resp, err := is.credentialsUseCase.VerifyAccessToken(ctx, req)
	defer func() {
		code := codes.OK
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				code = codes.Internal
			} else {
				code = st.Code()
			}
			metrics.ObserveVerifyAccessTokenRequest(time.Since(start), code)
		}
	}()
	return resp, err
}

func (is *AuthImplementationSever) Logout(ctx context.Context, req *desc.LogoutRequest) (*emptypb.Empty, error) {
	start := time.Now()
	resp, err := is.credentialsUseCase.Logout(ctx, req)
	defer func() {
		code := codes.OK
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				code = codes.Internal
			} else {
				code = st.Code()
			}
			metrics.ObserveLogoutRequest(time.Since(start), code)
		}
	}()

	return resp, err
}

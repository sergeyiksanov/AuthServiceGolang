package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// TOKEN ERRORS
	RevokedAccessToken  = status.Error(codes.Unauthenticated, "Access token revoked")
	RevokedRefreshToken = status.Error(codes.Unauthenticated, "Refresh token revoked")
	InvalidAccessToken  = status.Error(codes.Unauthenticated, "Invalid access token")
	InvalidRefreshToken = status.Error(codes.Unauthenticated, "Invalid refresh token")

	// CREDENTIALS ERRORS
	EmailAlreadyExists = status.Error(codes.AlreadyExists, "Email already exists")

	// OTHER ERRORS
	InternalServerError = status.Error(codes.Internal, "Internal server error")
)

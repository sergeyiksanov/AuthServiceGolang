package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// TOKEN ERRORS
	RevokedToken = status.Error(codes.Unauthenticated, "Token revoked")
	InvalidToken = status.Error(codes.Unauthenticated, "Invalid token")

	// CREDENTIALS ERRORS
	EmailAlreadyExists = status.Error(codes.AlreadyExists, "Email already exists")
	InvalidCredentials = status.Error(codes.NotFound, "Credentials not found")

	// OTHER ERRORS
	InternalServerError = status.Error(codes.Internal, "Internal server error")
)

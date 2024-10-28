package app

import (
	"AuthService/internal/api"
	"AuthService/internal/config"
	"AuthService/internal/repository"
	"AuthService/internal/service"
	"AuthService/internal/usecase"
	"gorm.io/gorm"
	"log"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig

	credentialsRepository *repository.CredentialsRepository

	tokensRepository *repository.TokensRepository

	credentialsUseCase *usecase.CredentialsUseCase

	authServerImpl *api.AuthImplementationSever

	credentialsService *service.CredentialsService

	tokensService *service.TokensService

	gormDB *gorm.DB
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("Failed to initialize gRPC config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) CredentialsRepository() *repository.CredentialsRepository {
	if s.credentialsRepository == nil {
		s.credentialsRepository = repository.NewCredentialsRepository()
	}

	return s.credentialsRepository
}

func (s *serviceProvider) CredentialsUseCase() *usecase.CredentialsUseCase {
	if s.credentialsUseCase == nil {
		s.credentialsUseCase = usecase.NewCredentialsUseCase(s.CredentialsService(), s.TokensService())
	}

	return s.credentialsUseCase
}

func (s *serviceProvider) CredentialsService() *service.CredentialsService {
	if s.credentialsService == nil {
		s.credentialsService = service.NewCredentialsService(s.GormDB(), s.CredentialsRepository(), s.TokensRepository())
	}

	return s.credentialsService
}

func (s *serviceProvider) TokensService() *service.TokensService {
	if s.tokensService == nil {
		s.tokensService = service.NewTokensService(s.GormDB(), s.CredentialsRepository(), s.TokensRepository())
	}

	return s.tokensService
}

func (s *serviceProvider) GormDB() *gorm.DB {
	if s.gormDB == nil {
		s.gormDB = config.NewDatabase()
	}

	return s.gormDB
}

func (s *serviceProvider) AuthServerImpl() *api.AuthImplementationSever {
	if s.authServerImpl == nil {
		s.authServerImpl = api.NewAuthImplementationSever(s.CredentialsUseCase())
	}

	return s.authServerImpl
}

func (s *serviceProvider) TokensRepository() *repository.TokensRepository {
	if s.tokensRepository == nil {
		s.tokensRepository = repository.NewTokensRepository()
	}

	return s.tokensRepository
}

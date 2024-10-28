package service

import (
	"AuthService/internal/convertor"
	"AuthService/internal/dto"
	"AuthService/internal/entity"
	"context"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CredentialsService struct {
	db        *gorm.DB
	crRepo    credentialsRepository
	tokenRepo tokensRepository
}

func NewCredentialsService(db *gorm.DB, crRepo credentialsRepository, tokensRepo tokensRepository) *CredentialsService {
	return &CredentialsService{
		db:        db,
		crRepo:    crRepo,
		tokenRepo: tokensRepo,
	}
}

func (cr *CredentialsService) CheckAlreadyExistsEmail(ctx context.Context, email string) (bool, error) {
	tx := cr.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	cnt, err := cr.crRepo.GetCountByEmail(tx, email)
	if err != nil {
		return false, err
	}

	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

func (cr *CredentialsService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func (cr *CredentialsService) ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (cr *CredentialsService) CreateCredentials(ctx context.Context, credentials entity.Credentials) error {
	tx := cr.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	credentialsDto := convertor.CredentialsEntityToCredentialsDto(credentials)

	if err := cr.crRepo.Create(tx, &credentialsDto); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (cr *CredentialsService) GetCredentialsByEmail(ctx context.Context, email string) (entity.Credentials, error) {
	tx := cr.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	credentialsDto := new(dto.CredentialsDto)
	if err := cr.crRepo.GetByEmail(tx, email, credentialsDto); err != nil {
		return entity.Credentials{}, err
	}

	return credentialsDto.ToCredentialsEntity(), nil
}

func (cr *CredentialsService) GetCredentialsById(ctx context.Context, id int64) (entity.Credentials, error) {
	tx := cr.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	credentialsDto := new(dto.CredentialsDto)
	if err := cr.crRepo.GetById(tx, credentialsDto, id); err != nil {
		return entity.Credentials{}, err
	}

	return credentialsDto.ToCredentialsEntity(), nil
}

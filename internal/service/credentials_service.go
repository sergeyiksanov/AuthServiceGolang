package service

import (
	"AuthService/internal/convertor"
	"AuthService/internal/dto"
	"AuthService/internal/entity"
	"AuthService/internal/external"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
)

type CredentialsService struct {
	db        *gorm.DB
	crRepo    credentialsRepository
	tokenRepo tokensRepository
	external  *external.NotificationExternal
}

func NewCredentialsService(db *gorm.DB, crRepo credentialsRepository, tokensRepo tokensRepository, external *external.NotificationExternal) *CredentialsService {
	return &CredentialsService{
		db:        db,
		crRepo:    crRepo,
		tokenRepo: tokensRepo,
		external:  external,
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

func generateRandomCode() string {
	letterBytes := "0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (cr *CredentialsService) SendConfirmRegistrationMailToEmail(email string) (string, error) {
	code := generateRandomCode()
	req := entity.EmailEventNotificationEntity{
		Name:  "Подтверждение регистрации",
		Title: "Подтвердите регистрацию",
		Body:  "Код подтверждения регистрации: " + code + ".\nНикому не сообщайте код.",
		Email: email,
	}

	if err := cr.external.SendEmailEventNotification(&req); err != nil {
		return "", err
	}

	return code, nil
}


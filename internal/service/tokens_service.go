package service

import (
	"AuthService/internal/dto"
	"AuthService/internal/utils"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

const secretKeyName = "JWT_SECRET_KEY"
const refreshLifeTimeName = "JWT_REFRESH_LIFE_TIME_DAY"
const accessLifeTimeName = "JWT_ACCESS_LIFE_TIME_MINUTE"

const (
	accessToken  = "access"
	refreshToken = "refresh"
)

type TokensService struct {
	db     *gorm.DB
	crRepo credentialsRepository
	tRepo  tokensRepository
}

type tokenInfo struct {
	tokenString string
	jti         string
	exp         int64
	typeToken   string
}

func NewTokensService(db *gorm.DB, crRepo credentialsRepository, tRepo tokensRepository) *TokensService {
	return &TokensService{
		db:     db,
		crRepo: crRepo,
		tRepo:  tRepo,
	}
}

func (ts *TokensService) RevokeTokenByJTI(ctx context.Context, jti string) error {
	tx := ts.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ts.tRepo.RevokeTokenByJTI(tx, jti); err != nil {
		return err
	}

	return nil
}

func (ts *TokensService) RevokeAllTokensWithBySubjectId(ctx context.Context, subjectId int64) error {
	tx := ts.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ts.tRepo.RevokeAllTokensWithBySubjectId(tx, subjectId); err != nil {
		return err
	}

	return tx.Commit().Error
}

func (ts *TokensService) GetTokenByJTI(ctx context.Context, jti string) (dto.TokenDto, error) {
	tx := ts.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	tokenDto := new(dto.TokenDto)
	if err := ts.tRepo.GetTokenByJTI(tx, jti, tokenDto); err != nil {
		return dto.TokenDto{}, err
	}

	return *tokenDto, nil
}

func (ts *TokensService) VerifyToken(ctx context.Context, tokenString string, expectedType string) (string, error) {
	tx := ts.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	secretKey := os.Getenv(secretKeyName)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		typeToken, ok := claims["type"].(string)
		if !ok {
			return "", utils.InvalidToken
		}

		if typeToken != expectedType {
			return "", utils.InvalidToken
		}

		jti, ok := claims["jti"].(string)
		if !ok {
			return "", utils.InvalidToken
		}

		tokenDto := new(dto.TokenDto)
		if err := ts.tRepo.GetTokenByJTI(tx, jti, tokenDto); err != nil {
			return "", err
		}

		if tokenDto.Revoked {
			return "", utils.RevokedToken
		}

		return jti, nil
	}

	return "", utils.InvalidToken
}

func (ts *TokensService) CreateAccessRefreshPairTokens(ctx context.Context, credentialsId int64, email string) (string, string, error) {
	tx := ts.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	refreshToken, err := ts.createJWTToken(credentialsId, email, "refresh")
	if err != nil {
		return "", "", err
	}

	accessToken, err := ts.createJWTToken(credentialsId, email, "access")
	if err != nil {
		return "", "", err
	}

	refreshTokenEntity := dto.TokenDto{
		JTI:       refreshToken.jti,
		SubjectId: credentialsId,
		TokenType: refreshToken.typeToken,
		Revoked:   false,
	}

	accessTokenEntity := dto.TokenDto{
		JTI:       accessToken.jti,
		SubjectId: credentialsId,
		TokenType: accessToken.typeToken,
		Revoked:   false,
	}

	if err := ts.tRepo.Create(tx, &accessTokenEntity); err != nil {
		return "", "", err
	}

	if err := ts.tRepo.Create(tx, &refreshTokenEntity); err != nil {
		return "", "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", "", err
	}

	return accessToken.tokenString, refreshToken.tokenString, nil
}

func (ts *TokensService) createJWTToken(credentialsId int64, email string, typeToken string) (tokenInfo, error) {
	secretKey := os.Getenv(secretKeyName)
	if len(secretKey) == 0 {
		log.Fatalf("%s is empty", secretKeyName)
		return tokenInfo{}, utils.InternalServerError
	}

	var (
		lifeTime int64
		err      error
		exp      int64
	)

	if typeToken == "refresh" {
		lifeTime, err = strconv.ParseInt(os.Getenv(refreshLifeTimeName), 0, 64)
		if err != nil || lifeTime < 1 {
			log.Fatalf("Invalid %s", refreshLifeTimeName)
			return tokenInfo{}, utils.InternalServerError
		}

		exp = time.Now().Add(time.Hour * 24 * time.Duration(lifeTime)).Unix()
	} else if typeToken == "access" {
		lifeTime, err = strconv.ParseInt(os.Getenv(accessLifeTimeName), 0, 64)
		if err != nil || lifeTime < 1 {
			log.Fatalf("Invalid %s", accessLifeTimeName)
			return tokenInfo{}, utils.InternalServerError
		}

		exp = time.Now().Add(time.Minute * time.Duration(lifeTime)).Unix()
	} else {
		return tokenInfo{}, utils.InternalServerError
	}

	jti := uuid.New().String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    credentialsId,
			"email": email,
			"exp":   exp,
			"jti":   jti,
			"type":  typeToken,
		},
	)

	tokenString, err := ts.getTokenString(token)

	return tokenInfo{
		tokenString: tokenString,
		jti:         jti,
		exp:         exp,
		typeToken:   typeToken,
	}, nil
}

func (ts *TokensService) getTokenString(token *jwt.Token) (string, error) {
	secretKey := os.Getenv(secretKeyName)
	if len(secretKey) == 0 {
		log.Fatalf("%s is empty", secretKeyName)
		return "", utils.InternalServerError
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", utils.InternalServerError
	}

	return tokenString, nil
}

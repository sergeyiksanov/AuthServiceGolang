package utils

import (
	"AuthService/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"time"
)

const secretKeyName = "JWT_SECRET_KEY"
const refreshLifeTimeName = "JWT_REFRESH_LIFE_TIME_DAY"
const accessLifeTimeName = "JWT_ACCESS_LIFE_TIME_MINUTE"

type tokenStruct struct {
	tokenString string
	jti         string
	exp         int64
	typeToken   string
}

func CreatePairTokens(userId int64, email string) (*entity.Token, *entity.Token, string, string, error) {
	refreshToken, err := createJWTToken(userId, email, "refresh")
	if err != nil {
		log.Fatalf("Failed to create refresh token: %v", err)
		return nil, nil, "", "", InternalServerError
	}

	accessToken, err := createJWTToken(userId, email, "access")
	if err != nil {
		log.Fatalf("Failed to create access token: %v", err)
		return nil, nil, "", "", InternalServerError
	}

	accessTokenPsql := &entity.Token{
		JTI:       accessToken.jti,
		SubjectId: userId,
		TokenType: accessToken.typeToken,
		Revoked:   false,
	}

	refreshTokenPsql := &entity.Token{
		JTI:       refreshToken.jti,
		SubjectId: userId,
		TokenType: refreshToken.typeToken,
		Revoked:   false,
	}

	return accessTokenPsql, refreshTokenPsql, accessToken.tokenString, refreshToken.tokenString, nil
}

func createJWTToken(userId int64, email string, typeToken string) (*tokenStruct, error) {
	secretKey := os.Getenv(secretKeyName)
	if len(secretKey) == 0 {
		log.Fatalf("Secret key not found in environment")
		return nil, InternalServerError
	}

	var lifeTime int
	var err error
	var exp int64
	if typeToken == "refresh" {
		lifeTime, err = strconv.Atoi(os.Getenv(refreshLifeTimeName))
		if err != nil {
			log.Fatalf("Invalid lifetime jwt token in env: %v", err)
			return nil, InternalServerError
		}
		exp = time.Now().Add(time.Hour * 24 * time.Duration(lifeTime)).Unix()
	} else if typeToken == "access" {
		lifeTime, err = strconv.Atoi(os.Getenv(accessLifeTimeName))
		if err != nil {
			log.Fatalf("Invalid lifetime jwt token in env: %v", err)
			return nil, InternalServerError
		}
		exp = time.Now().Add(time.Hour * 24 * time.Duration(lifeTime)).Unix()
	} else {
		return nil, InternalServerError
	}

	jti := uuid.New().String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    userId,
			"email": email,
			"exp":   exp,
			"jti":   jti,
			"type":  typeToken,
		},
	)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
		return nil, InternalServerError
	}

	return &tokenStruct{
		tokenString: tokenString,
		jti:         jti,
		exp:         exp,
		typeToken:   typeToken,
	}, nil
}

func VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv(secretKeyName)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("Failed to parse token: %v", err)
		return nil, InvalidAccessToken
	}

	if !token.Valid {
		log.Fatalf("Exp")
		return nil, InvalidAccessToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		typeToken, ok := claims["type"].(string)
		if !ok {
			log.Printf("None type")
			return nil, InvalidAccessToken
		}

		if typeToken != "access" {
			log.Printf("Invalid type token")
			return nil, InvalidAccessToken
		}
	} else {
		log.Printf("None claims")
		return nil, InvalidAccessToken
	}

	return token, nil
}

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv(secretKeyName)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, InvalidRefreshToken
	}

	if !token.Valid {
		return nil, InvalidRefreshToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		typeToken, ok := claims["type"].(string)
		if !ok {
			return nil, InvalidRefreshToken
		}

		if typeToken != "refresh" {
			return nil, InvalidRefreshToken
		}
	} else {
		return nil, InvalidRefreshToken
	}

	return token, nil
}

func GetJTIFromToken(jwtToken *jwt.Token) (string, error) {
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		jti, ok := claims["jti"].(string)
		if !ok {
			return "", InternalServerError
		}
		return jti, nil
	}
	return "", InternalServerError
}

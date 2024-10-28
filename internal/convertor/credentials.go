package convertor

import (
	"AuthService/internal/dto"
	"AuthService/internal/entity"
)

func CredentialsEntityToCredentialsDto(c entity.Credentials) dto.CredentialsDto {
	return dto.CredentialsDto{
		ID:       c.ID,
		Email:    c.Email,
		Password: c.Password,
	}
}

func CredentialsDtoToCredentialsEntity(c dto.CredentialsDto) entity.Credentials {
	return entity.Credentials{
		ID:       c.ID,
		Email:    c.Email,
		Password: c.Password,
	}
}

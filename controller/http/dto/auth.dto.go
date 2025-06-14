package dto

import "d-payroll/entity"

type LoginBodyDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDto struct {
	Token string `json:"token"`
}

func (l *LoginBodyDto) ToLoginEntity() *entity.Login {
	return &entity.Login{
		Username: l.Username,
		Password: l.Password,
	}
}

func (l *LoginResponseDto) FromAuthToken(authToken *entity.AuthToken) {
	l.Token = authToken.Token
}

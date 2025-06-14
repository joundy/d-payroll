package http

import (
	ctxresponse "d-payroll/controller/http/customctx"
	"d-payroll/controller/http/dto"
	internalerror "d-payroll/internal-error"
	authservice "d-payroll/service/auth"
	"d-payroll/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AuthHttp struct {
	http    *httpApp
	authSvc authservice.AuthService
}

func NewAuthHttp(http *httpApp, authSvc authservice.AuthService) *AuthHttp {
	authHttp := &AuthHttp{
		http:    http,
		authSvc: authSvc,
	}

	authHttp.http.App.Post("/login", authHttp.Login)

	return authHttp
}

func (a *AuthHttp) Login(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	loginBody := new(dto.LoginBodyDto)
	if err := c.BodyParser(loginBody); err != nil {
		return cc.BadRequest("Invalid request body")
	}

	err := utils.ValidateStruct(loginBody)
	if err != nil {
		return err
	}

	auth, err := a.authSvc.Login(c.Context(), loginBody.ToLoginEntity())
	if err != nil {
		if errors.Is(err, &internalerror.NotFoundError{}) {
			return cc.NotFound("User not found")
		}

		if errors.Is(err, &internalerror.InvalidCredentialsError{}) {
			return cc.Unauthorized("Invalid credentials")
		}

		return err
	}

	var response dto.LoginResponseDto
	response.FromAuthToken(auth)

	return cc.Ok(response, nil)
}

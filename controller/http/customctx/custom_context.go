package customctx

import (
	"d-payroll/entity"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CustomContext struct {
	*fiber.Ctx
}

func (c *CustomContext) Ok(data any, msg *string) error {
	message := "Success"
	if msg != nil {
		message = *msg
	}
	return c.Ctx.JSON(entity.HttpResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (c *CustomContext) jsonError(status int, msg string) error {
	return c.Ctx.Status(status).JSON(entity.HttpResponse{
		Success: false,
		Message: msg,
		Data:    nil,
	})
}

func (c *CustomContext) BadRequest(msg string) error {
	return c.jsonError(fiber.StatusBadRequest, msg)
}

func (c *CustomContext) Unauthorized(msg string) error {
	return c.jsonError(fiber.StatusUnauthorized, msg)
}

func (c *CustomContext) Forbidden(msg string) error {
	return c.jsonError(fiber.StatusForbidden, msg)
}

func (c *CustomContext) NotFound(msg string) error {
	return c.jsonError(fiber.StatusNotFound, msg)
}

func (c *CustomContext) UnprocessableEntity(msg string) error {
	return c.jsonError(fiber.StatusUnprocessableEntity, msg)
}

func (c *CustomContext) Conflict(msg string) error {
	return c.jsonError(fiber.StatusConflict, msg)
}

func (c *CustomContext) GetAuthPayload() (*entity.AuthTokenPayload, error) {
	if c.Locals("authPayload") != nil {
		authPayload := c.Locals("authPayload").(*entity.AuthTokenPayload)
		return authPayload, nil
	}

	return nil, fmt.Errorf("AuthPayload not found")
}

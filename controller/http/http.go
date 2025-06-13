package http

import (
	"d-payroll/config"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	"errors"
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

func (c *CustomContext) ErrorBadRequest(msg string) error {
	return c.Ctx.Status(fiber.StatusBadRequest).JSON(entity.HttpResponse{
		Success: false,
		Message: msg,
		Data:    nil,
	})
}

type httpApp struct {
	config *config.Config
	App    *fiber.App
}

func NewHttpApp(config *config.Config) *httpApp {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var validationError *internalerror.ValidationError
			if errors.As(err, &validationError) {
				return c.Status(fiber.StatusBadRequest).JSON(entity.HttpResponse{
					Success: false,
					Message: "Validation Error",
					Data:    validationError.Fields,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(entity.HttpResponse{
				Success: false,
				Message: "Internal Server Error",
			})
		},
	})

	return &httpApp{
		config: config,
		App:    app,
	}
}

func (h *httpApp) Listen() {
	h.App.Listen(fmt.Sprintf("%s:%d", h.config.Http.Host, h.config.Http.Port))
}

package http

import (
	"d-payroll/config"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

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

			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(entity.HttpResponse{
				Success: false,
				Message: "Internal Server Error",
			})
		},
	})

	app.Get("/_health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(entity.HttpResponse{
			Success: true,
			Message: "Ok",
		})
	})

	return &httpApp{
		config: config,
		App:    app,
	}
}

func (h *httpApp) Listen() {
	h.App.Listen(fmt.Sprintf("%s:%d", h.config.Http.Host, h.config.Http.Port))
}

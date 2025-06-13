package middleware

import (
	"d-payroll/config"
	ctxresponse "d-payroll/controller/http/customctx"
	"d-payroll/entity"
	"d-payroll/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Authorization(config *config.Config, roles []entity.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cc := ctxresponse.CustomContext{Ctx: c}

		token := c.Get("Authorization")
		if token == "" {
			return cc.Unauthorized("Missing Authorization header")
		}

		bearerToken := strings.TrimPrefix(token, "Bearer ")
		if bearerToken == "" {
			return cc.Unauthorized("Invalid Bearer token")
		}

		payload, err := utils.VerifyToken(config.Auth.JwtSecret, bearerToken)
		if err != nil || payload == nil {
			return cc.Unauthorized("Invalid or expired token")
		}

		if !utils.ArrContains(roles, payload.Role) {
			return cc.Forbidden("User is not allowed")
		}

		c.Locals("authPayload", payload)

		return c.Next()
	}
}

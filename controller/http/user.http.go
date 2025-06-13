package http

import (
	ctxresponse "d-payroll/controller/http/ctx-response"
	"d-payroll/controller/http/dto"
	"d-payroll/controller/http/middleware"
	"d-payroll/entity"
	userservice "d-payroll/service/user"
	"d-payroll/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHttp struct {
	userSvc userservice.UserService
}

func NewUserHttp(h *httpApp, userSvc userservice.UserService) {
	userHttp := &UserHttp{
		userSvc: userSvc,
	}

	h.App.Post("/users", middleware.Authorization(h.config, []entity.UserRole{entity.UserRoleEmployee}), userHttp.CreateUser)
	h.App.Get("/users/:id", middleware.Authorization(h.config, []entity.UserRole{entity.UserRoleEmployee}), userHttp.getUserById)
}

func (u *UserHttp) CreateUser(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	user := new(dto.CreateUserBodyDto)
	if err := c.BodyParser(user); err != nil {
		return cc.BadRequest("Invalid request body")
	}

	err := utils.ValidateStruct(user)
	if err != nil {
		return err
	}

	createdUser, err := u.userSvc.CreateUser(c.Context(), user.ToUserEntity())
	if err != nil {
		return err
	}

	var response dto.CreateUserResponseDto
	response.FromUserEntity(createdUser)

	return cc.Ok(response, nil)
}

func (u *UserHttp) getUserById(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return cc.BadRequest("Invalid ID param")
	}

	user, err := u.userSvc.GetUserById(c.Context(), idInt)
	if err != nil {
		return err
	}

	var response dto.GetUserByIdResponseDto
	response.FromUserEntity(user)

	return cc.Ok(response, nil)
}

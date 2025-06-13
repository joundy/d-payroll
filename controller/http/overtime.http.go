package http

import (
	ctxresponse "d-payroll/controller/http/customctx"
	"d-payroll/controller/http/dto"
	"d-payroll/controller/http/middleware"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	overtimeservice "d-payroll/service/overtime"
	"d-payroll/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OvertimeHttp struct {
	http        *httpApp
	overtimeSvc overtimeservice.OvertimeService
}

func NewOvertimeHttp(http *httpApp, overtimeSvc overtimeservice.OvertimeService) {
	overtimeHttp := &OvertimeHttp{
		http:        http,
		overtimeSvc: overtimeSvc,
	}

	overtimeHttp.http.App.Post("/overtimes", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleEmployee}), overtimeHttp.CreateOvertime)
	overtimeHttp.http.App.Post("/overtimes/:overtimeId/approve", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin}), overtimeHttp.ApproveOvertime)
	overtimeHttp.http.App.Get("/overtimes", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleEmployee, entity.UserRoleAdmin}), overtimeHttp.GetUserOvertimes)
}

func (o *OvertimeHttp) CreateOvertime(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	overtime := new(dto.CreateOvertimeBodyDto)
	if err := c.BodyParser(overtime); err != nil {
		return cc.BadRequest("Invalid request body")
	}

	err = utils.ValidateStruct(overtime)
	if err != nil {
		return err
	}

	createdOvertime, err := o.overtimeSvc.CreateOvertime(c.Context(), overtime.ToOvertimeEntity(authPayload.ID))
	if err != nil {
		return err
	}

	var response dto.OvertimeResponseDto
	response.FromOvertimeEntity(createdOvertime)

	return cc.Ok(response, nil)
}

func (o *OvertimeHttp) ApproveOvertime(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	overtimeId := c.Params("overtimeId")
	overtimeIdInt, err := strconv.ParseUint(overtimeId, 10, 32)
	if err != nil {
		return cc.BadRequest("Invalid overtime ID param")
	}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	err = o.overtimeSvc.ApproveOvertime(c.Context(), uint(overtimeIdInt), authPayload.ID)
	if err != nil {
		if errors.Is(err, &internalerror.OvertimeAlreadyApprovedError{}) {
			return cc.Conflict("Overtime already approved")
		}

		if errors.Is(err, &internalerror.NotFoundError{}) {
			return cc.NotFound("Overtime not found")
		}

		return err
	}

	return cc.Ok(nil, nil)
}

func (o *OvertimeHttp) GetUserOvertimes(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	userIdParam := c.Query("user_id")
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		return cc.BadRequest("Invalid user ID query")
	}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	if authPayload.Role == entity.UserRoleEmployee && authPayload.ID != uint(userId) {
		return cc.Unauthorized("Unauthorized to access other user's overtimes")
	}

	overtimes, err := o.overtimeSvc.GetOvertimesByUserID(c.Context(), uint(userId))
	if err != nil {
		return err
	}

	responses := make([]*dto.OvertimeResponseDto, len(overtimes))
	for i, overtime := range overtimes {
		var response dto.OvertimeResponseDto
		response.FromOvertimeEntity(overtime)
		responses[i] = &response
	}

	return cc.Ok(responses, nil)
}

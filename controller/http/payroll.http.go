package http

import (
	ctxresponse "d-payroll/controller/http/customctx"
	"d-payroll/controller/http/dto"
	"d-payroll/controller/http/middleware"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	payrollservice "d-payroll/service/payroll"
	"d-payroll/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PayrollHttp struct {
	http       *httpApp
	payrollSvc payrollservice.PayrollService
}

func NewPayrollHttp(http *httpApp, payrollSvc payrollservice.PayrollService) {
	payrollHttp := &PayrollHttp{
		http:       http,
		payrollSvc: payrollSvc,
	}

	payrollHttp.http.App.Post("/payrolls", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin}), payrollHttp.CreatePayroll)
	payrollHttp.http.App.Get("/payrolls", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin}), payrollHttp.GetUserPayrolls)
	payrollHttp.http.App.Post("/payrolls/:payrollId/roll", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin}), payrollHttp.RollPayroll)
	payrollHttp.http.App.Post("/payrolls/:payrollId/payslips", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin, entity.UserRoleEmployee}), payrollHttp.Payslips)

	payrollHttp.http.App.Post("/payrolls/:payrollId/payslip-summaries", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin}), payrollHttp.PayslipSummaries)
	payrollHttp.http.App.Post("/payrolls/:payrollId/total-take-home-pay", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin}), payrollHttp.PayslipTotalTakeHomePay)
}

func (p *PayrollHttp) CreatePayroll(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	payroll := new(dto.CreatePayrollBodyDto)
	if err := c.BodyParser(payroll); err != nil {
		return cc.BadRequest("Invalid request body")
	}

	err = utils.ValidateStruct(payroll)
	if err != nil {
		return err
	}

	createdPayroll, err := p.payrollSvc.CreatePayroll(c.Context(), payroll.ToPayrollEntity(authPayload.ID))
	if err != nil {
		return err
	}

	var response dto.PayrollResponseDto
	response.FromPayrollEntity(createdPayroll)

	return cc.Ok(response, nil)
}

func (p *PayrollHttp) GetUserPayrolls(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	payrolls, err := p.payrollSvc.GetPayrolls(c.Context())
	if err != nil {
		return err
	}

	responses := make([]*dto.PayrollResponseDto, len(payrolls))
	for i, payroll := range payrolls {
		var response dto.PayrollResponseDto
		response.FromPayrollEntity(payroll)
		responses[i] = &response
	}

	return cc.Ok(responses, nil)
}

func (p *PayrollHttp) RollPayroll(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	payrollId := c.Params("payrollId")
	payrollIdInt, err := strconv.ParseUint(payrollId, 10, 32)
	if err != nil {
		return cc.BadRequest("Invalid payroll ID param")
	}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	err = p.payrollSvc.RollPayroll(c.Context(), uint(payrollIdInt), authPayload.ID)
	if err != nil {
		if errors.Is(err, &internalerror.NotFoundError{}) {
			return cc.NotFound("Payroll not found")
		}

		if errors.Is(err, &internalerror.PayrollAlreadyRolledError{}) {
			return cc.Conflict("Payroll already rolled")
		}
		return err
	}

	return cc.Ok(nil, nil)
}

func (p *PayrollHttp) Payslips(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	payrollId := c.Params("payrollId")
	payrollIdInt, err := strconv.ParseUint(payrollId, 10, 32)
	if err != nil {
		return cc.BadRequest("Invalid payroll ID param")
	}

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
		return cc.Unauthorized("Unauthorized to access other user's payslips")
	}

	payslips, err := p.payrollSvc.GeneratePayslip(c.Context(), uint(payrollIdInt), uint(userId))
	if err != nil {
		if errors.Is(err, &internalerror.PayrollNotRolledError{}) {
			return cc.UnprocessableEntity("Payroll is not rolled yet")
		}

		return err
	}

	var response dto.PayslipDto
	response.FromPayslipEntity(payslips)

	return cc.Ok(response, nil)
}

func (p *PayrollHttp) PayslipSummaries(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	payrollId := c.Params("payrollId")
	payrollIdInt, err := strconv.ParseUint(payrollId, 10, 32)
	if err != nil {
		return cc.BadRequest("Invalid payroll ID param")
	}

	payrollSummaries, err := p.payrollSvc.GetPayslipSummaries(c.Context(), uint(payrollIdInt))
	if err != nil {
		if errors.Is(err, &internalerror.PayrollNotRolledError{}) {
			return cc.UnprocessableEntity("Payroll is not rolled yet")
		}

		return err
	}

	responses := make([]*dto.UserPayslipSummaryDto, len(payrollSummaries))
	for i, summary := range payrollSummaries {
		dtoSummary := &dto.UserPayslipSummaryDto{}
		dtoSummary.FromUserPayslipSummaryEntity(summary)
		responses[i] = dtoSummary
	}

	return cc.Ok(responses, nil)
}

func (p *PayrollHttp) PayslipTotalTakeHomePay(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	payrollId := c.Params("payrollId")
	payrollIdInt, err := strconv.ParseUint(payrollId, 10, 32)
	if err != nil {
		return cc.BadRequest("Invalid payroll ID param")
	}

	totalTakeHomePay, err := p.payrollSvc.GetTotalTakeHomePay(c.Context(), uint(payrollIdInt))
	if err != nil {
		return err
	}

	return cc.Ok(totalTakeHomePay, nil)
}

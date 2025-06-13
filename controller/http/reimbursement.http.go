package http

import (
	ctxresponse "d-payroll/controller/http/customctx"
	"d-payroll/controller/http/dto"
	"d-payroll/controller/http/middleware"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	reimbursementservice "d-payroll/service/reimbursement"
	"d-payroll/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReimbursementHttp struct {
	http             *httpApp
	reimbursementSvc reimbursementservice.ReimbursementService
}

func NewReimbursementHttp(http *httpApp, reimbursementSvc reimbursementservice.ReimbursementService) {
	reimbursementHttp := &ReimbursementHttp{
		http:             http,
		reimbursementSvc: reimbursementSvc,
	}

	reimbursementHttp.http.App.Post("/reimbursements", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleEmployee}), reimbursementHttp.CreateReimbursement)
	reimbursementHttp.http.App.Post("/reimbursements/:reimbursementId/approve", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleAdmin}), reimbursementHttp.ApproveReimbursement)
	reimbursementHttp.http.App.Get("/reimbursements", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleEmployee, entity.UserRoleAdmin}), reimbursementHttp.GetUserReimbursements)
}

func (r *ReimbursementHttp) CreateReimbursement(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	reimbursement := new(dto.CreateReimbursementBodyDto)
	if err := c.BodyParser(reimbursement); err != nil {
		return cc.BadRequest("Invalid request body")
	}

	err = utils.ValidateStruct(reimbursement)
	if err != nil {
		return err
	}

	createdReimbursement, err := r.reimbursementSvc.CreateReimbursement(c.Context(), reimbursement.ToReimbursementEntity(authPayload.ID))
	if err != nil {
		return err
	}

	var response dto.ReimbursementResponseDto
	response.FromReimbursementEntity(createdReimbursement)

	return cc.Ok(response, nil)
}

func (r *ReimbursementHttp) ApproveReimbursement(c *fiber.Ctx) error {
	cc := ctxresponse.CustomContext{Ctx: c}

	reimbursementId := c.Params("reimbursementId")
	reimbursementIdInt, err := strconv.ParseUint(reimbursementId, 10, 32)
	if err != nil {
		return cc.BadRequest("Invalid reimbursement ID param")
	}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	err = r.reimbursementSvc.ApproveReimbursement(c.Context(), uint(reimbursementIdInt), authPayload.ID)
	if err != nil {
		if errors.Is(err, &internalerror.ReimbursementAlreadyApprovedError{}) {
			return cc.Conflict("Reimbursement already approved")
		}

		if errors.Is(err, &internalerror.NotFoundError{}) {
			return cc.NotFound("Reimbursement not found")
		}

		return err
	}

	return cc.Ok(nil, nil)
}

func (r *ReimbursementHttp) GetUserReimbursements(c *fiber.Ctx) error {
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
		return cc.Unauthorized("Unauthorized")
	}

	reimbursements, err := r.reimbursementSvc.GetReimbursementsByUserID(c.Context(), uint(userId))
	if err != nil {
		return err
	}

	responses := make([]*dto.ReimbursementResponseDto, len(reimbursements))
	for i, reimbursement := range reimbursements {
		var response dto.ReimbursementResponseDto
		response.FromReimbursementEntity(reimbursement)
		responses[i] = &response
	}

	return cc.Ok(responses, nil)
}

package http

import (
	"d-payroll/controller/http/customctx"
	"d-payroll/controller/http/dto"
	"d-payroll/controller/http/middleware"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	attendanceservice "d-payroll/service/attendance"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AttendanceHttp struct {
	http          *httpApp
	attendanceSvc attendanceservice.AttendanceService
}

func NewAttendanceHttp(http *httpApp, attendanceSvc attendanceservice.AttendanceService) *AttendanceHttp {
	attendanceHttp := &AttendanceHttp{
		http:          http,
		attendanceSvc: attendanceSvc,
	}

	attendanceHttp.http.App.Post("/attendances/checkin", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleEmployee}), attendanceHttp.Checkin)
	attendanceHttp.http.App.Post("/attendances/checkout", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleEmployee}), attendanceHttp.Checkout)
	attendanceHttp.http.App.Get("/attendances", middleware.Authorization(http.config, []entity.UserRole{entity.UserRoleEmployee, entity.UserRoleAdmin}), attendanceHttp.GetAttendancesByUserID)

	return attendanceHttp
}

func (a *AttendanceHttp) Checkin(c *fiber.Ctx) error {
	cc := customctx.CustomContext{Ctx: c}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	attendance, err := a.attendanceSvc.Checkin(c.Context(), authPayload.ID)
	if err != nil {
		if errors.Is(err, &internalerror.AttendanceAlreadyCheckedInError{}) {
			return cc.Conflict("User already checked in")
		}

		if errors.Is(err, &internalerror.AttendanceWeekendError{}) {
			return cc.UnprocessableEntity("User cannot checked in on weekend")
		}
		return err
	}

	var response dto.AttendanceResponseDto
	response.FromUserAttendanceEntity(attendance)

	return cc.Ok(response, nil)
}

func (a *AttendanceHttp) Checkout(c *fiber.Ctx) error {
	cc := customctx.CustomContext{Ctx: c}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	attendance, err := a.attendanceSvc.Checkout(c.Context(), authPayload.ID)
	if err != nil {
		if errors.Is(err, &internalerror.AttendanceCannotCheckedOutError{}) {
			return cc.UnprocessableEntity("User cannot checked out because it is not checked in")
		}

		if errors.Is(err, &internalerror.AttendanceAlreadyCheckedOutError{}) {
			return cc.Conflict("User already checked out")
		}
		return err
	}

	var response dto.AttendanceResponseDto
	response.FromUserAttendanceEntity(attendance)

	return cc.Ok(response, nil)

}

func (a *AttendanceHttp) GetAttendancesByUserID(c *fiber.Ctx) error {
	cc := customctx.CustomContext{Ctx: c}

	authPayload, err := cc.GetAuthPayload()
	if err != nil {
		return err
	}

	attendances, err := a.attendanceSvc.GetAttendancesByUserID(c.Context(), authPayload.ID)
	if err != nil {
		return err
	}

	var responses []*dto.AttendanceResponseDto
	for _, attendance := range attendances {
		var response dto.AttendanceResponseDto
		response.FromUserAttendanceEntity(attendance)
		responses = append(responses, &response)
	}

	return cc.Ok(responses, nil)
}

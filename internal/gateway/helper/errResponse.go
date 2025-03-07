package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SuccessResponse(ctx *fiber.Ctx, status int, msg string, data any) error {
	return ctx.Status(status).JSON(fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func sendErrorResponse(ctx *fiber.Ctx, status int, message string, err error) error {
	if err != nil {
		sLogger.SLogger.Warn(message, "method", ctx.Method(), "path", ctx.Path(), "error", err.Error())
	} else {
		sLogger.SLogger.Warn(message, "method", ctx.Method(), "path", ctx.Path())
	}

	return ctx.Status(status).JSON(ErrorResponse{
		Status:  status,
		Message: message,
	})
}

func InternalServerError(ctx *fiber.Ctx, err error) error {
	sLogger.SLogger.Error("internal error", "method", ctx.Method(), "path", ctx.Path(), "err", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Status:  fiber.StatusInternalServerError,
		Message: "the server encountered a problem",
	})
}

func BadRequest(ctx *fiber.Ctx, err error) error {
	sLogger.SLogger.Error("bad request", "method", ctx.Method(), "path", ctx.Path(), "err", err.Error())
	return sendErrorResponse(ctx, fiber.StatusBadRequest, "bad request", err)
}

func NotFound(ctx *fiber.Ctx, err error) error {
	sLogger.SLogger.Warn("not found", "method", ctx.Method(), "path", ctx.Path(), "err", err.Error())
	return sendErrorResponse(ctx, fiber.StatusNotFound, "not found", err)
}

func ConflictResponse(ctx *fiber.Ctx, err error) error {
	sLogger.SLogger.Error("conflict", "method", ctx.Method(), "path", ctx.Path(), "err", err.Error())
	return ctx.Status(fiber.StatusConflict).JSON(ErrorResponse{
		Status:  fiber.StatusConflict,
		Message: err.Error(),
	})
}

func UnauthorizedErrorResponse(ctx *fiber.Ctx, err error) error {
	sLogger.SLogger.Warn("unauthorized", "method", ctx.Method(), "path", ctx.Path(), "err", err.Error())
	return sendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", err)
}

func UnauthorizedBasicErrorResponse(ctx *fiber.Ctx, err error) error {
	sLogger.SLogger.Warn("unauthorized", "method", ctx.Method(), "err", err.Error())
	ctx.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	return sendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", err)
}

func ForbiddenResponse(ctx *fiber.Ctx) error {
	sLogger.SLogger.Warn("forbidden", "method", ctx.Method(), "path", ctx.Path(), "err", "Forbidden")
	return sendErrorResponse(ctx, fiber.StatusForbidden, "forbidden", nil)
}

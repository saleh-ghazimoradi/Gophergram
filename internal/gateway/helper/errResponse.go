package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
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
	sLogger.SLogger.Error("not found", "method", ctx.Method(), "path", ctx.Path(), "err", err.Error())
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
	sLogger.SLogger.Error("unauthorized", "method", ctx.Method(), "path", ctx.Path(), "err", err.Error())
	return sendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", err)
}

func UnauthorizedBasicErrorResponse(ctx *fiber.Ctx, err error) error {
	sLogger.SLogger.Error("unauthorized", "method", ctx.Method(), "err", err.Error())
	ctx.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	return sendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", err)
}

func ForbiddenResponse(ctx *fiber.Ctx) error {
	sLogger.SLogger.Error("forbidden", "method", ctx.Method(), "path", ctx.Path(), "err", "Forbidden")
	return sendErrorResponse(ctx, fiber.StatusForbidden, "forbidden", nil)
}

func RateLimitExceededResponse(ctx *fiber.Ctx, retryAfter string) error {
	ctx.Set("Retry-After", retryAfter)
	sLogger.SLogger.Warn("rate limit exceeded", "method", ctx.Method(), "path", ctx.Path())
	return ctx.Status(fiber.StatusTooManyRequests).JSON(ErrorResponse{
		Status:  fiber.StatusTooManyRequests,
		Message: "rate limit exceeded, retry after: " + retryAfter,
	})
}

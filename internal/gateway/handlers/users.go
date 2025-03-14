package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

type UserHandler struct {
	userService     service.UserService
	followerService service.FollowService
}

func (u *UserHandler) RegisterUser(ctx *fiber.Ctx) error {
	var payload dto.RegisterUser
	if err := ctx.BodyParser(&payload); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := helper.Validator.Struct(payload); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := u.userService.Create(ctx.Context(), &payload); err != nil {
		switch err {
		case repository.ErrDuplicateEmail:
			return helper.BadRequest(ctx, err)
		case repository.ErrDuplicateUsername:
			return helper.BadRequest(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	return helper.SuccessResponse(ctx, fiber.StatusCreated, "success", nil)
}

func (u *UserHandler) ActivateUserHandler(ctx *fiber.Ctx) error {
	token := ctx.Params("token")

	if err := u.userService.Activate(ctx.Context(), token); err != nil {
		switch err {
		case repository.ErrsNotFound:
			return helper.NotFound(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	return helper.SuccessResponse(ctx, fiber.StatusNoContent, "success", nil)
}

func (u *UserHandler) GetUserHandler(ctx *fiber.Ctx) error {
	user := u.GetUserFromContext(ctx)
	return helper.SuccessResponse(ctx, fiber.StatusOK, "Success", user)
}

func (u *UserHandler) FollowUserHandler(ctx *fiber.Ctx) error {
	followerUser := u.GetUserFromContext(ctx)

	var payload dto.FollowUser
	if err := ctx.BodyParser(&payload); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := u.followerService.Follow(ctx.Context(), followerUser.ID, payload.UserId); err != nil {
		switch {
		case errors.Is(err, repository.ErrConflict):
			return helper.ConflictResponse(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	return helper.SuccessResponse(ctx, fiber.StatusNoContent, "Success", nil)
}

func (u *UserHandler) UnfollowUserHandler(ctx *fiber.Ctx) error {
	unfollowedUser := u.GetUserFromContext(ctx)

	var payload dto.UnfollowUser
	if err := ctx.BodyParser(&payload); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := u.followerService.Unfollow(ctx.Context(), unfollowedUser.ID, payload.UserId); err != nil {
		return helper.InternalServerError(ctx, err)
	}

	return helper.SuccessResponse(ctx, fiber.StatusNoContent, "Success", nil)
}

func (u *UserHandler) UsersContextMiddleware(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return helper.BadRequest(ctx, err)
	}

	user, err := u.userService.GetById(ctx.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrsNotFound):
			return helper.NotFound(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	ctx.Locals(userCtx, user)
	return ctx.Next()
}

func (u *UserHandler) GetUserFromContext(ctx *fiber.Ctx) *service_models.Users {
	user, ok := ctx.Locals(userCtx).(*service_models.Users)
	if !ok {
		sLogger.SLogger.Warn("user not found in context")
		return nil
	}
	return user
}

func NewUserHandler(userService service.UserService, followerService service.FollowService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		followerService: followerService,
	}
}

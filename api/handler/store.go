package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/e-commercer-api/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	"github.com/samber/do"
)

type storeHandler struct {
	i            *do.Injector
	storeService domain.StoreService
	userService  domain.UserService
}

func NewStoreHandler(i *do.Injector) (domain.StoreHandler, error) {
	storeService, err := do.Invoke[domain.StoreService](i)
	if err != nil {
		return nil, err
	}

	userService, err := do.Invoke[domain.UserService](i)
	if err != nil {
		return nil, err
	}

	return &storeHandler{
		i:            i,
		storeService: storeService,
		userService:  userService,
	}, nil
}

func (s *storeHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("func", "Create"),
		slog.String("handler", "store"),
	)

	log.Info("Initializing store create process")

	var storePayload domain.StorePayload
	if err := ctx.Bind(&storePayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := storePayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	storeResponse, err := s.storeService.Create(ctx.Request().Context(), storePayload)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("Store created successfully")

	return ctx.JSON(http.StatusCreated, storeResponse)
}

func (s *storeHandler) GetAll(ctx echo.Context) error {
	log := slog.With(
		slog.String("func", "GetAll"),
		slog.String("handler", "store"),
	)

	log.Info("Initializing get all stores process")

	storeResponse, err := s.storeService.GetAll(ctx.Request().Context())
	if err != nil {
		log.Error("Failed to get all stores", slog.String("error", err.Error()))

		switch {
		case errors.Is(err, domain.ErrUserNotFoundInContext):
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "Forbidden",
				Detail: "User not found in context. Please log in again.",
			})
		case errors.Is(err, domain.ErrStoresNotFound):
			return ctx.NoContent(http.StatusNoContent)
		default:
			return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Oops! Something went wrong while processing your request. Please try again later.",
			})
		}
	}

	log.Info("Successfully retrieved all stores")
	return ctx.JSON(http.StatusOK, storeResponse)
}

func (s *storeHandler) UpdateName(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "store"),
		slog.String("func", "UpdateName"),
	)

	log.Info("Initializing store name update process")

	param := ctx.Param("storeId")

	storeID, err := uuid.Parse(param)
	if err != nil {
		log.Warn("Invalid params", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	var storeNameUpdatePayload domain.StoreNameUpdatePayload
	if err := ctx.Bind(&storeNameUpdatePayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := storeNameUpdatePayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	if err := s.userService.CheckStatus(ctx.Request().Context()); err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailNotConfirmed):
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "unauthorized",
				Detail: "You need to confirm your email to use this feature",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Oops! Something went wrong while processing your request. Please try again later.",
			})
		}
	}

	if err := s.storeService.UpdateName(ctx.Request().Context(), storeID, storeNameUpdatePayload); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("Store name updated successfully")
	return ctx.NoContent(http.StatusOK)
}

func (s *storeHandler) Delete(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "store"),
		slog.String("func", "UpdateName"),
	)

	log.Info("Initializing delete store process")

	param := ctx.Param("storeId")

	storeID, err := uuid.Parse(param)
	if err != nil {
		log.Warn("Invalid params", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	if err := s.userService.CheckStatus(ctx.Request().Context()); err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailNotConfirmed):
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "Unauthorized",
				Detail: "You need to confirm your email to use this feature",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Oops! Something went wrong while processing your request. Please try again later.",
			})
		}
	}

	if err := s.storeService.Delete(ctx.Request().Context(), storeID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("Store deleted successfully")
	return ctx.NoContent(http.StatusNoContent)
}

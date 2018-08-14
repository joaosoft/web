package web

import (
	"db-migration/services"
	"net/http"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/validator"
	"github.com/labstack/echo"
)

type Controller struct {
	logger     logger.ILogger
	interactor *services.Interactor
}

func NewController(logger logger.ILogger, interactor *services.Interactor) *Controller {
	return &Controller{
		logger:     logger,
		interactor: interactor,
	}
}

func (controller *Controller) GetMigrationHandler(ctx echo.Context) error {
	request := GetMigrationRequest{
		IdMigration: ctx.Param("id"),
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		return ctx.JSON(http.StatusBadRequest, errs)
	}

	if process, err := controller.interactor.GetMigration(request.IdMigration); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if process == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, process)
	}
}

func (controller *Controller) GetMigrationsHandler(ctx echo.Context) error {
	if processes, err := controller.interactor.GetMigrations(ctx.QueryParams()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if processes == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, processes)
	}
}

func (controller *Controller) CreateMigrationHandler(ctx echo.Context) error {
	request := CreateMigrationRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		newErr := errors.New("0", err)
		controller.logger.WithFields(map[string]interface{}{"error": err, "cause": newErr.Cause()}).
			Error("error getting body").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if errs := validator.Validate(request.Body); !errs.IsEmpty() {
		newErr := errors.New("0", errs)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	newMigration := services.Migration{
		IdMigration: request.Body.IdMigration,
	}
	if err := controller.interactor.CreateMigration(&newMigration); err != nil {
		newErr := errors.New("0", err)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error creating process %s", request.Body.IdMigration).ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		return ctx.NoContent(http.StatusCreated)
	}
}

func (controller *Controller) DeleteMigrationHandler(ctx echo.Context) error {
	request := DeleteMigrationRequest{
		IdMigration: ctx.Param("id"),
	}

	if errs := validator.Validate(request); !errs.IsEmpty() {
		newErr := errors.New("0", errs)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := controller.interactor.DeleteMigration(request.IdMigration); err != nil {
		newErr := errors.New("0", err)
		controller.logger.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error deleting process by id %s", request.IdMigration).ToErr(newErr)
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (controller *Controller) DeleteMigrationsHandler(ctx echo.Context) error {
	if err := controller.interactor.DeleteMigrations(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

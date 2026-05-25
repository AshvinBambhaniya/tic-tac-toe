package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/AshvinBambhaniya/tic-tac-toe/constants"
	"github.com/AshvinBambhaniya/tic-tac-toe/models"
	"github.com/AshvinBambhaniya/tic-tac-toe/pkg/structs"
	"github.com/AshvinBambhaniya/tic-tac-toe/services"
	"github.com/AshvinBambhaniya/tic-tac-toe/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gofiber/fiber/v2"
)

// UserController for user controllers
type UserController struct {
	userService *services.UserService
	logger      *zap.Logger
}

// NewUserController returns a user
func NewUserController(goqu *goqu.Database, logger *zap.Logger) (*UserController, error) {
	userModel, err := models.InitUserModel(goqu)
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(&userModel)

	return &UserController{
		userService: userSvc,
		logger:      logger,
	}, nil
}

func (ctrl *UserController) GetUser(c *fiber.Ctx) error {
	userIDStr := c.Params(constants.ParamUid)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidUserID)
	}

	user, err := ctrl.userService.GetUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.UserNotExist)
		}
		ctrl.logger.Error("error while get user by id", zap.Any("id", userID), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
	}
	return utils.JSONSuccess(c, http.StatusOK, user)
}

func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {

	var userReq structs.ReqRegisterUser

	err := json.Unmarshal(c.Body(), &userReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(userReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	user, err := ctrl.userService.RegisterUser(models.User{
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
		Email:     userReq.Email,
		Password:  userReq.Password,
		Roles:     userReq.Roles})
	if err != nil {
		ctrl.logger.Error("error while insert user", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrInsertUser)
	}

	return utils.JSONSuccess(c, http.StatusCreated, user)
}

func (ctrl *UserController) GetMe(c *fiber.Ctx) error {
	userIDStr := c.Locals(constants.ContextUid).(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidUserID)
	}

	user, err := ctrl.userService.GetUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.UserNotExist)
		}
		ctrl.logger.Error("error while get current user", zap.Any("id", userID), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
	}
	return utils.JSONSuccess(c, http.StatusOK, user)
}

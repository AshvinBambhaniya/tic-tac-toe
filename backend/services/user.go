package services

import (
	"github.com/AshvinBambhaniya/tic-tac-toe/models"
	"github.com/AshvinBambhaniya/tic-tac-toe/utils"
	"github.com/google/uuid"
)

type UserService struct {
	userModel *models.UserModel
}

func NewUserService(userModel *models.UserModel) *UserService {
	return &UserService{
		userModel: userModel,
	}
}

func (userSvc *UserService) RegisterUser(user models.User) (models.User, error) {
	hashedPassword, err := utils.PasswordHash(user.Password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPassword

	user, err = userSvc.userModel.InsertUser(user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (userSvc *UserService) GetUser(userId uuid.UUID) (models.User, error) {
	user, err := userSvc.userModel.GetById(userId)
	return user, err
}

// Authenticate verify identity using email, and password.
// On successful validtion it'll return the user
func (userSvc *UserService) Authenticate(email, password string) (models.User, error) {
	user, err := userSvc.userModel.GetUserByEmail(email)
	if err != nil {
		return user, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return models.User{}, nil // Return empty user if password doesn't match, controller handles error
	}

	return user, nil
}

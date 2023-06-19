package service

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserService interface {
	Register(user *model.User) (model.User, error)
	Login(user *model.User) (token *string, err error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userService struct {
	userRepo     repo.UserRepository
	sessionsRepo repo.SessionRepository
}

func NewUserService(userRepository repo.UserRepository, sessionsRepo repo.SessionRepository) UserService {
	return &userService{userRepository, sessionsRepo}
}

func (s *userService) Register(user *model.User) (model.User, error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return *user, err
	}

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}

	user.CreatedAt = time.Now()

	newUser, err := s.userRepo.CreateUser(*user)
	if err != nil {
		return *user, err
	}

	return newUser, nil
}

func (s *userService) Login(user *model.User) (token *string, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if dbUser.ID == 0 || dbUser.Email == "" {
		return nil, errors.New("user not found")
	}

	if dbUser.Password != user.Password || dbUser.Email != user.Email {
		return nil, errors.New("wrong email or password")
	}

	terminationOfTheTime := time.Now().Add(20 * time.Minute)
	makeClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		Email: dbUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: terminationOfTheTime.Unix(),
		},
	})

	tokenString, err := makeClaims.SignedString(model.JwtKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *userService) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	category, result := s.userRepo.GetUserTaskCategory()
	if result != nil {
		return nil, result
	}

	return category, nil
}

package userservice

import (
	"GoGameApp/entity"
	"GoGameApp/param"
	"GoGameApp/pkg/password"
	"GoGameApp/pkg/richerror"
	"fmt"
)

// --- serivce defining part
type Repository interface {
	CreateUser(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	repo Repository
	auth AuthGenerator
}

func New(repo Repository, authGenerator AuthGenerator) Service {
	return Service{repo: repo, auth: authGenerator}
}

// --- service business logic part
func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
	const op = "userservice.Register"

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("can't hash password -> %w", err)
	}

	// create new user in storage
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    hashedPassword,
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return param.RegisterResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"request": req, "created_user": createdUser})
	}

	// return created user
	return param.RegisterResponse{
		User: param.UserInfo{
			ID:          createdUser.ID,
			Name:        createdUser.Name,
			PhoneNumber: createdUser.PhoneNumber,
		},
	}, nil
}

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userservice.Login"

	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"phone_number": req.PhoneNumber})
	}

	if !password.VerifyPassword(user.Password, req.Password) {
		return param.LoginResponse{}, fmt.Errorf("invalid credentials")

	}

	// jwt create token
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return param.LoginResponse{
		Tokens: param.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken},
		User: param.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber},
	}, nil
}

func (s Service) Profile(req param.ProfileRequest) (param.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return param.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"request": req})
	}

	return param.ProfileResponse{User: param.UserInfo{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
	}}, nil
}

package userservice

import (
	"GoGameApp/dto"
	"GoGameApp/entity"
	"GoGameApp/pkg/password"
	"GoGameApp/pkg/richerror"
	"fmt"
)


//--- serivce defining part
type Repository interface {
	CreateUser(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
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


//--- service business logic part
func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	const op = "userservice.Register"

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("can't hash password -> %w", err)
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
		return dto.RegisterResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"request": req, "created_user": createdUser})
	}

	// return created user
	return dto.RegisterResponse{
		User: dto.UserInfo{
			ID:          createdUser.ID,
			Name:        createdUser.Name,
			PhoneNumber: createdUser.PhoneNumber,
		},
	}, nil
}


func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"

	//TODO: separate existence and get user by phone number method
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return dto.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	if !password.VerifyPassword(user.Password, req.Password) {
		return dto.LoginResponse{}, fmt.Errorf("invalid credentials")

	}

	// jwt create token
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return dto.LoginResponse{
		Tokens: dto.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken},
		User: dto.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber},
	}, nil
}



func (s Service) Profile(req dto.ProfileRequest) (dto.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return dto.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"request": req})
	}

	return dto.ProfileResponse{User: dto.UserInfo{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
	}}, nil
}

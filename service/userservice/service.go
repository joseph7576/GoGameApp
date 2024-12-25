package userservice

import (
	"GoGameApp/entity"
	"GoGameApp/pkg/errmsg"
	"GoGameApp/pkg/password"
	"GoGameApp/pkg/phonenumber"
	"GoGameApp/pkg/richerror"
	"fmt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
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

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	const op = "userservice.Register"

	//TODO: verify phone number
	// validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessage(errmsg.ErrMsgInvalidPhoneNumber).WithMeta(map[string]any{"request": req})
	}

	// check uniqueness of phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, richerror.New(op).WithErr(err).
				WithMeta(map[string]any{"request": req})
		}

		if !isUnique {
			return RegisterResponse{}, richerror.New(op).WithMessage(errmsg.ErrMsgPhoneNumberNotUnique).
				WithKind(richerror.KindBadRequest).WithMeta(map[string]any{"request": req})
		}
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3")
	}

	// validate password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password length should be greater than 8")
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("can't hash password -> %w", err)
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
		return RegisterResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"request": req, "created_user": createdUser})
	}

	// return created user
	return RegisterResponse{
		User: UserInfo{
			ID:          createdUser.ID,
			Name:        createdUser.Name,
			PhoneNumber: createdUser.PhoneNumber,
		},
	}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	Tokens Tokens   `json:"tokens"`
	User   UserInfo `json:"user"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "userservice.Login"

	//TODO: separate existence and get user by phone number method
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	if !password.VerifyPassword(user.Password, req.Password) {
		return LoginResponse{}, fmt.Errorf("invalid credentials")

	}

	// jwt create token
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return LoginResponse{
		Tokens: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken},
		User: UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber},
	}, nil
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	User UserInfo `json:"user"`
}

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]any{"request": req})
	}

	return ProfileResponse{User: UserInfo{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
	}}, nil
}

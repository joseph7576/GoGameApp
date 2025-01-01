package matchingservice

import (
	"GoGameApp/entity"
	"GoGameApp/param"
	"GoGameApp/pkg/richerror"
	"time"
)

type Repo interface {
	AddToWaitingList(userID uint, category entity.Category) error
}

type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	repo   Repo
	config Config
}

func New(repo Repo,
	config Config) Service {
	return Service{
		repo:   repo,
		config: config,
	}
}

func (s Service) AddToWaitingList(req param.AddToWaitingListRequest) (param.AddToWaitingListResponse, error) {
	const op = "matchingservice.AddToWaitingList"

	err := s.repo.AddToWaitingList(req.UserID, req.Category)
	if err != nil {
		return param.AddToWaitingListResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return param.AddToWaitingListResponse{Timeout: s.config.WaitingTimeout}, nil
}

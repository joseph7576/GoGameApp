package backofficeuserservice

import "GoGameApp/entity"

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) ListAllUsers() ([]entity.User, error) {
	list := make([]entity.User, 0)

	list = append(list, entity.User{
		ID:          1,
		PhoneNumber: "fake",
		Name:        "fake",
		Password:    "fake",
		Role:        entity.AdminRole,
	})

	return list, nil
}

package uservalidator

import (
	"GoGameApp/dto"
	"GoGameApp/entity"
	"GoGameApp/pkg/errmsg"
	"GoGameApp/pkg/richerror"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	phoneNumberRegex = "^09[0-9]{9}$"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{
		repo: repo,
	}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	//TODO: config the params for validation
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),

		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`))),

		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)),
			validation.By(v.checkPhoneNumberUniqueness)),
	); err != nil {

		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				fieldErrors[key] = value.Error()
			}
		}

		return fieldErrors, richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessage(errmsg.ErrMsgInvalidInput).WithMeta(map[string]any{"request": req})
	}

	return nil, nil
}

func (v Validator) checkPhoneNumberUniqueness(value any) error {
	phoneNumber := value.(string)

	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {
		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrMsgPhoneNumberNotUnique)
		}
	}

	return nil
}

func (v Validator) ValidateLoginRequest(req dto.LoginRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)),
			validation.By(v.checkPhoneNumberExist)),

		validation.Field(&req.PhoneNumber, validation.Required),
	); err != nil {

		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				fieldErrors[key] = value.Error()
			}
		}

		return fieldErrors, richerror.New(op).WithKind(richerror.KindBadRequest).
			WithMessage(errmsg.ErrMsgInvalidInput).WithMeta(map[string]any{"request": req})
	}

	return nil, nil
}

func (v Validator) checkPhoneNumberExist(value any) error {
	phoneNumber := value.(string)

	_, err := v.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

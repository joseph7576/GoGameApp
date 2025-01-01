package matchingvalidator

import (
	"GoGameApp/entity"
	"GoGameApp/param"
	"GoGameApp/pkg/errmsg"
	"GoGameApp/pkg/richerror"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddToWaitingList(req param.AddToWaitingListRequest) (map[string]string, error) {
	const op = "matchingvalidator.ValidateAddToWaitingList"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Category, validation.Required,
			validation.By(v.isCategoryValid)),
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

func (v Validator) isCategoryValid(value any) error {
	category := value.(entity.Category)

	if !category.IsValid() {
		return fmt.Errorf(errmsg.ErrMsgCategoryNoValid)
	}

	return nil
}

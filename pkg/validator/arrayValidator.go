package validator

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ArrayValidator struct {
	elementsValidator FieldValidator
	Elements          json.RawMessage
	Required          bool
	MaxElements       *int
	MinElements       *int
}

func (v *ArrayValidator) Type() string {
	return "array"
}

func (v *ArrayValidator) IsRequired() bool {
	return v.Required
}

func (v *ArrayValidator) Validate(json interface{}, position string) error {
	jsonArray, ok := json.([]interface{})

	if !ok {
		return ValidationError{"This field is not an array", position}
	}

	compositeError := CompositeValidationError{}
	returnError := false
	for key, value := range jsonArray {
		fieldPosition := position + "/" + fmt.Sprint(key)
		fieldError := v.elementsValidator.Validate(value, fieldPosition)

		if fieldError != nil {
			compositeError.errors = append(compositeError.errors, fieldError)
			returnError = true
		}
	}

	if v.MaxElements != nil && len(jsonArray) > *v.MaxElements {
		compositeError.errors = append(compositeError.errors, ValidationError{
			fmt.Sprintf("This array can contain at most %d elements", *v.MaxElements),
			position,
		})
		returnError = true
	}
	if v.MinElements != nil && len(jsonArray) < *v.MinElements {
		compositeError.errors = append(compositeError.errors, ValidationError{
			fmt.Sprintf("This array must contain at least %d elements", *v.MinElements),
			position,
		})
		returnError = true
	}

	if returnError {
		return compositeError
	} else {
		return nil
	}
}

func (v *ArrayValidator) ValidatePatch(patch Patch, position string) error {
	if patch.IsRootPosition() {
		switch patch.op {
		case "remove":
			if v.IsRequired() {
				return ValidationError{"Cannot remove a required field", position}
			}

			return nil
		default:
			return v.Validate(patch.value, position)
		}
	}

	field, err := patch.UnshiftPosition()
	if err != nil {
		return ValidationError{err.Error(), position}
	}

	// Check that the path specifies an array element
	if field != "-" {
		value, err := strconv.Atoi(field)
		if err != nil {
			return ValidationError{"Array position is invalid, it should be either - or a number", position}
		}

		// FIXME: the length limit can be bypassed by using the position -
		if v.MaxElements != nil && *v.MaxElements <= value {
			return ValidationError{fmt.Sprintf("This array can contain at most %d elements", *v.MaxElements), position + "/" + field}
		}
	}

	return v.elementsValidator.ValidatePatch(patch, position+"/"+field)
}

func (v *ArrayValidator) InitializeAfterUnmarshaling(customTypes map[string]bool, rootValidator *Validator) error {
	validator, err := UnmarshalValidator(v.Elements, customTypes, rootValidator)

	if err != nil {
		return err
	}

	v.elementsValidator = validator

	return nil
}

package validator

import (
	"encoding/json"
)

type ObjectValidator struct {
	keyValidators map[string]FieldValidator
	Fields        map[string]json.RawMessage
	Required      bool
}

func (v *ObjectValidator) Type() string {
	return "object"
}

func (v *ObjectValidator) IsRequired() bool {
	return v.Required
}

func (v *ObjectValidator) Validate(json interface{}, position string) error {
	jsonObject, ok := json.(map[string]interface{})

	if !ok {
		return ValidationError{"This field is not an object", position}
	}

	compositeError := CompositeValidationError{}
	returnError := false
	for key, validator := range v.keyValidators {
		if value, ok := jsonObject[key]; !ok {
			if validator.IsRequired() {
				compositeError.errors = append(compositeError.errors, ValidationError{"Required field " + key + " is missing", position})
				returnError = true
			}
		} else {
			fieldPosition := position + "/" + key
			fieldError := validator.Validate(value, fieldPosition)

			if fieldError != nil {
				compositeError.errors = append(compositeError.errors, fieldError)
				returnError = true
			}
		}
	}

	for key := range jsonObject {
		if _, ok := v.keyValidators[key]; !ok {
			compositeError.errors = append(compositeError.errors, ValidationError{"Field " + key + " is not specified in the schema", position})
			returnError = true
		}
	}

	if returnError {
		return compositeError
	} else {
		return nil
	}
}

// TODO: change validation depending on Patch operation
func (v *ObjectValidator) ValidatePatch(patch Patch, position string) error {
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

	if validator, ok := v.keyValidators[field]; ok {
		fieldPosition := position + "/" + field
		return validator.ValidatePatch(patch, fieldPosition)
	} else {
		return ValidationError{"Field " + field + " is not specified in the schema", position}
	}
}

func (v *ObjectValidator) InitializeAfterUnmarshaling() error {
	v.keyValidators = make(map[string]FieldValidator, len(v.Fields))

	for key, value := range v.Fields {
		validator, err := UnmarshalValidator(value)

		if err != nil {
			return err
		}

		v.keyValidators[key] = validator
	}

	return nil
}

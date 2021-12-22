package validator

import (
	"encoding/json"
	"fmt"
)

type FieldValidator interface {
	Type() string
	Validate(interface{}, string) error
	ValidatePatch(Patch, string) error
	InitializeAfterUnmarshaling() error
	IsRequired() bool
}

type Validator struct {
	fieldValidator FieldValidator
}

func (v *Validator) UnmarshalJSON(data []byte) error {
	validator, err := UnmarshalValidator(data)

	if err != nil {
		return err
	}

	v.fieldValidator = validator
	return nil
}

func (v *Validator) Validate(jsonSource []byte) error {
	var decodedJson interface{}
	err := json.Unmarshal([]byte(jsonSource), &decodedJson)

	if err != nil {
		return fmt.Errorf("Failed to parse json:\n" + err.Error())
	}

	return v.fieldValidator.Validate(decodedJson, "$")
}

func (v *Validator) ValidatePatches(jsonPatches []byte) error {
	var patches []Patch
	err := json.Unmarshal([]byte(jsonPatches), &patches)
	if err != nil {
		return fmt.Errorf("Could not parse json patches: " + err.Error())
	}

	compositeError := CompositeValidationError{}
	returnError := false
	for id, patch := range patches {
		err := v.fieldValidator.ValidatePatch(patch, fmt.Sprintf("Patch %d ", id))

		if err != nil {
			compositeError.errors = append(compositeError.errors, err)
			returnError = true
		}
	}

	if returnError {
		return compositeError
	} else {
		return nil
	}
}

func UnmarshalValidator(data []byte) (FieldValidator, error) {
	var validatorType struct {
		Type string `json:"type"`
	}

	err := json.Unmarshal(data, &validatorType)
	if err != nil {
		return nil, fmt.Errorf("Passed json is invalid: \n" + err.Error())
	}

	var validator FieldValidator
	switch validatorType.Type {
	case "object":
		var objectValidator ObjectValidator
		json.Unmarshal(data, &objectValidator)
		validator = &objectValidator
		break
	case "array":
		var arrayValidator ArrayValidator
		json.Unmarshal(data, &arrayValidator)

		validator = &arrayValidator
		break
	case "string":
		var stringValidator StringValidator
		json.Unmarshal(data, &stringValidator)

		validator = &stringValidator
		break
	case "float":
		var floatValidator FloatValidator
		json.Unmarshal(data, &floatValidator)
		validator = &floatValidator
		break
	case "any":
		var anyValidator AnyValidator
		json.Unmarshal(data, &anyValidator)
		validator = &anyValidator
		break
	default:
		return nil, fmt.Errorf("Passed validation schema is missing the type field or the type is not supported")
	}

	err = validator.InitializeAfterUnmarshaling()
	if err != nil {
		return nil, err
	}

	return validator, nil
}

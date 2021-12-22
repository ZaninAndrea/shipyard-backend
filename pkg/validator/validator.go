package validator

import (
	"encoding/json"
	"fmt"
)

type FieldValidator interface {
	Type() string
	Validate(interface{}, string) error
	ValidatePatch(Patch, string) error
	InitializeAfterUnmarshaling(map[string]bool, *Validator) error
	IsRequired() bool
}

type Validator struct {
	fieldValidator FieldValidator
	customTypes map[string]FieldValidator
}

type CustomValidator struct{
	sourceValidator *Validator
	fieldName string
}

func (v *CustomValidator) Type() string{
	return v.fieldName
}

func (v *CustomValidator) Validate(json interface{}, position string) error {
	return v.sourceValidator.customTypes[v.fieldName].Validate(json, position)
}

func (v *CustomValidator) ValidatePatch(patch Patch, position string) error {
	return v.sourceValidator.customTypes[v.fieldName].ValidatePatch(patch, position)
}

func (v *CustomValidator) InitializeAfterUnmarshaling(customTypes map[string]bool, rootValidator *Validator) error {
	return nil
}

func (v *CustomValidator) IsRequired() bool {
	return v.sourceValidator.customTypes[v.fieldName].IsRequired()
}

func (v *Validator) UnmarshalJSON(data []byte) error {
	// Unmarshal custom types
	var validatorType struct {
		CustomTypes map[string]json.RawMessage `json:"customTypes"`
	}
	err := json.Unmarshal(data, &validatorType)
	
	typesList := make(map[string]bool)
	for key := range validatorType.CustomTypes{
		typesList[key] = true
	}

	if err != nil {
		return fmt.Errorf("Passed json is invalid: \n" + err.Error())
	}

	v.customTypes = make(map[string]FieldValidator)
	for key, source := range validatorType.CustomTypes{
		keyValidator, err := UnmarshalValidator(source, typesList, v)
		if err != nil{
			return err
		}

		v.customTypes[key] = keyValidator
	}

	// Unmarshal validator
	validator, err := UnmarshalValidator(data, typesList, v)

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

func UnmarshalValidator(data []byte, customTypes map[string]bool, rootValidator *Validator) (FieldValidator, error) {
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
	case "array":
		var arrayValidator ArrayValidator
		json.Unmarshal(data, &arrayValidator)

		validator = &arrayValidator
	case "string":
		var stringValidator StringValidator
		json.Unmarshal(data, &stringValidator)

		validator = &stringValidator
	case "float":
		var floatValidator FloatValidator
		json.Unmarshal(data, &floatValidator)
		validator = &floatValidator
	case "any":
		var anyValidator AnyValidator
		json.Unmarshal(data, &anyValidator)
		validator = &anyValidator
	default:
		if _, ok := customTypes[validatorType.Type]; ok{
			validator = &CustomValidator{sourceValidator: rootValidator, fieldName: validatorType.Type}
		}else{
			return nil, fmt.Errorf("Passed validation schema is missing the type field or the type (%s) is not supported", validatorType.Type)
		}
	}

	err = validator.InitializeAfterUnmarshaling(customTypes, rootValidator)
	if err != nil {
		return nil, err
	}

	return validator, nil
}

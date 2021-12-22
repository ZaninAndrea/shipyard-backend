package validator

import "fmt"

type StringValidator struct {
	Required bool
	MaxChars *int
	MinChars *int
}

func (v *StringValidator) Type() string {
	return "string"
}

func (v *StringValidator) IsRequired() bool {
	return v.Required
}

func (v *StringValidator) Validate(json interface{}, position string) error {
	jsonString, ok := json.(string)

	if !ok {
		return ValidationError{"This field is not an string", position}
	}

	if v.MaxChars != nil && len(jsonString) > *v.MaxChars {
		return ValidationError{fmt.Sprintf("This string is longer than maxChars (%d)", *v.MaxChars), position}
	}
	if v.MinChars != nil && len(jsonString) < *v.MinChars {
		return ValidationError{fmt.Sprintf("This string is shorter than minChars (%d)", *v.MinChars), position}
	}

	return nil
}

func (v *StringValidator) ValidatePatch(patch Patch, position string) error {
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
	} else {
		return ValidationError{"Cannot access a field inside a string", position}
	}
}

func (v *StringValidator) InitializeAfterUnmarshaling() error {
	return nil
}

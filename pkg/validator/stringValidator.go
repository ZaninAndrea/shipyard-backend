package validator

import (
	"fmt"
	"regexp"
)

type StringValidator struct {
	Required      bool
	MaxChars      *int
	MinChars      *int
	RegexMatch    *string // checks that the string matches the given regex
	NoRegexMatch  *string // checks that the string doesn't match the given regex
	AllowedValues *[]string
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
	if v.RegexMatch != nil {
		r, err := regexp.Compile(*v.RegexMatch)
		if err != nil {
			return ValidationError{"The validator for this field has an invalid RegexMatch", position}
		}

		if !r.MatchString(jsonString) {
			return ValidationError{fmt.Sprintf("This string does not match the RegexMatch field: %s", *v.RegexMatch), position}
		}
	}
	if v.NoRegexMatch != nil {
		r, err := regexp.Compile(*v.NoRegexMatch)
		if err != nil {
			return ValidationError{"The validator for this field has an invalid NoRegexMatch", position}
		}

		if r.MatchString(jsonString) {
			return ValidationError{fmt.Sprintf("This string matches the NoRegexMatch field: %s", *v.NoRegexMatch), position}
		}
	}
	if v.AllowedValues != nil {
		found := false
		for _, allowedValue := range *v.AllowedValues {
			if allowedValue == jsonString {
				found = true
				break
			}
		}

		if !found {
			return ValidationError{"The value is not in the AllowedValues", position}
		}
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

func (v *StringValidator) InitializeAfterUnmarshaling(customTypes map[string]bool, rootValidator *Validator) error {
	return nil
}

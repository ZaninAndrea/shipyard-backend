package validator

import (
	"encoding/json"
	"testing"
)

func TestEmptyStringValidation(t *testing.T) {
	const jsonValidator = `{
		"type": "string"
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	t.Run("Valid String", func(t *testing.T) {
		err = v.Validate([]byte(`"afd asf ajsdkf"`))
		if err != nil {
			t.Error(err)
		}
	})
}

func TestMinMaxStringValidation(t *testing.T) {
	const jsonValidator = `{
		"type": "string",
		"minChars": 3,
		"maxChars": 10
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	t.Run("Valid String", func(t *testing.T) {
		err = v.Validate([]byte(`"abc"`))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Invalid String (under min)", func(t *testing.T) {
		err = v.Validate([]byte(`"a"`))
		if err == nil {
			t.Error("The value -1 is below the min but wasn't rejected")
		}
	})
	t.Run("Invalid String (above max)", func(t *testing.T) {
		err = v.Validate([]byte(`"abcdefghijklmopqrstuv"`))
		if err == nil {
			t.Error("The value 12 is above the max but wasn't rejected")
		}
	})
}

func TestRegexMatchStringValidation(t *testing.T) {
	const jsonValidator = `{
		"type": "string",
		"regexMatch": "^\\d{9}$"
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	t.Run("Valid String", func(t *testing.T) {
		err = v.Validate([]byte(`"123456789"`))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Invalid String (doesn't match RegexMatch)", func(t *testing.T) {
		err = v.Validate([]byte(`"abc"`))
		if err == nil {
			t.Error("The value abc doesn't match the regex but wasn't rejected")
		}
	})
}

func TestNoRegexMatchStringValidation(t *testing.T) {
	const jsonValidator = `{
		"type": "string",
		"noRegexMatch": "^\\d{9}$"
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	t.Run("Invalid String", func(t *testing.T) {
		err = v.Validate([]byte(`"123456789"`))
		if err == nil {
			t.Error("The string 123456789 matches NoRegexMatch but wan't rejected")
		}
	})
	t.Run("Valid string", func(t *testing.T) {
		err = v.Validate([]byte(`"abc"`))
		if err != nil {
			t.Error(err)
		}
	})
}

func TestAllowedValuesStringValidation(t *testing.T) {
	const jsonValidator = `{
		"type": "string",
		"allowedValues": ["a", "b", "c"]
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	t.Run("Invalid String", func(t *testing.T) {
		err = v.Validate([]byte(`"z"`))
		if err == nil {
			t.Error("The string z is not in the AllowedValues but wan't rejected")
		}
	})
	t.Run("Valid string", func(t *testing.T) {
		err = v.Validate([]byte(`"a"`))
		if err != nil {
			t.Error(err)
		}
	})
}

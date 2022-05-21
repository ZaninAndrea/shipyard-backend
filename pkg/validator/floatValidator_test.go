package validator

import (
	"encoding/json"
	"testing"
)

func TestEmptyFloatValidation(t *testing.T) {
	const jsonValidator = `{
		"type": "float"
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	t.Run("Valid Float", func(t *testing.T) {
		err = v.Validate([]byte("5"))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Invalid Float (string)", func(t *testing.T) {
		err = v.Validate([]byte("a"))
		if err == nil {
			t.Error("The string 'a' should not be recognized as a valid float")
		}
	})
	t.Run("Invalid Float (incorrect numeric format)", func(t *testing.T) {
		err = v.Validate([]byte("5.2.3"))
		if err == nil {
			t.Error("The string '5.2.3' should not be recognized as a valid float")
		}
	})
}

func TestMinMaxFloatValidation(t *testing.T) {
	const jsonValidator = `{
		"type": "float",
		"min": 0,
		"max": 10
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	t.Run("Valid Float", func(t *testing.T) {
		err = v.Validate([]byte("5"))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Invalid Float (under min)", func(t *testing.T) {
		err = v.Validate([]byte("-1"))
		if err == nil {
			t.Error("The value -1 is below the min but wasn't rejected")
		}
	})
	t.Run("Invalid Float (above max)", func(t *testing.T) {
		err = v.Validate([]byte("12"))
		if err == nil {
			t.Error("The value 12 is above the max but wasn't rejected")
		}
	})

}

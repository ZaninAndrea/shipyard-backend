package validator

type FloatValidator struct {
	Required  bool
	Min       *float64 // Enforces a >= constraint
	StrictMin *float64 // Enforces a > constraint
	Max       *float64 // Enforces a <= constraint
	StrictMax *float64 // Enforces a <
}

func (v *FloatValidator) Type() string {
	return "float"
}

func (v *FloatValidator) IsRequired() bool {
	return v.Required
}

func (v *FloatValidator) Validate(json interface{}, position string) error {
	value, ok := json.(float64)

	if !ok {
		return ValidationError{"This field is not a float", position}
	}

	if v.Min != nil && value < *v.Min {
		return ValidationError{"The value is below the Min", position}
	}
	if v.Max != nil && value > *v.Max {
		return ValidationError{"The value is above the Max", position}
	}
	if v.StrictMin != nil && value <= *v.StrictMin {
		return ValidationError{"The value is below or equal to the StrictMin", position}
	}
	if v.StrictMax != nil && value > *v.StrictMax {
		return ValidationError{"The value is above or equal to the StrictMax", position}
	}

	return nil
}

func (v *FloatValidator) ValidatePatch(patch Patch, position string) error {
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
		return ValidationError{"Cannot access a field inside a float", position}
	}
}

func (v *FloatValidator) InitializeAfterUnmarshaling(customTypes map[string]bool, rootValidator *Validator) error {
	return nil
}

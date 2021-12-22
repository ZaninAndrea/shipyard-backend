package validator

type FloatValidator struct {
	Required bool
}

func (v *FloatValidator) Type() string {
	return "float"
}

func (v *FloatValidator) IsRequired() bool {
	return v.Required
}

func (v *FloatValidator) Validate(json interface{}, position string) error {
	_, ok := json.(float64)

	if !ok {
		return ValidationError{"This field is not a float", position}
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

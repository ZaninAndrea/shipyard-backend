package validator

type AnyValidator struct {
	Required bool
}

func (v *AnyValidator) Type() string {
	return "any"
}

func (v *AnyValidator) IsRequired() bool {
	return v.Required
}

func (v *AnyValidator) Validate(json interface{}, position string) error {
	return nil
}

func (v *AnyValidator) ValidatePatch(patch Patch, position string) error {
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

	return nil
}

func (v *AnyValidator) InitializeAfterUnmarshaling(customTypes map[string]bool, rootValidator *Validator) error {
	return nil
}

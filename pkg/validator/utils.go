package validator

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ValidationError struct {
	message  string
	position string
}

func (e ValidationError) Error() string {
	return "[" + e.position + "] " + e.message
}

type CompositeValidationError struct {
	errors []error
}

func (e CompositeValidationError) Error() string {
	message := ""

	for _, err := range e.errors {
		if message != "" {
			message += "\n" + err.Error()
		} else {
			message += err.Error()
		}
	}

	return message
}

type Patch struct {
	op    string
	path  string
	value interface{}
}

func (p *Patch) UnmarshalJSON(data []byte) error {
	var patchMap map[string]interface{}
	err := json.Unmarshal(data, &patchMap)

	if err != nil {
		return fmt.Errorf("Failed to parse patches json, one or more of the patches is not an object")
	}

	if value, ok := patchMap["op"]; ok {
		opString, ok := value.(string)

		if !ok {
			return fmt.Errorf("Patch field \"op\" is not a string")
		}

		p.op = opString
	} else {
		return fmt.Errorf("Patch field \"op\" not specified")
	}

	if value, ok := patchMap["path"]; ok {
		pathString, ok := value.(string)

		if !ok {
			return fmt.Errorf("Patch field \"path\" is not a string")
		}

		p.path = pathString
	} else {
		return fmt.Errorf("Patch field \"path\" not specified")
	}

	if value, ok := patchMap["value"]; ok {
		p.value = value
	}

	return nil
}

func (p *Patch) IsRootPosition() bool {
	return p.path == ""
}

func (p *Patch) UnshiftPosition() (string, error) {
	if p.IsRootPosition() {
		panic("UnshiftPosition should not be called if patch is at root position")
	}

	if p.path[0] != '/' {
		return "", fmt.Errorf("Invalid patch path")
	}

	p.path = p.path[1:]
	parts := strings.SplitAfterN(p.path, "/", 2)

	switch len(parts) {
	// Case 0 is when the original path was /
	case 0:
		p.path = ""
		return "", nil

	// Case 1 is when the original path was /something
	case 1:
		p.path = ""
		return parts[0], nil

	// Case 2 is when the original path was /something/somethingelse
	case 2:
		s := strings.Replace(parts[0], "/", "", 1)
		p.path = "/" + parts[1]

		return s, nil
	}

	panic("This should never be reached")
}

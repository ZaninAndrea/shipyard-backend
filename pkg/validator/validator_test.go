package validator

import (
	"encoding/json"
	"testing"
)

func TestSourceValidation(t *testing.T) {
	const jsonSource string = `{
		"name": "Andrea",
		"email": "andrea@igloo.ooo",
		"decks":[
			{
				"name": "Analisi",
				"repetitionCount": 59
			}
		]
	}
	`

	const jsonValidator = `{
		"type": "object",
		"fields": {
			"name": {
				"type": "string",
				"required": false,
				"maxChars": 20
			},
			"email": {
				"type": "string",
				"required": true
			},
			"decks":{
				"type":"array",
				"maxElements": 1,
				"elements":{
					"type":"object",
					"fields":{
						"name": {
							"type": "string"
						},
						"repetitionCount":{
							"type": "float"
						}
					}
				}
			}
		}	
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	err = v.Validate([]byte(jsonSource))
	if err != nil {
		t.Error(err)
	}
}

func TestPatchValidation(t *testing.T) {
	const patchSource string = `[
		{ "op": "add", "path": "/name", "value": "Giorgio" },
		{ "op": "replace", "path": "/decks/0/repetitionCount", "value": 5 },
		{ "op": "remove", "path": "/name" }
	]`

	const jsonValidator = `{
		"type": "object",
		"fields": {
			"name": {
				"type": "string",
				"required": false,
				"maxChars": 20
			},
			"email": {
				"type": "string",
				"required": true
			},
			"decks":{
				"type":"array",
				"maxElements": 1,
				"elements":{
					"type":"object",
					"fields":{
						"name": {
							"type": "string"
						},
						"repetitionCount":{
							"type": "float"
						}
					}
				}
			}
		}	
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	err = v.ValidatePatches([]byte(patchSource))
	if err != nil {
		t.Error(err)
	}
}

func TestPatchStringValidation(t *testing.T) {
	const patchSource string = `[
		{ "op": "add", "path": "/name", "value": "Giorgio asdf asdf asdf asdf" }
	]`

	const jsonValidator = `{
		"type": "object",
		"fields": {
			"name": {
				"type": "string",
				"required": false,
				"maxChars": 20
			}
		}	
	}`

	var v Validator
	err := json.Unmarshal([]byte(jsonValidator), &v)
	if err != nil {
		panic(err)
	}

	err = v.ValidatePatches([]byte(patchSource))
	if err == nil {
		t.Error("Didn't throw the string length over maxChars error")
	}
}

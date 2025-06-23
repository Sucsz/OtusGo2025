package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
)

// createStringValidator creates a validator for string fields based on the validation rule.
// Supported rules: "len", "regexp", "in".
// Returns an error if the rule parameters are invalid or the rule is unknown.
func createStringValidator(rule ValidatorRule) (Validator[string], error) {
	switch rule.Name {
	case "len":
		length, err := strconv.Atoi(rule.Param)
		if err != nil {
			return nil, fmt.Errorf("invalid len parameter: %w", err)
		}
		return LengthValidator{length}, nil

	case "regexp":
		re, err := regexp.Compile(rule.Param)
		if err != nil {
			return nil, fmt.Errorf("invalid regexp: %w", err)
		}
		return RegexpValidator{re}, nil

	case "in":
		return InValidator[string]{rule.Params}, nil

	default:
		return nil, fmt.Errorf("unknown validator %s", rule.Name)
	}
}

// createIntValidator creates a validator for integer fields based on the validation rule.
// Supported rules: "min", "max", "in".
// Returns an error if the rule parameters are invalid or the rule is unknown.
func createIntValidator(rule ValidatorRule) (Validator[int], error) {
	switch rule.Name {
	case "min":
		minP, err := strconv.Atoi(rule.Param)
		if err != nil {
			return nil, fmt.Errorf("invalid min parameter: %w", err)
		}
		return MinValidator{minP}, nil

	case "max":
		maxP, err := strconv.Atoi(rule.Param)
		if err != nil {
			return nil, fmt.Errorf("invalid max parameter: %w", err)
		}
		return MaxValidator{maxP}, nil

	case "in":
		// Convert all allowed values to integers and create inclusion validator.
		var allowed []int
		for _, param := range rule.Params {
			num, err := strconv.Atoi(param)
			if err != nil {
				return nil, fmt.Errorf("invalid in parameter: %w", err)
			}
			allowed = append(allowed, num)
		}
		return InValidator[int]{allowed}, nil

	default:
		return nil, fmt.Errorf("unknown validator %s", rule.Name)
	}
}

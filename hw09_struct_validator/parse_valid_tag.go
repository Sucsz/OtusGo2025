package hw09structvalidator

import (
	"reflect"
	"strings"
)

// ValidatorRule represents a single validation rule with its parameters.
type ValidatorRule struct {
	Name   string
	Param  string
	Params []string
}

// ParseValidTag parses the "validate" struct tag and returns a slice of ValidatorRules.
// The tag format is "rule1:param1|rule2:param2,param3|...".
func ParseValidTag(tag reflect.StructTag) ([]ValidatorRule, error) {
	validTag, ok := tag.Lookup("validate")
	if !ok {
		return nil, nil
	}

	validTag = strings.TrimSpace(validTag)
	if validTag == "" {
		return nil, nil
	}

	rules := strings.Split(validTag, "|")
	validatorRules := make([]ValidatorRule, 0, len(rules))

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		// Split rule into name and parameters
		parts := strings.SplitN(rule, ":", 2)
		validatorRule := ValidatorRule{
			Name: strings.TrimSpace(parts[0]),
		}

		// If rule has parameters
		if len(parts) == 2 {
			validatorRule.Param = strings.TrimSpace(parts[1])
			if validatorRule.Name == "in" {
				validatorRule.Params = strings.Split(validatorRule.Param, ",")
				for i := range validatorRule.Params {
					validatorRule.Params[i] = strings.TrimSpace(validatorRule.Params[i])
				}
			}
		}

		validatorRules = append(validatorRules, validatorRule)
	}

	return validatorRules, nil
}

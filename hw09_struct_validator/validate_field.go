package hw09structvalidator

import (
	"fmt"
	"reflect"
)

// ValidateField validates a struct field against given validation rules.
// It returns a slice of ValidationError if validation fails, or an error if validation setup fails.
// Supports string, int, and slices of these types.
//
//nolint:gocognit
func ValidateField(fieldName string, fieldValue reflect.Value, rules []ValidatorRule) ([]ValidationError, error) {
	var errors []ValidationError

	if !fieldValue.IsValid() {
		return nil, nil
	}

	//nolint:exhaustive
	switch fieldValue.Kind() {
	case reflect.String:
		value := fieldValue.String()
		for _, rule := range rules {
			validator, err := createStringValidator(rule)
			if err != nil {
				return nil, err
			}
			if err := validator.Validate(value); err != nil {
				errors = append(errors, ValidationError{fieldName, err})
			}
		}

	case reflect.Int:
		value := int(fieldValue.Int())
		for _, rule := range rules {
			validator, err := createIntValidator(rule)
			if err != nil {
				return nil, err
			}
			if err := validator.Validate(value); err != nil {
				errors = append(errors, ValidationError{fieldName, err})
			}
		}

	//nolint:exhaustive
	case reflect.Slice:
		elemType := fieldValue.Type().Elem()
		switch elemType.Kind() {
		case reflect.String:
			slice := fieldValue.Interface().([]string)
			for _, rule := range rules {
				validator, err := createStringValidator(rule)
				if err != nil {
					return nil, err
				}
				for i, item := range slice {
					if err := validator.Validate(item); err != nil {
						errors = append(errors, ValidationError{
							Field: fmt.Sprintf("%s[%d]", fieldName, i),
							Err:   err,
						})
					}
				}
			}

		case reflect.Int:
			slice := fieldValue.Interface().([]int)
			for _, rule := range rules {
				validator, err := createIntValidator(rule)
				if err != nil {
					return nil, err
				}
				for i, item := range slice {
					if err := validator.Validate(item); err != nil {
						errors = append(errors, ValidationError{
							Field: fmt.Sprintf("%s[%d]", fieldName, i),
							Err:   err,
						})
					}
				}
			}

		// Return error if trying to validate unsupported slice type.
		default:
			if len(rules) > 0 {
				return nil, fmt.Errorf("unsupported slice element type: %s", elemType.Kind().String())
			}
		}

	// Return error if trying to validate unsupported type.
	default:
		if len(rules) > 0 {
			return nil, fmt.Errorf("unsupported field type: %s", fieldValue.Kind().String())
		}
	}

	return errors, nil
}

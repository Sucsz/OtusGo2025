package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("%s: %v\n", err.Field, err.Err))
	}
	return sb.String()
}

type Validator[T any] interface {
	Validate(value T) error
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldValue := val.Field(i)
		rules, err := ParseValidTag(field.Tag)
		if err != nil {
			return fmt.Errorf("tag parsing error: %w", err)
		}
		if len(rules) == 0 {
			continue
		}

		errs, err := ValidateField(field.Name, fieldValue, rules)
		if err != nil {
			return err
		}
		validationErrors = append(validationErrors, errs...)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

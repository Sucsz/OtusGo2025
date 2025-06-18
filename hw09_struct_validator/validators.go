package hw09structvalidator

import (
	"fmt"
	"regexp"
)

// LengthValidator validates that a string has exact required length.
type LengthValidator struct {
	Length int
}

func (v LengthValidator) Validate(value string) error {
	if len(value) != v.Length {
		return fmt.Errorf("length must be %d", v.Length)
	}
	return nil
}

// RegexpValidator validates string against regular expression pattern.
type RegexpValidator struct {
	Pattern *regexp.Regexp
}

func (v RegexpValidator) Validate(value string) error {
	if !v.Pattern.MatchString(value) {
		return fmt.Errorf("must match pattern %s", v.Pattern.String())
	}
	return nil
}

// MinValidator validates that integer is greater or equal to minimum value.
type MinValidator struct {
	Min int
}

func (v MinValidator) Validate(value int) error {
	if value < v.Min {
		return fmt.Errorf("value must be >= %d", v.Min)
	}
	return nil
}

// MaxValidator validates that integer is less or equal to maximum value.
type MaxValidator struct {
	Max int
}

func (v MaxValidator) Validate(value int) error {
	if value > v.Max {
		return fmt.Errorf("value must be <= %d", v.Max)
	}
	return nil
}

// InValidator validates that value is in allowed set of values.
// Generic type T must be comparable (supports == and != operations).
type InValidator[T comparable] struct {
	AllowedValues []T
}

func (v InValidator[T]) Validate(value T) error {
	for _, allowed := range v.AllowedValues {
		if value == allowed {
			return nil
		}
	}
	return fmt.Errorf("value not in allowed set: %v", v.AllowedValues)
}

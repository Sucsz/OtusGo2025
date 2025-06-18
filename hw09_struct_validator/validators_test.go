package hw09structvalidator

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLengthValidator(t *testing.T) {
	t.Run("valid length", func(t *testing.T) {
		v := LengthValidator{5}
		require.NoError(t, v.Validate("abcde"))
	})

	t.Run("invalid length", func(t *testing.T) {
		v := LengthValidator{5}
		require.EqualError(t, v.Validate("abcd"), "length must be 5")
	})
}

func TestRegexpValidator(t *testing.T) {
	t.Run("valid pattern", func(t *testing.T) {
		v := RegexpValidator{regexp.MustCompile(`^\d+$`)}
		require.NoError(t, v.Validate("123"))
	})

	t.Run("invalid pattern", func(t *testing.T) {
		v := RegexpValidator{regexp.MustCompile(`^\d+$`)}
		err := v.Validate("abc")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must match pattern")
	})
}

func TestInValidator(t *testing.T) {
	t.Run("valid string value", func(t *testing.T) {
		v := InValidator[string]{[]string{"a", "b", "c"}}
		require.NoError(t, v.Validate("b"))
	})

	t.Run("invalid string value", func(t *testing.T) {
		v := InValidator[string]{[]string{"a", "b", "c"}}
		err := v.Validate("d")
		require.Error(t, err)
		require.Contains(t, err.Error(), "value not in allowed set")
	})

	t.Run("valid int value", func(t *testing.T) {
		v := InValidator[int]{[]int{1, 2, 3}}
		require.NoError(t, v.Validate(2))
	})
}

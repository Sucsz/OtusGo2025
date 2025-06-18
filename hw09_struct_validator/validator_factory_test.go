package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateStringValidator(t *testing.T) {
	t.Run("len validator", func(t *testing.T) {
		v, err := createStringValidator(ValidatorRule{Name: "len", Param: "5"})
		require.NoError(t, err)
		require.IsType(t, LengthValidator{}, v)
	})

	t.Run("invalid len parameter", func(t *testing.T) {
		_, err := createStringValidator(ValidatorRule{Name: "len", Param: "abc"})
		require.Error(t, err)
	})
}

func TestCreateIntValidator(t *testing.T) {
	t.Run("min validator", func(t *testing.T) {
		v, err := createIntValidator(ValidatorRule{Name: "min", Param: "5"})
		require.NoError(t, err)
		require.IsType(t, MinValidator{}, v)
	})

	t.Run("invalid min parameter", func(t *testing.T) {
		_, err := createIntValidator(ValidatorRule{Name: "min", Param: "abc"})
		require.Error(t, err)
	})
}

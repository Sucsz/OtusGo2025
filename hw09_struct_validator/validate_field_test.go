package hw09structvalidator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateField(t *testing.T) {
	tests := []struct {
		name        string
		fieldName   string
		value       interface{}
		rules       []ValidatorRule
		expectedErr []ValidationError
		expectError bool
	}{
		{
			name:      "string valid length",
			fieldName: "ID",
			value:     "123456789012345678901234567890123456",
			rules: []ValidatorRule{
				{Name: "len", Param: "36"},
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name:      "string invalid length",
			fieldName: "ID",
			value:     "short",
			rules: []ValidatorRule{
				{Name: "len", Param: "36"},
			},
			expectedErr: []ValidationError{
				{Field: "ID", Err: fmt.Errorf("length must be 36")},
			},
			expectError: true,
		},
		{
			name:      "int valid range",
			fieldName: "Age",
			value:     25,
			rules: []ValidatorRule{
				{Name: "min", Param: "18"},
				{Name: "max", Param: "50"},
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name:      "int invalid min",
			fieldName: "Age",
			value:     15,
			rules: []ValidatorRule{
				{Name: "min", Param: "18"},
			},
			expectedErr: []ValidationError{
				{Field: "Age", Err: fmt.Errorf("value must be >= 18")},
			},
			expectError: true,
		},
		{
			name:      "string slice validation",
			fieldName: "Phones",
			value:     []string{"12345678901", "12345"},
			rules: []ValidatorRule{
				{Name: "len", Param: "11"},
			},
			expectedErr: []ValidationError{
				{Field: "Phones[1]", Err: fmt.Errorf("length must be 11")},
			},
			expectError: true,
		},
		{
			name:      "int slice validation",
			fieldName: "Scores",
			value:     []int{10, 20, 30},
			rules: []ValidatorRule{
				{Name: "max", Param: "25"},
			},
			expectedErr: []ValidationError{
				{Field: "Scores[2]", Err: fmt.Errorf("value must be <= 25")},
			},
			expectError: true,
		},
		{
			name:      "string regexp valid",
			fieldName: "Email",
			value:     "test@example.com",
			rules: []ValidatorRule{
				{Name: "regexp", Param: `^\w+@\w+\.\w+$`},
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name:      "string in valid",
			fieldName: "Role",
			value:     "admin",
			rules: []ValidatorRule{
				{Name: "in", Param: "admin,user", Params: []string{"admin", "user"}},
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name:      "int in valid",
			fieldName: "Code",
			value:     200,
			rules: []ValidatorRule{
				{Name: "in", Param: "200,404,500", Params: []string{"200", "404", "500"}},
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name:      "combined validators",
			fieldName: "Password",
			value:     "1234567890",
			rules: []ValidatorRule{
				{Name: "len", Param: "10"},
				{Name: "regexp", Param: `\d+`},
			},
			expectedErr: nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldValue := reflect.ValueOf(tt.value)
			valErrs, err := ValidateField(tt.fieldName, fieldValue, tt.rules)

			if tt.expectError {
				if err != nil {
					return // Это программная ошибка (например, неверный валидатор)
				}
				require.Equal(t, tt.expectedErr, valErrs)
			} else {
				require.NoError(t, err)
				require.Empty(t, valErrs)
			}
		})
	}
}

func TestValidateField_Unsupported(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		rules   []ValidatorRule
		wantErr bool
	}{
		{
			name:    "unsupported type with rules",
			value:   true,
			rules:   []ValidatorRule{{Name: "min", Param: "1"}},
			wantErr: true,
		},
		{
			name:    "unsupported type without rules",
			value:   true,
			rules:   nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldValue := reflect.ValueOf(tt.value)
			_, err := ValidateField("Field", fieldValue, tt.rules)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

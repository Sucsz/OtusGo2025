package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	StringRegexpLenCombined struct {
		Value string `validate:"regexp:\\d+|len:20"`
	}
)

//nolint:funlen
func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "short_id",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be 36")},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    15,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("value must be >= 18")},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    55,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("value must be <= 50")},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    25,
				Email:  "invalid-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: fmt.Errorf("must match pattern ^\\w+@\\w+\\.\\w+$")},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    25,
				Email:  "test@example.com",
				Role:   "invalid_role",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Role", Err: fmt.Errorf("value not in allowed set: [admin stuff]")},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"short"},
			},
			expectedErr: ValidationErrors{
				{Field: "Phones[0]", Err: fmt.Errorf("length must be 11")},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 400,
				Body: "Bad Request",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("value not in allowed set: [200 404 500]")},
			},
		},
		{
			in: StringRegexpLenCombined{
				Value: "12345678901234567890",
			},
			expectedErr: nil,
		},
		{
			in: StringRegexpLenCombined{
				Value: "not_digits",
			},
			expectedErr: ValidationErrors{
				{Field: "Value", Err: fmt.Errorf("must match pattern \\d+")},
				{Field: "Value", Err: fmt.Errorf("length must be 20")},
			},
		},
		{
			in: StringRegexpLenCombined{
				Value: "1234567890",
			},
			expectedErr: ValidationErrors{
				{Field: "Value", Err: fmt.Errorf("length must be 20")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)

			var valErrs ValidationErrors
			if !errors.As(err, &valErrs) {
				t.Errorf("expected validation errors, got %T", err)
				return
			}

			var expectedErrs ValidationErrors
			if !errors.As(tt.expectedErr, &expectedErrs) {
				t.Errorf("invalid test setup: expectedErr must be ValidationErrors")
				return
			}

			if len(valErrs) != len(expectedErrs) {
				t.Errorf("expected %d errors, got %d: %v", len(expectedErrs), len(valErrs), valErrs)
				return
			}

			expectedMap := make(map[string]bool)
			for _, e := range expectedErrs {
				key := e.Field + ":" + e.Err.Error()
				expectedMap[key] = true
			}

			for _, actual := range valErrs {
				key := actual.Field + ":" + actual.Err.Error()
				if !expectedMap[key] {
					t.Errorf("unexpected error: %s", key)
				}
			}
		})
	}
}

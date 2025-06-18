package hw09structvalidator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseValidTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected []ValidatorRule
	}{
		{
			name:     "no tag",
			tag:      "",
			expected: nil,
		},
		{
			name:     "empty tag",
			tag:      `validate:""`,
			expected: nil,
		},
		{
			name: "single rule",
			tag:  `validate:"len:5"`,
			expected: []ValidatorRule{
				{Name: "len", Param: "5"},
			},
		},
		{
			name: "multiple rules",
			tag:  `validate:"min:5|max:10"`,
			expected: []ValidatorRule{
				{Name: "min", Param: "5"},
				{Name: "max", Param: "10"},
			},
		},
		{
			name: "in validator with multiple values",
			tag:  `validate:"in:admin,user,guest"`,
			expected: []ValidatorRule{
				{
					Name:   "in",
					Param:  "admin,user,guest",
					Params: []string{"admin", "user", "guest"},
				},
			},
		},
		{
			name: "multiple rules with spaces",
			tag:  `validate:"min: 5 | max: 10 "`,
			expected: []ValidatorRule{
				{Name: "min", Param: "5"},
				{Name: "max", Param: "10"},
			},
		},
		{
			name: "empty rule",
			tag:  `validate:"min:5||max:10"`,
			expected: []ValidatorRule{
				{Name: "min", Param: "5"},
				{Name: "max", Param: "10"},
			},
		},
		{
			name: "rule without parameter",
			tag:  `validate:"min"`,
			expected: []ValidatorRule{
				{Name: "min", Param: ""},
			},
		},
		{
			name: "rule with colon but without parameter",
			tag:  `validate:"min:"`,
			expected: []ValidatorRule{
				{Name: "min", Param: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := reflect.StructTag(tt.tag)
			rules, err := ParseValidTag(tag)
			require.NoError(t, err)
			require.Equal(t, tt.expected, rules)
		})
	}
}

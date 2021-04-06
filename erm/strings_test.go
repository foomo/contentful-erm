package erm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLatinLettersOnly(t *testing.T) {
	for _, tt := range []struct {
		input          string
		expectedOutput string
	}{
		{"Hello World!", "HelloWorld"},
		{"something", "something"},
		{"!.-123", ""},
		{"!.-x123", "x"},
		{"", ""},
		{"ёшкар юла", ""},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.expectedOutput, latinLettersOnly(tt.input))
		})
	}
}

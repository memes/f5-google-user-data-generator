package generators_test

import (
	"errors"
	"net"
	"testing"

	"go.uber.org/goleak"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestShaveMustache(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "empty",
			text:     "",
			expected: "",
		},
		{
			name:     "whitespace",
			text:     "  \t\n",
			expected: "",
		},
		{
			name:     "leading-ws",
			text:     " {{{ test }}}",
			expected: "test",
		},
		{
			name:     "trailing-ws",
			text:     "{{{ test }}} \n",
			expected: "test",
		},
		{
			name:     "no-ws",
			text:     "{{{test}}}",
			expected: "test",
		},
		{
			name:     "extra-ws-internal",
			text:     "{{{    test \t\t\n}}}",
			expected: "test",
		},
		{
			name:     "double-handlebars",
			text:     "{{ test }}",
			expected: "test",
		},
		{
			name:     "single-handlebars",
			text:     "{ test }",
			expected: "test",
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			actual := generators.ShaveMustache(tst.text)
			if tst.expected != actual {
				t.Errorf("Expected %q, got %q", tst.expected, actual)
			}
		})
	}
}

func TestChomp(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "empty",
			text:     "",
			expected: "",
		},
		{
			name:     "simple",
			text:     "test",
			expected: "test",
		},
		{
			name:     "whitespace",
			text:     "  \t\n",
			expected: "  \t",
		},
		{
			name:     "cr-lf",
			text:     "test\r\n",
			expected: "test",
		},
		{
			name:     "multiple",
			text:     "test\n\n",
			expected: "test",
		},
		{
			name:     "leading",
			text:     "\r\ntest",
			expected: "\r\ntest",
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			actual := generators.Chomp(tst.text)
			if tst.expected != actual {
				t.Errorf("Expected %q, got %q", tst.expected, actual)
			}
		})
	}
}

func TestToYAML(t *testing.T) {
	tests := []struct {
		name          string
		obj           any
		expected      string
		expectedError error
	}{
		{
			name:     "nil",
			obj:      nil,
			expected: "",
		},
		{
			name:     "empty struct",
			obj:      struct{}{},
			expected: "{}\n",
		},
		{
			name:     "empty array",
			obj:      []struct{}{},
			expected: "[]\n",
		},
		{
			name:     "empty map",
			obj:      map[string]string{},
			expected: "{}\n",
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			result, err := generators.ToYAML(tst.obj)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("ToYAML returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case tst.expected != result:
				t.Errorf("Expected %q, got %q", tst.expected, result)
			}
		})
	}
}

func TestVipIdentifier(t *testing.T) {
	tests := []struct {
		name          string
		vip           string
		expected      string
		expectedError any
	}{
		{
			name:          "empty",
			expectedError: net.ParseError{},
		},
		{
			name:          "ipv4-invalid",
			vip:           "500.400.300.200",
			expectedError: net.ParseError{},
		},
		{
			name:     "ipv4-simple",
			vip:      "10.0.10.10",
			expected: "vip_10_0_10_10_32",
		},
		{
			name:     "ipv4-32",
			vip:      "10.0.10.10/32",
			expected: "vip_10_0_10_10_32",
		},
		{
			name:     "ipv4-24",
			vip:      "10.0.10.10/24",
			expected: "vip_10_0_10_0_24",
		},
		{
			name:     "ipv6-simple",
			vip:      "2001:cafe::10",
			expected: "vip_2001_cafe_10_128",
		},
		{
			name:     "ipv6-128",
			vip:      "2001:cafe::10/128",
			expected: "vip_2001_cafe_10_128",
		},
		{
			name:     "ipv4-64",
			vip:      "2001:cafe::10/64",
			expected: "vip_2001_cafe_64",
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			result, err := generators.VipIdentifier(tst.vip)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("VipIdentifier returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.As(err, &tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case tst.expected != result:
				t.Errorf("Expected %q, got %q", tst.expected, result)
			}
		})
	}
}

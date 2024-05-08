package runtimeinit_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/runtimeinit"
)

func TestInterfaceBuilder_Index(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      int
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      0,
			expectedError: nil,
		},
		{
			name:          "one",
			index:         1,
			expected:      0,
			expectedError: generators.ErrInvalidInterfaceIndex,
		},
		{
			name:          "two",
			index:         2,
			expected:      2,
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      3,
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      4,
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      5,
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      6,
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      7,
			expectedError: nil,
		},
		{
			name:          "eight",
			index:         8,
			expected:      0,
			expectedError: generators.ErrInvalidInterfaceIndex,
		},
		{
			name:          "negative",
			index:         -1,
			expected:      0,
			expectedError: generators.ErrInvalidInterfaceIndex,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.Index():
				t.Errorf("Expected %q, got %q", tst.expected, iface.Index())
			}
		})
	}
}

func TestInterfaceBuilder_Name(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "external",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "internal",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "internal2",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "internal3",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "internal4",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "internal5",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "internal6",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.Name():
				t.Errorf("Expected %q, got %q", tst.expected, iface.Name())
			}
		})
	}
}

func TestInterfaceBuilder_SelfIPIdentifier(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "external_self_ip",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "internal_self_ip",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "internal2_self_ip",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "internal3_self_ip",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "internal4_self_ip",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "internal5_self_ip",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "internal6_self_ip",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.SelfIPIdentifier():
				t.Errorf("Expected %q, got %q", tst.expected, iface.SelfIPIdentifier())
			}
		})
	}
}

func TestInterfaceBuilder_Address(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "{{{ EXTERNAL_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "{{{ INTERNAL_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "{{{ INTERNAL2_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "{{{ INTERNAL3_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "{{{ INTERNAL4_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "{{{ INTERNAL5_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "{{{ INTERNAL6_ADDRESS }}}",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.Address():
				t.Errorf("Expected %q, got %q", tst.expected, iface.Address())
			}
		})
	}
}

func TestInterfaceBuilder_VLANTag(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "4094",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "4092",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "4091",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "4090",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "4089",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "4088",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "4087",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.VLANTag():
				t.Errorf("Expected %q, got %q", tst.expected, iface.VLANTag())
			}
		})
	}
}

func TestInterfaceBuilder_MTU(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "{{{ EXTERNAL_MTU }}}",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "{{{ INTERNAL_MTU }}}",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "{{{ INTERNAL2_MTU }}}",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "{{{ INTERNAL3_MTU }}}",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "{{{ INTERNAL4_MTU }}}",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "{{{ INTERNAL5_MTU }}}",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "{{{ INTERNAL6_MTU }}}",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.MTU():
				t.Errorf("Expected %q, got %q", tst.expected, iface.MTU())
			}
		})
	}
}

func TestInterfaceBuilder_TMMName(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "1.0",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "1.2",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "1.3",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "1.4",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "1.5",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "1.6",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "1.7",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.TMMName():
				t.Errorf("Expected %q, got %q", tst.expected, iface.TMMName())
			}
		})
	}
}

func TestInterfaceBuilder_AllowService(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      []string
		expectedError error
	}{
		{
			name:  "zero",
			index: 0,
			expected: []string{
				"tcp:80",
				"tcp:443",
				"tcp:4353",
				"udp:1026",
			},
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      nil,
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      nil,
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      nil,
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      nil,
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      nil,
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      nil,
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)

			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && !reflect.DeepEqual(tst.expected, iface.AllowService()):
				t.Errorf("Expected %q, got %q", tst.expected, iface.AllowService())
			}
		})
	}
}

func TestInterfaceBuilder_GatewayAddress(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "{{{ EXTERNAL_GATEWAY_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "{{{ INTERNAL_GATEWAY_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "{{{ INTERNAL2_GATEWAY_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "{{{ INTERNAL3_GATEWAY_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "{{{ INTERNAL4_GATEWAY_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "{{{ INTERNAL5_GATEWAY_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "{{{ INTERNAL6_GATEWAY_ADDRESS }}}",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.GatewayAddress():
				t.Errorf("Expected %q, got %q", tst.expected, iface.GatewayAddress())
			}
		})
	}
}

func TestInterfaceBuilder_GatewayRouteName(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "external_gw_rt",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "internal_gw_rt",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "internal2_gw_rt",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "internal3_gw_rt",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "internal4_gw_rt",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "internal5_gw_rt",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "internal6_gw_rt",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.GatewayRouteName():
				t.Errorf("Expected %q, got %q", tst.expected, iface.GatewayRouteName())
			}
		})
	}
}

func TestInterfaceBuilder_NetworkAddress(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "{{{ EXTERNAL_NETWORK_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "{{{ INTERNAL_NETWORK_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "{{{ INTERNAL2_NETWORK_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "{{{ INTERNAL3_NETWORK_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "{{{ INTERNAL4_NETWORK_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "{{{ INTERNAL5_NETWORK_ADDRESS }}}",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "{{{ INTERNAL6_NETWORK_ADDRESS }}}",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.NetworkAddress():
				t.Errorf("Expected %q, got %q", tst.expected, iface.NetworkAddress())
			}
		})
	}
}

func TestInterfaceBuilder_NetworkBitmask(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "{{{ EXTERNAL_NETWORK_BITMASK }}}",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "{{{ INTERNAL_NETWORK_BITMASK }}}",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "{{{ INTERNAL2_NETWORK_BITMASK }}}",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "{{{ INTERNAL3_NETWORK_BITMASK }}}",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "{{{ INTERNAL4_NETWORK_BITMASK }}}",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "{{{ INTERNAL5_NETWORK_BITMASK }}}",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "{{{ INTERNAL6_NETWORK_BITMASK }}}",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.NetworkBitmask():
				t.Errorf("Expected %q, got %q", tst.expected, iface.NetworkBitmask())
			}
		})
	}
}

func TestInterfaceBuilder_NetworkRouteName(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		expected      string
		expectedError error
	}{
		{
			name:          "zero",
			index:         0,
			expected:      "default",
			expectedError: nil,
		},
		{
			name:          "two",
			index:         2,
			expected:      "internal_nt_rt",
			expectedError: nil,
		},
		{
			name:          "three",
			index:         3,
			expected:      "internal2_nt_rt",
			expectedError: nil,
		},
		{
			name:          "four",
			index:         4,
			expected:      "internal3_nt_rt",
			expectedError: nil,
		},
		{
			name:          "five",
			index:         5,
			expected:      "internal4_nt_rt",
			expectedError: nil,
		},
		{
			name:          "six",
			index:         6,
			expected:      "internal5_nt_rt",
			expectedError: nil,
		},
		{
			name:          "seven",
			index:         7,
			expected:      "internal6_nt_rt",
			expectedError: nil,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			iface, err := runtimeinit.InterfaceBuilder(tst.index)
			switch {
			case tst.expectedError == nil && err != nil:
				t.Errorf("InterfaceBuilder returned an unexpected error: %v", err)
			case tst.expectedError != nil && !errors.Is(err, tst.expectedError):
				t.Errorf("Expected error %v, got error %v", tst.expectedError, err)
			case iface != nil && tst.expected != iface.NetworkRouteName():
				t.Errorf("Expected %q, got %q", tst.expected, iface.NetworkRouteName())
			}
		})
	}
}

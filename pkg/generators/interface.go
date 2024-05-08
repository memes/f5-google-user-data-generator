package generators

import (
	"errors"
	"fmt"
)

// Defines the properties that define a BIG-IP network interface.
type Interface interface {
	Index() int
	Name() string
	SelfIPIdentifier() string
	Address() string
	VLANTag() string
	MTU() string
	TMMName() string
	AllowService() []string
	GatewayAddress() string
	GatewayRouteName() string
	NetworkAddress() string
	NetworkBitmask() string
	NetworkRouteName() string
}

// Defines an InterfaceBuilder function that can create an object that implements
// the Interface interface for the network interface at 0-based ordinal index.
type InterfaceBuilder func(index int) (Interface, error)

var ErrInvalidInterfaceIndex = errors.New("interface index must be an integer between 0 and 7, excluding 1")

type staticInterface struct {
	index int
}

func (si staticInterface) Index() int {
	return si.index
}

func (si staticInterface) Name() string {
	switch si.index {
	case 0:
		return "external"
	case 2:
		return "internal"
	}
	return fmt.Sprintf("internal%d", si.index-1)
}

func (si staticInterface) SelfIPIdentifier() string {
	return fmt.Sprintf("%s_self_ip", si.Name())
}

func (si staticInterface) Address() string {
	return fmt.Sprintf("10.%d.0.10", si.index)
}

func (si staticInterface) VLANTag() string {
	return fmt.Sprintf("%d", 4094-si.index)
}

func (si staticInterface) MTU() string {
	return "1460"
}

func (si staticInterface) TMMName() string {
	return fmt.Sprintf("1.%d", si.index)
}

func (si staticInterface) AllowService() []string {
	if si.index == 0 {
		return []string{
			"tcp:80",
			"tcp:443",
			"tcp:4353",
			"udp:1026",
		}
	}
	return nil
}

func (si staticInterface) GatewayAddress() string {
	return fmt.Sprintf("10.%d.0.1", si.index)
}

func (si staticInterface) GatewayRouteName() string {
	return fmt.Sprintf("%s_gw_rt", si.Name())
}

func (si staticInterface) NetworkAddress() string {
	return fmt.Sprintf("10.%d.0.0", si.index)
}

func (si staticInterface) NetworkBitmask() string {
	return "24"
}

func (si staticInterface) NetworkRouteName() string {
	if si.index == 0 {
		return "default"
	}
	return fmt.Sprintf("%s_nt_rt", si.Name())
}

// Generates an Interface implementation for the data-plane network interface at
// zero-based ordinal index which will have a standard name, a fixed IPv4
// address, etc.
func StaticInterfaceBuilder(index int) (Interface, error) {
	if index < 0 || index == 1 || index > 7 {
		return nil, ErrInvalidInterfaceIndex
	}
	return staticInterface{
		index: index,
	}, nil
}

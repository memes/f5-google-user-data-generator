package runtimeinit

import (
	"fmt"
	"strings"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
)

type runtimeinitInterface struct {
	index int
}

func (ri runtimeinitInterface) Index() int {
	return ri.index
}

func (ri runtimeinitInterface) Name() string {
	switch ri.index {
	case 0:
		return "external"
	case 2:
		return "internal"
	}
	return fmt.Sprintf("internal%d", ri.index-1)
}

func (ri runtimeinitInterface) SelfIPIdentifier() string {
	return fmt.Sprintf("%s_self_ip", ri.Name())
}

func (ri runtimeinitInterface) Address() string {
	return fmt.Sprintf("{{{ %s_ADDRESS }}}", strings.ToUpper(ri.Name()))
}

func (ri runtimeinitInterface) VLANTag() string {
	return fmt.Sprintf("%d", 4094-ri.index)
}

func (ri runtimeinitInterface) MTU() string {
	return fmt.Sprintf("{{{ %s_MTU }}}", strings.ToUpper(ri.Name()))
}

func (ri runtimeinitInterface) TMMName() string {
	return fmt.Sprintf("1.%d", ri.index)
}

func (ri runtimeinitInterface) AllowService() []string {
	if ri.index == 0 {
		return []string{
			"tcp:80",
			"tcp:443",
			"tcp:4353",
			"udp:1026",
		}
	}
	return nil
}

func (ri runtimeinitInterface) GatewayAddress() string {
	return fmt.Sprintf("{{{ %s_GATEWAY_ADDRESS }}}", strings.ToUpper(ri.Name()))
}

func (ri runtimeinitInterface) GatewayRouteName() string {
	return fmt.Sprintf("%s_gw_rt", ri.Name())
}

func (ri runtimeinitInterface) NetworkAddress() string {
	return fmt.Sprintf("{{{ %s_NETWORK_ADDRESS }}}", strings.ToUpper(ri.Name()))
}

func (ri runtimeinitInterface) NetworkBitmask() string {
	return fmt.Sprintf("{{{ %s_NETWORK_BITMASK }}}", strings.ToUpper(ri.Name()))
}

func (ri runtimeinitInterface) NetworkRouteName() string {
	if ri.index == 0 {
		return "default"
	}
	return fmt.Sprintf("%s_nt_rt", ri.Name())
}

// Generates an Interface implementation for the data-plane network interface at
// zero-based ordinal index which will have a standard name, a fixed IPv4
// address, etc.
func InterfaceBuilder(index int) (generators.Interface, error) {
	if index < 0 || index == 1 || index > 7 {
		return nil, generators.ErrInvalidInterfaceIndex
	}
	return runtimeinitInterface{
		index: index,
	}, nil
}

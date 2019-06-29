// +build !linux

package wgcreate

import (
	"golang.zx2c4.com/wireguard/device"
)

func Create(preferredInterfaceName string, mtu uint32, shouldRecreate bool, logger *device.Logger) (resultName string, err error) {
	return createUserspace(preferredInterfaceName, mtu, logger)
}
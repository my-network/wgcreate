// +build !linux

package wgcreate

import (
	"github.com/xaionaro-go/errors"
	"golang.zx2c4.com/wireguard/device"
)

func Create(preferredInterfaceName string, mtu uint32, shouldRecreate bool, logger *device.Logger) (resultName string, err error) {
	defer func() { err = errors.Wrap(err, preferredInterfaceName, mtu, shouldRecreate) }()

	return createUserspace(preferredInterfaceName, mtu, logger)
}
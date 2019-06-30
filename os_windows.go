// +build windows

package wgcreate

import (
	"net"

	"github.com/xaionaro-go/errors"
)

func AddIP(ifaceName string, newIP net.IP, newSubnet net.IPNet) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName, newIP, newSubnet) }()

	return ErrNotSupported
}

func ResetIPs(ifaceName string) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName) }()

	return ErrNotSupported
}

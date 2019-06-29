// +build darwin

package wgcreate

import (
	"net"
	"os/exec"

	"github.com/xaionaro-go/errors"
)

func AddIP(ifaceName string, newIP net.IP, newSubnet net.IPNet) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName, newIP, newSubnet) }()

	subnet := newSubnet
	subnet.IP = newIP

	err = exec.Command("/sbin/ifconfig", ifaceName, subnet.String(), newIP.String()).Run()
	if err != nil {
		return errors.Wrap(err, ifaceName, subnet.String(), newIP.String())
	}


	return
}

func ResetIPs(ifaceName string) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName) }()

	// TODO: implement it

	return
}
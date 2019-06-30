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

	err = exec.Command("/sbin/route", "add", "-net", subnet.String(), "-interface", ifaceName).Run()
	if err != nil {
		return errors.Wrap(err, subnet.String(), ifaceName)
	}

	return
}

func ResetIPs(ifaceName string) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName) }()

	// TODO: implement it

	return
}

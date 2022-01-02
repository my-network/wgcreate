//go:build darwin
// +build darwin

package wgcreate

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"

	"github.com/xaionaro-go/errors"
)

var (
	assignedAddrs   = map[string]map[string]bool{}
	assignedSubnets = map[string]map[string]bool{}
)

func AddIP(ifaceName string, newIP net.IP, newSubnet net.IPNet) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName, newIP, newSubnet) }()

	subnet := newSubnet
	subnet.IP = newIP

	err = exec.Command("/sbin/ifconfig", ifaceName, "alias", newIP.String(), newIP.String()).Run()
	if err != nil {
		return errors.Wrap(err, ifaceName, subnet.String(), newIP.String())
	}

	if assignedAddrs[ifaceName] == nil {
		assignedAddrs[ifaceName] = map[string]bool{}
	}
	assignedAddrs[ifaceName][newIP.String()] = true

	if m := assignedSubnets[ifaceName]; m == nil || !m[subnet.String()] {
		err = exec.Command("/sbin/route", "add", "-net", subnet.String(), "-interface", ifaceName).Run()
		if err != nil {
			return errors.Wrap(err, subnet.String(), ifaceName)
		}
		if assignedSubnets[ifaceName] == nil {
			assignedSubnets[ifaceName] = map[string]bool{}
		}
		assignedSubnets[ifaceName][subnet.String()] = true
	}

	return
}

func ResetIPs(ifaceName string) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName) }()

	for k := range assignedAddrs[ifaceName] {
		err = exec.Command("/sbin/ifconfig", ifaceName, "-alias", k).Run()
		if err != nil {
			return errors.Wrap(err, ifaceName)
		}
		delete(assignedAddrs[ifaceName], k)
	}

	return
}

var (
	correctInterfacePattern, _ = regexp.Compile(`^utun[0-9]+$`)
)

func findFreeUtunName() (string, error) {
	for i := 0; i < 256; i++ {
		ifaceName := fmt.Sprintf("utun%d", i)
		output, err := exec.Command("/sbin/ifconfig", ifaceName).CombinedOutput()
		if strings.Index(string(output), "does not exist") >= 0 {
			return ifaceName, nil
		}
		if err != nil {
			return "", errors.Wrap(err, ifaceName, string(output))
		}
	}
	return "", ErrNoFreeInterface
}

func Create(preferredInterfaceName string, mtu uint32, shouldRecreate bool, logger *_device.Logger) (resultName string, err error) {
	defer func() { err = errors.Wrap(err, preferredInterfaceName, mtu, shouldRecreate) }()

	doDefaultsTriesIncreaseNofile()

	if !correctInterfacePattern.MatchString(preferredInterfaceName) {
		preferredInterfaceName, err = findFreeUtunName()
		if err != nil {
			return
		}
	}

	return createUserspace(preferredInterfaceName, mtu, logger)
}

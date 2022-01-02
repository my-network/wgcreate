//go:build linux
// +build linux

package wgcreate

import (
	"net"
	"strconv"
	"syscall"

	"github.com/lorenzosaino/go-sysctl"
	"golang.zx2c4.com/wireguard/device"

	"github.com/xaionaro-go/errors"
	"github.com/xaionaro-go/netlink"
)

func findLink(ifaceName string) (link netlink.Link, err error) {
	defer func() { err = errors.Wrap(err, ifaceName) }()

	links, err := netlink.LinkList()
	if err != nil {
		return
	}
	for _, link := range links {
		if link.Attrs().Name == ifaceName {
			return link, nil
		}
	}

	return nil, ErrInterfaceNotFound
}

func sysctlGetValue(key string) (intValue int64, err error) {
	value, err := sysctl.Get(key)
	if err != nil {
		return
	}

	intValue, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		return
	}

	return
}

func sysctlIncreaseTo(key string, value int64, logger *device.Logger) {
	oldValue, err := sysctlGetValue(key)
	if err != nil {
		logger.Verbosef(`unable to get current sysctl value by key "%v": %v`, key, err)
		return
	}

	if value <= oldValue {
		return
	}

	err = sysctl.Set(key, strconv.FormatInt(value, 10))
	if err != nil {
		logger.Verbosef(`unable to set sysctl value by key "%v" to "%v": %v`, key, value, err)
		return
	}

	return
}

func Create(preferredInterfaceName string, mtu uint32, shouldRecreate bool, logger *device.Logger) (resultName string, err error) {
	defer func() { err = errors.Wrap(err, preferredInterfaceName, mtu, shouldRecreate) }()

	doDefaultsTriesIncreaseNofile()

	sysctlRequiredValues := map[string]int64{
		"net.ipv4.igmp_max_memberships": 256,
		"fs.inotify.max_user_instances": 256,
		"fs.inotify.max_user_watches":   65536,
	}

	for k, v := range sysctlRequiredValues {
		sysctlIncreaseTo(k, v, logger)
	}

	if shouldRecreate {
		link, err := findLink(preferredInterfaceName)

		if err != nil && !err.(*errors.Error).Has(ErrInterfaceNotFound) {
			return "", err
		}

		if link != nil {
			err = netlink.LinkDel(link)
			if err != nil {
				return "", err
			}
		}
	}

	err = netlink.LinkAdd(&netlink.Wireguard{
		LinkAttrs: netlink.LinkAttrs{
			MTU:    int(mtu),
			Name:   preferredInterfaceName,
			TxQLen: 1000,

			Flags:     net.FlagUp | net.FlagMulticast | net.FlagBroadcast,
			OperState: netlink.OperUp,
		},
	})

	if err == syscall.ENOTSUP {
		logger.Verbosef(`There is no in-kernel support of wireguard on this system. It could negatively affect performance. To avoid it install kernel module "wireguard".`)

		// Fallback to userspace implementation
		resultName, _, err = createUserspace(preferredInterfaceName, mtu, logger)
		if err != nil {
			return
		}

		var link netlink.Link
		link, err = netlink.LinkByName(resultName)
		if err != nil {
			return
		}

		err = netlink.LinkSetUp(link)
		if err != nil {
			return
		}
	}

	if err != nil {
		return
	}

	resultName = preferredInterfaceName
	return
}

func AddIP(ifaceName string, newIP net.IP, newSubnet net.IPNet) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName, newIP, newSubnet) }()

	link, err := findLink(ifaceName)
	if err != nil {
		return
	}

	subnet := newSubnet
	subnet.IP = newIP
	err = netlink.AddrAdd(link, &netlink.Addr{
		IPNet: &subnet,
	})
	if err != nil {
		return
	}

	return
}

func ResetIPs(ifaceName string) (err error) {
	defer func() { err = errors.Wrap(err, ifaceName) }()

	link, err := findLink(ifaceName)
	if err != nil {
		return
	}

	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return
	}

	for _, addr := range addrs {
		err = netlink.AddrDel(link, &addr)
		if err != nil {
			return
		}
	}

	return
}

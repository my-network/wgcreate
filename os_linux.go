// +build linux

package wgcreate

import (
	e "errors"
	"net"

	"github.com/xaionaro-go/errors"
	"github.com/xaionaro-go/netlink"
)

var (
	ErrInterfaceNotFound = e.New(`interface not found`)
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

func Create(preferredInterfaceName string, mtu uint32, shouldRecreate bool) (resultName string, err error) {
	defer func() { err = errors.Wrap(err, preferredInterfaceName, mtu, shouldRecreate) }()

	if shouldRecreate {
		link, err := findLink(preferredInterfaceName)

		if err != nil && err.(errors.SmartError).OriginalError() != ErrInterfaceNotFound {
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
			MTU:          int(mtu),
			Name:         preferredInterfaceName,

			Flags:        net.FlagUp | net.FlagMulticast | net.FlagBroadcast,
			OperState:    netlink.OperUp,
		},
	})
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
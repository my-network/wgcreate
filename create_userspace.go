package wgcreate

import (
	"syscall"

	"github.com/hashicorp/go-multierror"
	"github.com/xaionaro-go/errors"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
)

func doDefaultsTriesIncreaseNofile() {
	for _, nofileLimit := range []uint64{4096, 12000, 65536, 524288} {
		tryIncreaseNofileTo(nofileLimit)
	}
}

func tryIncreaseNofileTo(newLimit uint64) {
	nofileLimit := &syscall.Rlimit{}
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, nofileLimit)
	if err != nil {
		return
	}
	if nofileLimit.Cur < newLimit {
		nofileLimit.Cur = newLimit
		if nofileLimit.Max < newLimit {
			nofileLimit.Max = newLimit
		}
		_ = syscall.Setrlimit(syscall.RLIMIT_NOFILE, nofileLimit)
	}
}

func createUserspace(ifaceName string, mtu uint32, logger *device.Logger) (resultIfaceName string, closeFunc func() error, err error) {
	defer func() { err = errors.Wrap(err, ifaceName, mtu) }()

	tunDev, err := tun.CreateTUN(ifaceName, int(mtu))
	if err != nil {
		return
	}

	resultIfaceName, err = tunDev.Name()
	if err != nil {
		return
	}

	wgDev := device.NewDevice(tunDev, conn.NewDefaultBind(), logger)

	logger.Verbosef("userspace device started")

	uapiFile, err := ipc.UAPIOpen(resultIfaceName)
	if err != nil {
		wgDev.Close()
		return
	}

	uapi, err := ipc.UAPIListen(resultIfaceName, uapiFile)
	if err != nil {
		wgDev.Close()

		return
	}

	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				logger.Verbosef("unable to accept UAPI connection", err)
				return
			}
			go wgDev.IpcHandle(conn)
		}
	}()

	logger.Verbosef("UAPI started")
	closeFunc = func() error {
		var err error
		if addErr := uapiFile.Close(); addErr != nil {
			err = multierror.Append(err, addErr)
		}
		if addErr := uapi.Close(); addErr != nil {
			err = multierror.Append(err, addErr)
		}
		wgDev.Close()
		return err
	}

	return
}

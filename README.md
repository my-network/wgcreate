There is package [https://github.com/WireGuard/wgctrl-go](https://github.com/WireGuard/wgctrl-go) which does everything to manage WireGuard interfaces except creating them. And authors of the package says:

> This package implements WireGuard configuration protocol operations, enabling the configuration of existing WireGuard devices. Operations such as creating WireGuard devices, or applying IP addresses to those devices, are out of scope for this package.

So `wgcreate` just creates WireGuard interfaces to be managed by [https://github.com/WireGuard/wgctrl-go](https://github.com/WireGuard/wgctrl-go).

It uses kernel-space support if it available. If not then it uses userspace implementation of wireguard.

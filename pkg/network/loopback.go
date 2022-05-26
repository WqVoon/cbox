package network

import (
	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/log"
)

// 启动 loopback 设备
func SetupLookback() {
	dev, err := netlink.LinkByName("lo")
	if err != nil {
		log.Errorln("failed to get loopback, err:", err)
	}

	if err := netlink.LinkSetUp(dev); err != nil {
		log.Errorln("failed to set up loopback, err:", err)
	}
}

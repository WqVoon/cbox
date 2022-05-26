package network

import (
	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/log"
)

// 清空一个设备的所有地址
func CleanAddr(dev netlink.Link) {
	addrList, err := netlink.AddrList(dev, netlink.FAMILY_V4)
	if err != nil {
		log.Errorln("failed to get addr list, err:", err)
	}

	for _, addr := range addrList {
		if err := netlink.AddrDel(dev, &addr); err != nil {
			log.Errorln("failed to delete addr, err:", err)
		}
	}
}

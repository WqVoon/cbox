package network

import (
	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/config"
	"github.com/wqvoon/cbox/pkg/log"
)

var bridge netlink.Link

// 初始化网桥设备，如果不存在则创建，否则更新其地址并将其启动
func InitBridge() {
	bridgeConfig := config.GetNetworkConfig()

	var err error
	bridge, err = netlink.LinkByName(bridgeConfig.Name)
	if err != nil { // 这里假定出错一定是因为不存在对应的设备
		linkAttrs := netlink.NewLinkAttrs()
		linkAttrs.Name = bridgeConfig.Name
		bridge = &netlink.Bridge{LinkAttrs: linkAttrs}

		if err := netlink.LinkAdd(bridge); err != nil {
			log.Errorln("failed to add bridge, err:", err)
		}
	}

	bridgeAddr, err := netlink.ParseAddr(bridgeConfig.Addr)
	if err != nil {
		log.Errorln("failed to parse bridge addr, err:", err)
	}

	CleanAddr(bridge)
	if err := netlink.AddrAdd(bridge, bridgeAddr); err != nil {
		log.Errorln("failed to set addr for bridge, err:", err)
	}

	if err := netlink.LinkSetUp(bridge); err != nil {
		log.Errorln("failed to set up bridge, err:", err)
	}
}

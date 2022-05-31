package network

import (
	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/config"
	"github.com/wqvoon/cbox/pkg/log"
)

var Bridge *Device

// 获取或创建一个网桥
func GetOrCreateBridge(name string) *Device {
	rawBridge, err := netlink.LinkByName(name)
	if err != nil { // 这里假定出错一定是因为不存在对应的设备
		linkAttrs := netlink.NewLinkAttrs()
		linkAttrs.Name = name
		rawBridge = &netlink.Bridge{LinkAttrs: linkAttrs}

		if err := netlink.LinkAdd(rawBridge); err != nil {
			log.Errorln("failed to create bridge, err:", err)
		}
	}
	return &Device{rawDevice: rawBridge}
}

// 初始化网桥设备，如果不存在则创建，否则更新其地址并将其启动
func InitBridge() {
	cfg := config.GetNetworkConfig()

	tmpBridge := GetOrCreateBridge(cfg.Name)
	{
		tmpBridge.SetAddress(cfg.BridgeCIDR)
		tmpBridge.SetUp()
	}

	Bridge = tmpBridge
}

package network

import (
	"net"

	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/log"
)

// 设置默认路由
func SetDefaultRoute(via net.IP, dev *Device) {
	route := netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: dev.rawDevice.Attrs().Index, // 设置走哪个设备，所以 src 字段就不设置了
		Gw:        via,                         // 网关地址
		Dst:       nil,                         // 置空 dst 字段，所以是 default 路由
	}

	if err := netlink.RouteAdd(&route); err != nil {
		log.Errorln("failed to add default route, err:", err)
	}
}

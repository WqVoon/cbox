package network

import (
	"net"

	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/config"
	"github.com/wqvoon/cbox/pkg/log"
)

// 创建 containerID 对应的容器要使用的 veth pair 设备
func CreateVethPairFor(containerID string) {
	veth0 := "veth0_" + containerID[:6]
	veth1 := "veth1_" + containerID[:6]
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = veth0
	veth0Struct := &netlink.Veth{
		LinkAttrs: linkAttrs,
		PeerName:  veth1,
	}
	if err := netlink.LinkAdd(veth0Struct); err != nil {
		log.Errorf("failed to add %s, err: %v\n", veth0, err)
	}

	if err := netlink.LinkSetUp(veth0Struct); err != nil {
		log.Errorln("failed to setup veth, err:", err)
	}

	if err := netlink.LinkSetMaster(veth0Struct, bridge); err != nil {
		log.Errorln("failed to plug veth into bridge, err:", err)
	}
}

// 设置 veth1 设备的网络命名空间
func SetVeth1NS(containerID string, nsFd int) {
	vethName := "veth1_" + containerID[:6]
	vethDev, err := netlink.LinkByName(vethName)
	if err != nil {
		log.Errorln("failed to get veth1, err:", err)
	}

	if err := netlink.LinkSetNsFd(vethDev, nsFd); err != nil {
		log.Errorln("failed to set ns for veth1, err:", err)
	}
}

func ConfigAndUpVeth1(containerID string, ip string) {
	vethName := "veth1_" + containerID[:6]
	vethDev, err := netlink.LinkByName(vethName)
	if err != nil {
		log.Errorln("failed to get veth1, err:", err)
	}

	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		log.Errorln("failed to parse ip address, err:", err)
	}

	if err := netlink.AddrAdd(vethDev, addr); err != nil {
		log.Errorf("Error assigning IP to veth1: %v\n", err)
	}

	if err := netlink.LinkSetUp(vethDev); err != nil {
		log.Errorln("failed to setup veth, err:", err)
	}

	// TODO: 这部分逻辑应该集成到 Bridge 中
	bridgeIp, _, err := net.ParseCIDR(config.GetNetworkConfig().BridgeCIDR)
	if err != nil {
		log.Errorln("failed to parse bridge cidr, err:", err)
	}

	route := netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: vethDev.Attrs().Index,
		Gw:        bridgeIp,
		Dst:       nil,
	}
	if err := netlink.RouteAdd(&route); err != nil {
		log.Errorln("failed to add route, err:", err)
	}
}

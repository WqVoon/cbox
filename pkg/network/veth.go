package network

import (
	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/log"
)

type VethPair struct {
	HostPeer *Device // 宿主机的一端
	CntrPeer *Device // 容器的一端
}

// 创建 containerID 对应的容器要使用的 veth pair 设备
func CreateVethPairFor(containerID string) *VethPair {
	hostName := "host_" + containerID[:6]
	cntrName := "cntr_" + containerID[:6]

	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = hostName

	hostVeth := &netlink.Veth{
		LinkAttrs: linkAttrs,
		PeerName:  cntrName,
	}

	if err := netlink.LinkAdd(hostVeth); err != nil {
		log.Errorln("failed to create veth, err:", err)
	}

	cntrVeth, err := netlink.LinkByName(cntrName)
	if err != nil {
		log.Errorln("failed to get veth, err:", err)
	}

	return &VethPair{
		HostPeer: &Device{rawDevice: hostVeth},
		CntrPeer: &Device{rawDevice: cntrVeth},
	}
}

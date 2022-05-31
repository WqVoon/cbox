package network

import (
	"net"

	"github.com/vishvananda/netlink"
	"github.com/wqvoon/cbox/pkg/log"
)

type Device struct {
	rawDevice netlink.Link
	address   *netlink.Addr // 假定一个设备仅有一个地址
}

// 为设备设置地址
func (d *Device) SetAddress(cidr string) {
	addr, err := netlink.ParseAddr(cidr)
	if err != nil {
		log.Errorln("failed to parse addr, err:", err)
	}

	CleanAddr(d.rawDevice)
	if err := netlink.AddrAdd(d.rawDevice, addr); err != nil {
		log.Errorln("failed to set addr for device, err:", err)
	}

	d.address = addr
}

// 返回设备的 ip
func (d *Device) GetIP() net.IP {
	return d.address.IP
}

// 启动设备
func (d *Device) SetUp() {
	if err := netlink.LinkSetUp(d.rawDevice); err != nil {
		log.Errorln("failed to set up device, err:", err)
	}
}

// 设置主设备，这里用于连接网桥
func (d *Device) SetMaster(master *Device) {
	if err := netlink.LinkSetMaster(d.rawDevice, master.rawDevice); err != nil {
		log.Errorln("failed to set master for device, err:", err)
	}
}

// 根据 fd 设置设备的网络命名空间
func (d *Device) SetNamespace(nsFd int) {
	if err := netlink.LinkSetNsFd(d.rawDevice, nsFd); err != nil {
		log.Errorln("failed to set namespace, err:", err)
	}
}

package address

import (
	"crypto/rand"
	"net"
)

// 创建一个 Mac 地址，使用随机数来近似唯一
func CreateMACAddress() net.HardwareAddr {
	hw := make(net.HardwareAddr, 6)
	hw[0] = 0x02 // 本地生成的单播地址
	hw[1] = 0xcb // 随便弄的，做一个区分
	rand.Read(hw[2:])
	return hw
}

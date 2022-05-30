package address

import (
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// 键是 IPv4Address.RawVal
type usedIPType map[uint32]struct{}

var (
	// 当前已经被使用的 ip 地址，由 rootdir/used_ip 文件解析而来
	usedIP usedIPType
	// rootdir/used_ip 文件的实际地址
	filePath string
)

func InitUsedIP() {
	filePath = rootdir.GetUsedIPPath()
	utils.GetObjFromJsonFile(filePath, &usedIP)
	if usedIP == nil {
		usedIP = make(usedIPType)
	}
}

// 根据 uint32 类型的 ip 形式来查询是否已被使用
func AddressIsUsed(ip uint32) bool {
	_, used := usedIP[ip]
	return used
}

// 占用一个 ip 地址，即将这个地址写入 used_ip
func RequireIPByUint32(ip uint32) {
	if _, used := usedIP[ip]; used {
		return
	}

	usedIP[ip] = struct{}{}
	utils.SaveObjToJsonFile(filePath, usedIP)
}

// 释放一个 ip 地址，即将这个地址从 used_ip 中移除
func ReleaseIPByString(ip string) {
	if !IsValidIPv4(ip) {
		return
	}

	addr := ParseIPv4FromString(ip)

	if _, used := usedIP[addr.RawVal]; used {
		delete(usedIP, addr.RawVal)
		utils.SaveObjToJsonFile(filePath, usedIP)
	}
}

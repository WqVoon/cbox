package address

import (
	"fmt"
	"net"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/utils"
)

// IPv4 的地址，由 ParseIPv4 函数得到
type IPv4Address struct {
	RawIP  net.IP // 保证长度为 4
	Mask   string // 形如 `20` 这样的字符串
	RawVal uint32 // RawIP 的 byte 对应的数字，比如 1.2.3.4 = 1*2^24 + 2*2^16 + 3*2^8 + 4
}

const (
	byte0Base uint32 = 24
	byte1Base uint32 = 16
	byte2Base uint32 = 8
	byte3Base uint32 = 0
)

// 解析字符串获取 IPv4Address 结构
func ParseIPv4FromString(s string) *IPv4Address {
	splitedS := strings.Split(s, "/")
	if len(splitedS) != 2 {
		log.Errorln("failed to parse", s, "to ip v4 address")
	}
	ipStr, maskStr := splitedS[0], splitedS[1]

	rawIP := net.ParseIP(ipStr)
	if rawIP == nil {
		log.Errorln("failed to parse", s, "to ip v4 address")
	}

	rawIP = rawIP.To4()
	return &IPv4Address{
		RawIP: rawIP,
		Mask:  maskStr,
		RawVal: (uint32(rawIP[0]) << byte0Base) +
			(uint32(rawIP[1]) << byte1Base) +
			(uint32(rawIP[2]) << byte2Base) +
			(uint32(rawIP[3]) << byte3Base),
	}
}

// 解析 uint32 获取 IPv4Address 结构，maskStr 为 `20` 这样的数字内容的字符串
func ParseIPv4FromUint32(u uint32, maskStr string) *IPv4Address {
	rawIP := net.IPv4(
		byte(u>>byte0Base),
		byte(u<<byte2Base>>byte0Base),
		byte(u<<byte1Base>>byte0Base),
		byte(u<<byte0Base>>byte0Base),
	).To4()

	return &IPv4Address{
		RawIP:  rawIP,
		Mask:   maskStr,
		RawVal: u,
	}
}

func (addr *IPv4Address) String() string {
	rawIp := addr.RawIP
	return fmt.Sprintf("%d.%d.%d.%d/%s",
		rawIp[0], rawIp[1], rawIp[2], rawIp[3], addr.Mask)
}

// 获取一个空闲的 IP 地址，start 和 end 都是 `x.x.x.x/x` 这样的格式
func GetIPAddress(start, end string) string {
	startAddr := ParseIPv4FromString(start)
	endAddr := ParseIPv4FromString(end)

	if startAddr.RawVal >= endAddr.RawVal {
		log.Errorln("error value of ip range, start should less than end")
	}

	usedIp := map[uint32]struct{}{}
	utils.GetObjFromJsonFile(rootdir.GetUsedIPPath(), &usedIp)

	for ip := startAddr.RawVal; ip <= endAddr.RawVal; ip++ {
		if !AddressIsUsed(ip) {
			RequireIPByUint32(ip)
			return ParseIPv4FromUint32(ip, startAddr.Mask).String()
		}
	}

	log.Errorln("no avaliable ip address")
	return ""
}

// 判断字符串 s 的内容是否是一个有效的 ip v4 地址，这个地址会携带 mask
func IsValidIPv4(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

package network

import "github.com/wqvoon/cbox/pkg/network/dns"

func Init() {
	dns.Init()
	InitBridge()
}

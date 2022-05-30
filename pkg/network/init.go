package network

import (
	"github.com/wqvoon/cbox/pkg/network/address"
	"github.com/wqvoon/cbox/pkg/network/dns"
)

func Init() {
	dns.Init()
	address.InitUsedIP()
	InitBridge()
}

# 开启本机的转发功能
sysctl net.ipv4.conf.all.forwarding=1

# 配置 iptables 规则，让属于网桥网络的数据包在发往外网时经过 DNAT
# 网桥名与网络号由 rootdir 中 config.json 进行配置
iptables -t nat -I POSTROUTING 1 -s 172.29.0.0/16 ! -o cbox0 -j MASQUERADE

# 除此之外，诸如 virtualbox 的环境可能会在 filter 表的 FORWARD 链中加入一些规则导致非本地地址不进行转发
# 这可能会导致虽然可以 ping 外部的 ip 但不能使用域名来 ping 的玄学问题，此时可以使用类似下面的规则来解决
iptables -t filter -I FORWARD 1 -i cbox0 -j RETURN
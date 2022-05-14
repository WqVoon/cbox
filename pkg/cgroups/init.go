package cgroups

import "github.com/wqvoon/cbox/pkg/config"

var (
	Cpu *CPUCGroup
	Mem *MemCGroup
	Pid *PIDCGroup
)

func Init() {
	cc := config.GetCgroupConfig()
	if !cc.Enable { // 如果未启用那么直接结束这个初始化，此时上面 var 中的三个变量均不可使用
		return
	}

	Cpu = GetCPUCGroupByPath(cc.CPUCgroupPath).GetOrCreateSubCGroup(cc.Name)
	{
		Cpu.SetNotifyOnRelease(true)
		Cpu.SetCPULimit(cc.CPULimit)
	}

	Mem = GetMemCGroupByPath(cc.MemCgroupPath).GetOrCreateSubCGroup(cc.Name)
	{
		Mem.SetNotifyOnRelease(true)
		Mem.SetMemLimit(cc.MemLimit)
	}

	Pid = GetPIDCGroupByPath(cc.PIDCgroupPath).GetOrCreateSubCGroup(cc.Name)
	{
		Pid.SetNotifyOnRelease(true)
		Pid.SetProcessLimit(cc.ProcessLimit)
	}
}

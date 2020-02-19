package main

import (
	"os"
	//"errors"
	pb "github.com/kata-containers/agent/protocols/grpc"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//grpcStatus "google.golang.org/grpc/status"
)

func (s *sandbox) GetPodStats() (*pb.PodStats, error) {
	//Buddy, err := getPodBuddyStats()
	//if err != nil {
	//	return nil, err
	//}

	//CpuStats, err := getPodCpuStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//IoStats, err := getPodIOStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//IrqStats, err := getPodIrqStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//LoadStats, err := getPodLoadStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//MemStats, err := getPodMemStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//MemxStats, err := getPodMemxStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//PartStats, err := getPodPartitionStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//PcswStats, err := getPodPcswStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//PercpuStats, err := getPodPercpuStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//PercpuSchedStats, err := getPodPercpuSchedStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//SoftirqStats, err := getPodSoftirqStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//SwapStats, err := getPodSwapStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//TcpStats, err := getPodTcpStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//TcpxStats, err := getPodTcpxStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//TrafficStats, err := getPodTrafficStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//UdpStats, err := getPodUdpStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//VmstatStats, err := getPodVmstatStats()
	//if err != nil {
	//	return nil, err
	//}
	//
	//GpuInfo, err := getPodGpuInfo()
	//if err != nil {
	//	return nil, err
	//}

	return &pb.PodStats {
		//Buddy: Buddy,
		//CpuStats,
		//IoStats,
		//IrqStats,
		//LoadStats,
		//MemStats,
		//MemxStats,
		//PartStats,
		//PcswStats,
		//PercpuStats,
		//PercpuSchedStats,
		//SoftirqStats,
		//SwapStats,
		//TcpStats,
		//TcpxStats,
		//TrafficStats,
		//UdpStats,
		//VmstatStats,
		//GpuInfo,
	}, nil
}

func (c *container) GetContainerStats() (*pb.ContainerStats, error) {
	base, err := c.GetContainerBaseStats()
	if err != nil {
		return nil, err
	}

	var pstats []*pb.ProcessStats
	for _, p := range c.processes {
		pstat, err := p.GetProcessStats()
		if err != nil {
			continue
		}

		pstats = append(pstats, pstat)
	}

	return &pb.ContainerStats {
		BaseStats: base,
		ProcStats: pstats,
	}, nil
}

// ajdk_score come from environment of init process
func (c *container) GetContainerBaseStats() (*pb.ContainerBaseStats, error) {
	//ContainerId := c.container.ID()
	//
	//// get cgroup paths and then read out info from cgroup fs
	//state, err := c.container.State()
	//cgroups := state.CgroupPaths

	//CgroupCpu, err := getContainerCgroupCpu(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// get ajdk_score
	//CgroupCpu.JdkMaxProcessLimit = c.getAjdkScore()
	//
	//CgroupCpushare, err := getContainerCgroupCpushare(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupIo, err := getContainerCgroupIo(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupIotime, err := getContainerCgroupIoTime(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupIox, err := getContainerCgroupIox(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupLoad, err := getContainerCgroupLoad(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupMem, err := getContainerCgroupMem(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupMemdelay, err := getContainerCgroupMemDelay(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupLatencyStat, err := getContainerCgroupLatencyStats(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupMemx, err := getContainerCgroupMemx(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupPcsw, err := getContainerCgroupPcsw(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupPercpu, err := getContainerCgroupPercpu(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupSchedbvt, err := getContainerCgroupSchedbvt(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupSchedcg, err := getContainerCgroupSchedcg(cgroups)
	//if err != nil {
	//	return nil,err
	//}
	//
	//CgroupTask, err := getContainerCgroupTask(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupTcpStats, err := getContainerCgroupTcpStats(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupTcpxStats, err := getContainerCgroupTcpxStats(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupTrafficStats, err := getContainerCgroupTrafficStats(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupUdpStats, err := getContainerCgroupUdpStats(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//CgroupVm, err := getContianerCgroupVm(cgroups)
	//if err != nil {
	//	return nil, err
	//}
	//
	//PartitionStats, err := getContainerPartitionStats(cgroups)
	//if err != nil {
	//	return nil, err
	//}

	return &pb.ContainerBaseStats {
		//ContainerId,
		//CgroupCpu,
		//CgroupCpushare,
		//CgroupIo,
		//CgroupIotime,
		//CgroupIox,
		//CgroupLoad,
		//CgroupMem,
		//CgroupMemdelay,
		//CgroupLatencyStat,
		//CgroupMemx,
		//CgroupPcsw,
		//CgroupPercpu,
		//CgroupSchedbvt,
		//CgroupSchedcg,
		//CgroupTask,
		//CgroupTcpStats,
		//CgroupTcpxStats,
		//CgroupTafficStats,
		//CgroupUdpStats,
		//CgroupVm,
		//PartitionStats,
	}, nil
}

func (p *process) GetProcessStats() (*pb.ProcessStats, error) {
	Pid, err := p.process.Pid()
	if err != nil {
		return nil, err
	}

	PidStats, err := getProcessPidStats(Pid)
	if err != nil {
		return nil, err
	}

	CgroupSched, err := getProcessCgroupSched(Pid)
	if err != nil {
		return nil, err
	}

	ProcCpuStats, err := getProcessProcCpuStats(Pid)
	if err != nil {
		return nil, err
	}

	ProcIoStats, err := getProcessProcIoStats(Pid)
	if err != nil {
		return nil, err
	}

	ProcMemStats, err := getProcessProcMemStats(Pid)
	if err != nil {
		return nil, err
	}

	return &pb.ProcessStats {
		Pid: int32(Pid),
		PidStats: PidStats,
		CgroupSched: CgroupSched,
		ProcCpuStats: ProcCpuStats,
		ProcIoStats: ProcIoStats,
		ProcMemStats: ProcMemStats,
	}, nil
}

const (
	aliProcPath = "/sys/fs/cgroup/cpuacct/cpuacct.proc_stat"
)

// true if ali kernel, false on other kernels
func aliKernel() (bool) {
	_, err := os.Stat(aliProcPath)
	return err == nil || os.IsExist(err)
}

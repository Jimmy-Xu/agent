package main

import (
	"errors"
	pb "github.com/kata-containers/agent/protocols/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func (s *sandbox) GetPodStats() (*pb.PodStats, error) {
	Buddy, err := getPodBuddyStats()
	if err != nil {
		return nil, err
	}

	CpuStats, err := getPodCpuStats()
	if err != nil {
		return nil, err
	}

	IoStats, err := getPodIOStats()
	if err != nil {
		return nil, err
	}

	IrqStats, err := getPodIrqStats()
	if err != nil {
		return nil, err
	}

	LoadStats, err := getPodLoadStats()
	if err != nil {
		return nil, err
	}

	MemStats, err := getPodMemStats()
	if err != nil {
		return nil, err
	}

	MemxStats, err := getPodMemxStats()
	if err != nil {
		return nil, err
	}

	PartStats, err := getPodPartitionStats()
	if err != nil {
		return nil, err
	}

	PcswStats, err := getPodPcswStats()
	if err != nil {
		return nil err
	}

	PercpuStats, err := getPodPercpuStats()
	if err != nil {
		return nil, err
	}

	PercpuSchedStats, err := getPodPercpuSchedStats()
	if err != nil {
		return nil, err
	}

	SoftirqStats, err := getPodSoftirqStats()
	if err != nil {
		return nil, err
	}

	SwapStats, err := getPodSwapStats()
	if err != nil {
		return nil, err
	}

	TcpStats, err := getPodTcpStats()
	if err != nil {
		return nil, err
	}

	TcpxStats, err := getPodTcpxStats()
	if err != nil {
		return nil, err
	}

	TrafficStats, err := getPodTrafficStats()
	if err != nil {
		return nil, err
	}

	UdpStats, err := getPodUdpStats()
	if err != nil {
		return nil, err
	}

	VmstatStats, err := getPodVmstatStats()
	if err != nil {
		return nil, err
	}

	GpuInfo, err := getPodGpuInfo()
	if err != nil {
		return nil, err
	}

	return &PodStats {
		Buddy,
		CpuStats,
		IoStats,
		IrqStats,
		LoadStats,
		MemStats,
		MemxStats,
		PartStats,
		PcswStats,
		PercpuStats,
		PercpuSchedStats,
		SoftirqStats,
		SwapStats,
		TcpStats,
		TcpxStats,
		TrafficStats,
		UdpStats,
		VmstatStats,
		GpuInfo,
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

		append(pstats, pstat)
	}

	return &pb.ContainerStats {
		BaseStats: base,
		ProcStats: pstats,
	}, nil
}

func (c *container) GetContainerBaseStats() (*pb.ContainerBaseStats, error) {
	ContainerId := c.container.ID()

	// get cgroup paths and then read out info from cgroup fs
	state, err := c.container.State()
	cgroups := state.CgroupPaths
}

func (p *process) GetProcessStats() (*pb.ProcessStats, error) {
}

package main

import (
	pb "github.com/kata-containers/agent/protocols/grpc"
)

func getProcessPidStats(pid int) (*pb.ProcessPidStats, error) {
	var (
		userCpu    uint64
		sysCpu     uint64
		mem        uint64
		readBytes  uint64
		writeBytes uint64
		fdCnt      uint64
		minFlt     uint64
		majFlt     uint64
		thread     uint64
	)
	return &pb.ProcessPidStats{
		UserCpu:    userCpu,
		SysCpu:     sysCpu,
		Mem:        mem,
		ReadBytes:  readBytes,
		WriteBytes: writeBytes,
		Fdcnt:      fdCnt,
		Minflt:     minFlt,
		Majflt:     majFlt,
		Thread:     thread,
	}, nil
}

func getProcessCgroupSched(pid int) (*pb.ProcessCgroupSched, error) {
	var (
		tasksDelayMs   []uint64
		bvtDelay       uint64
		noiseKickDelay uint64
		lossTime       uint64
	)
	return &pb.ProcessCgroupSched{
		TasksDelayMs:   tasksDelayMs,
		BvtDelay:       bvtDelay,
		NoiseKickDelay: noiseKickDelay,
		LossTime:       lossTime,
	}, nil
}

func getProcessProcCpuStats(pid int) (*pb.ProcessProcCpuStats, error) {
	var (
		userCpu    uint64
		sysCpu     uint64
		blkioDelay uint64
		runOnCpu   uint64
		waitOnRq   uint64
		nThreads   uint64
	)
	return &pb.ProcessProcCpuStats{
		UserCpu:    userCpu,
		SysCpu:     sysCpu,
		BlkioDelay: blkioDelay,
		RunOnCpu:   runOnCpu,
		WaitOnRq:   waitOnRq,
		NThreads:   nThreads,
	}, nil
}

func getProcessProcIoStats(pid int) (*pb.ProcessProcIOStats, error) {
	var (
		rchar  uint64
		wchar  uint64
		rbytes uint64
		wbytes uint64
		iolat  uint64
		syscr  uint64
		syscw  uint64
	)
	return &pb.ProcessProcIOStats{
		Rchar:  rchar,
		Wchar:  wchar,
		Rbytes: rbytes,
		Wbytes: wbytes,
		Iolat:  iolat,
		Syscr:  syscr,
		Syscw:  syscw,
	}, nil
}

func getProcessProcMemStats(pid int) (*pb.ProcessProcMemStats, error) {
	var (
		usage    uint64
		cache    uint64
		anon     uint64
		shmem    uint64
		lock     uint64
		swp      uint64
		totalSwp uint64
	)
	return &pb.ProcessProcMemStats{
		Usage:    usage,
		Cache:    cache,
		Anon:     anon,
		Shmem:    shmem,
		Lock:     lock,
		Swp:      swp,
		TotalSwp: totalSwp,
	}, nil
}

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	pb "github.com/kata-containers/agent/protocols/grpc"
)

const (
	VMSTAT       = "/proc/vmstat"
	STAT         = "/proc/stat"
	SCHEDSTAT    = "/proc/schedstat"
	MEMINFO      = "/proc/meminfo"
	LOADAVG      = "/proc/loadavg"
	DISKSTATS    = "/proc/diskstats"
	NET_DEV      = "/proc/net/dev"
	NET_SNMP     = "/proc/net/snmp"
	APACHERT     = "/tmp/apachert.mmap"
	TCP          = "/proc/net/tcp"
	NETSTAT      = "/proc/net/netstat"
	SOCKSTAT     = "/proc/net/sockstat"
	BUDDYINFO    = "/proc/buddyinfo"
	PROCSTAT     = "/sys/fs/cgroup/cpuacct/cpuacct.proc_stat"
	MEMCGSTAT    = "/sys/fs/cgroup/memory/memory.stat"
	MEMCGFAILCNT = "/sys/fs/cgroup/memory/memory.failcnt"
	CGTASKS      = "/sys/fs/cgroup/freezer/tasks"
	SOFTIRQS     = "/proc/softirqs"
	INTERRUPT    = "/proc/interrupts"
	IO_FILE      = "/proc/diskstats"

	ROOT_CPU_CGROUP_STAT = "/sys/fs/cgroup/cpu/cpuacct.proc_stat_v2"
)

func GetFileContent1MBAsStringLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		agentLog.WithError(err).Errorf("open file: %v ", filePath)
		return nil, err
	}
	defer file.Close()

	var chunk []byte
	buffer := make([]byte, 1024*1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		agentLog.WithError(err).Errorf("read file: %v ", filePath)
		return nil, err
	}
	chunk = append(chunk, buffer[:n]...)

	s := string(chunk)
	result := []string{}
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		result = append(result, lineStr)
	}
	return result, nil
}

func getPodBuddyStats() (map[uint16]*pb.PodBuddyStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(BUDDYINFO)
	if err != nil {
		return nil, err
	}
	pbs := make(map[uint16]*pb.PodBuddyStats)
	for _, eachLine := range strLines {
		if strings.HasPrefix(eachLine, "Node") {
			fields := strings.Fields(eachLine)
			node, err := strconv.Atoi(strings.Trim(fields[1], ","))
			if err != nil {
				return nil, err
			}
			tn := uint16(node)

			statInt32Values := []int32{}
			for _, each := range fields[4:] {
				intValue, err := strconv.Atoi(each)
				if err != nil {
					return nil, err
				}
				statInt32Values = append(statInt32Values, int32(intValue))
			}
			if len(statInt32Values) < 11 {
				return nil, fmt.Errorf("invalid buddyinfo")
			}

			if _, ok := pbs[tn]; !ok {
				pbs[tn] = &pb.PodBuddyStats{
					Order0:  0,
					Order1:  0,
					Order2:  0,
					Order3:  0,
					Order4:  0,
					Order5:  0,
					Order6:  0,
					Order7:  0,
					Order8:  0,
					Order9:  0,
					Order10: 0,
				}
			}
			pbs[tn].Order0 += statInt32Values[0]
			pbs[tn].Order1 += statInt32Values[1]
			pbs[tn].Order2 += statInt32Values[2]
			pbs[tn].Order3 += statInt32Values[3]
			pbs[tn].Order4 += statInt32Values[4]
			pbs[tn].Order5 += statInt32Values[5]
			pbs[tn].Order6 += statInt32Values[6]
			pbs[tn].Order7 += statInt32Values[7]
			pbs[tn].Order8 += statInt32Values[8]
			pbs[tn].Order9 += statInt32Values[9]
			pbs[tn].Order10 += statInt32Values[10]
		}
	}
	return pbs, err
}

func getPodCpuStats() ([]*pb.PodCpuStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(STAT)
	if err != nil {
		return nil, err
	}
	pcs := []*pb.PodCpuStats{}
	totalCpuStat := &pb.PodCpuStats{}
	for _, eachLine := range strLines {
		if strings.HasPrefix(eachLine, "cpu") {
			fields := strings.Fields(eachLine)
			statUint64Values := []uint64{}
			for _, eachField := range fields[1:] {
				v, err := strconv.ParseUint(eachField, 10, 64)
				if err != nil {
					return nil, err
				}
				statUint64Values = append(statUint64Values, v)
			}
			if len(statUint64Values) < 10 {
				return nil, fmt.Errorf("invalid cpu stat info")
			}
			tmpPcs := &pb.PodCpuStats{
				CpuUser:    statUint64Values[0],
				CpuNice:    statUint64Values[1],
				CpuSys:     statUint64Values[2],
				CpuIdle:    statUint64Values[3],
				CpuIowait:  statUint64Values[4],
				CpuSteal:   statUint64Values[5],
				CpuHardirq: statUint64Values[6],
				CpuSoftirq: statUint64Values[7],
				CpuGuest:   statUint64Values[8],
			}
			if fields[0] == "cpu" {
				totalCpuStat = tmpPcs
			} else {
				pcs = append(pcs, tmpPcs)
			}
		}
	}
	totalCpuStat.CpuNumber = uint64(len(pcs))
	pcs = append(pcs, totalCpuStat)
	return pcs, err
}

func getPodIOStats() ([]*pb.PodIOStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(IO_FILE)
	if err != nil {
		return nil, err
	}
	pis := []*pb.PodIOStats{}
	for _, eachLine := range strLines {
		if eachLine == "" {
			continue
		}

		fields := strings.Fields(eachLine)

		major, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		minor, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		name := fields[2]
		pi := &pb.PartInfo{
			Major: uint32(major),
			Minor: uint32(minor),
			Name:  name,
		}

		statUint64Values := []uint64{}
		for _, eachField := range fields[3:] {
			v, err := strconv.ParseUint(eachField, 10, 64)
			if err != nil {
				return nil, err
			}
			statUint64Values = append(statUint64Values, v)
		}
		if len(statUint64Values) < 10 {
			return nil, fmt.Errorf("invalid disk stat info")
		}
		bi := &pb.BlkIOInfo{
			RdIos:     statUint64Values[0],
			RdMerges:  statUint64Values[1],
			RdSectors: statUint64Values[2],
			RdTicks:   statUint64Values[3],
			WrIos:     statUint64Values[4],
			WrMerges:  statUint64Values[5],
			WrSectors: statUint64Values[6],
			WrTicks:   statUint64Values[7],
			Ticks:     statUint64Values[8],
			Aveq:      statUint64Values[9],
		}
		tmp := &pb.PodIOStats{
			Part:  pi,
			Blkio: bi,
		}
		pis = append(pis, tmp)
	}
	return pis, err
}

func getPodIrqStats() ([]*pb.PodIrqStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(INTERRUPT)
	if err != nil {
		return nil, err
	}
	pis := []*pb.PodIrqStats{}
	for _, eachLine := range strLines {
		if eachLine == "" {
			continue
		}

		fields := strings.Fields(eachLine)

		in, err := strconv.Atoi(strings.Trim(fields[0], ":"))
		if err != nil {
			continue
		}

		ic := []uint64{}
		for _, eachField := range fields[1 : 1+runtime.NumCPU()] {
			v, err := strconv.ParseUint(eachField, 10, 64)
			if err != nil {
				return nil, err
			}
			ic = append(ic, v)
		}

		tmp := &pb.PodIrqStats{
			IrqNum:   uint32(in),
			IrqCount: ic,
		}
		pis = append(pis, tmp)
	}
	return pis, err
}

func getPodLoadStats() (*pb.PodLoadStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(LOADAVG)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(strLines[0])

	tmpSlice := []uint32{}
	for _, eachField := range fields[:3] {
		v, err := strconv.ParseFloat(eachField, 64)
		if err != nil {
			return nil, err
		}
		tmpSlice = append(tmpSlice, uint32(v*100))
	}

	pls := &pb.PodLoadStats{
		LoadAvg_1:  tmpSlice[0],
		LoadAvg_5:  tmpSlice[1],
		LoadAvg_15: tmpSlice[2],
	}

	a := strings.Split(fields[3], "/")
	nrr, err := strconv.ParseUint(a[0], 10, 64)
	if err != nil {
		return nil, err
	}
	nrt, err := strconv.ParseUint(a[1], 10, 64)
	if err != nil {
		return nil, err
	}
	if nrr > 0 {
		// Do not take current process into account
		nrr -= 1
	}
	pls.NrRunning = nrr
	pls.NrThreads = uint32(nrt)

	strLines, err = GetFileContent1MBAsStringLines(ROOT_CPU_CGROUP_STAT)
	if err != nil {
		return nil, err
	}
	for _, eachLine := range strLines {
		if strings.HasPrefix(eachLine, "running load average(1min)") {
			v, err := strconv.Atoi(strings.Fields(eachLine)[3])
			if err != nil {
				return nil, err
			}
			pls.RunningLoadAvg_1 = uint32(v)
		}
		if strings.HasPrefix(eachLine, "running load average(5min)") {
			v, err := strconv.Atoi(strings.Fields(eachLine)[3])
			if err != nil {
				return nil, err
			}
			pls.RunningLoadAvg_5 = uint32(v)
		}
		if strings.HasPrefix(eachLine, "running load average(15min)") {
			v, err := strconv.Atoi(strings.Fields(eachLine)[3])
			if err != nil {
				return nil, err
			}
			pls.RunningLoadAvg_15 = uint32(v)
		}
	}

	return pls, err
}

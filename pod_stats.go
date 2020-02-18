package main

import (
	"fmt"
	"io"
	"os"
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

package main

import (
	"fmt"
	"strconv"
	"strings"

	pb "github.com/kata-containers/agent/protocols/grpc"
)

func getProcessProcCpuStats(Pid int) (*pb.ProcessProcCpuStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(fmt.Sprintf("/proc/%d/stat", Pid))
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(strLines[0])

	if len(fields) < 42 {
		return nil, fmt.Errorf("invalid process stat info")
	}

	tmpSlice := []uint64{}
	for _, eachField := range [5]string{fields[13], fields[15], fields[14], fields[16], fields[41]} {
		v, err := strconv.ParseUint(eachField, 10, 64)
		if err != nil {
			return nil, err
		}
		tmpSlice = append(tmpSlice, v)
	}

	ppcs := &pb.ProcessProcCpuStats{
		UserCpu:    tmpSlice[0] + tmpSlice[1],
		SysCpu:     tmpSlice[2] + tmpSlice[3],
		BlkioDelay: tmpSlice[4],
	}

	strLines, err = GetFileContent1MBAsStringLines(fmt.Sprintf("/proc/%d/status", Pid))
	if err != nil {
		return nil, err
	}
	for _, eachLine := range strLines {
		if strings.HasPrefix(eachLine, "Threads:") {
			v, err := strconv.ParseUint(strings.Fields(eachLine)[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ppcs.NThreads = v
			break
		}
	}

	strLines, err = GetFileContent1MBAsStringLines(fmt.Sprintf("/proc/%d/schedstat", Pid))
	if err != nil {
		return nil, err
	}
	fields = strings.Fields(strLines[0])
	oncpu, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return nil, err
	}
	wait_on_rq, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return nil, err
	}
	ppcs.RunOnCpu = oncpu
	ppcs.WaitOnRq = wait_on_rq

	return ppcs, err
}

func getProcessProcIoStats(Pid int) (*pb.ProcessProcIOStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(fmt.Sprintf("/proc/%d/io", Pid))
	if err != nil {
		return nil, err
	}

	tmpMap := make(map[string]uint64)
	for _, eachLine := range strLines {
		if eachLine == "" {
			continue
		}
		fields := strings.Fields(eachLine)
		v, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, err
		}
		tmpMap[fields[0]] = v
	}
	ppis := &pb.ProcessProcIOStats{
		Rchar:  tmpMap["rchar:"],
		Wchar:  tmpMap["wchar:"],
		Rbytes: tmpMap["read_bytes:"],
		Wbytes: tmpMap["write_bytes:"],
		Syscr:  tmpMap["syscr:"],
		Syscw:  tmpMap["syscw:"],
	}

	strLines, err = GetFileContent1MBAsStringLines(fmt.Sprintf("/proc/%d/stat", Pid))
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(strLines[0])
	if len(fields) < 43 {
		return nil, fmt.Errorf("invalid process stat info")
	}
	v, err := strconv.ParseUint(fields[41], 10, 64)
	if err != nil {
		return nil, err
	}
	ppis.Iolat = v

	return ppis, err
}

func getProcessProcMemStats(Pid int) (*pb.ProcessProcMemStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(fmt.Sprintf("/proc/%d/status", Pid))
	if err != nil {
		return nil, err
	}

	tmpMap := make(map[string]uint64)
	for _, eachLine := range strLines {
		fields := strings.Fields(eachLine)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "VmRSS:" ||
			fields[0] == "RssFile:" ||
			fields[0] == "RssAnon:" ||
			fields[0] == "RssShmem:" ||
			fields[0] == "VmLck:" ||
			fields[0] == "VmSwap:" {
			v, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return nil, err
			}
			tmpMap[fields[0]] = v
		}
	}
	ppms := &pb.ProcessProcMemStats{
		Usage: tmpMap["VmRSS:"],
		Cache: tmpMap["RssFile:"],
		Anon:  tmpMap["RssAnon:"],
		Shmem: tmpMap["RssShmem:"],
		Lock:  tmpMap["VmLck:"],
		Swp:   tmpMap["VmSwap:"],
	}

	strLines, err = GetFileContent1MBAsStringLines("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	for _, eachLine := range strLines {
		if strings.HasPrefix(eachLine, "SwapTotal:") {
			v, err := strconv.ParseUint(strings.Fields(eachLine)[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ppms.TotalSwp = v
			break
		}
	}

	return ppms, err
}

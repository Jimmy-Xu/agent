package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	pb "github.com/kata-containers/agent/protocols/grpc"
	"github.com/sirupsen/logrus"
)

const (
	PidIO        = "/proc/%d/io"
	PidFd        = "/proc/%d/fd"
	PidStat      = "/proc/%d/stat"
	PidStatus    = "/proc/%d/status"
	PidSchedStat = "/proc/%d/schedstat"
	PidCgroup    = "/proc/%d/cgroup"

	ProcMemInfo = "/proc/meminfo"

	CfsBaseDir = "/sys/fs/cgroup"
	CfsStatics = "cpuacct.sched_cfs_statistics"

	ReadSysFsBufSize  = 1024 * 512
	ReadDirentBufSize = 4096 * 25
)

// struct for /proc/{pid}/stat
type ProcStat struct {
	PID        int
	TaskName   string
	State      string
	PPID       int
	PGID       int
	SID        int
	TtyNr      int
	TtyPgrp    int
	TaskFlags  uint
	MinFlt     uint
	CMinFlt    uint
	MajFlt     uint
	CMajFlt    uint
	UTime      uint
	STime      uint
	CUTime     uint
	CSTime     uint
	Priority   int
	Nice       int
	NumThreads int
	StartTime  uint64
	VSize      uint
	MMRss      int
}

// struct for /proc/{pid}/io
type ProcIO struct {
	RChar               uint64 `json:"rchar"`
	WChar               uint64 `json:"wchar"`
	SyscR               uint64 `json:"syscr"`
	SyscW               uint64 `json:"syscw"`
	ReadBytes           uint64 `json:"read_bytes"`
	WriteBytes          uint64 `json:"write_bytes"`
	CancelledWriteBytes int64  `json:"cancelled_write_bytes"`
}

// struct for /sys/fs/cgroup/cpu/.../cpuacct.sched_cfs_statistics
type CgroupSchedLine2 struct {
	BvtDelay       uint64
	NoiseKickDelay uint64
	LossTime       uint64
}
type CgroupSchedLine4 struct {
	TasksDelayMs0  uint64
	TasksDelayMs1  uint64
	TasksDelayMs2  uint64
	TasksDelayMs3  uint64
	TasksDelayMs4  uint64
	TasksDelayMs5  uint64
	TasksDelayMs6  uint64
	TasksDelayMs7  uint64
	TasksDelayMs8  uint64
	TasksDelayMs9  uint64
	TasksDelayMs10 uint64
	TasksDelayMs11 uint64
}

func getProcessProcCpuStats(Pid int) (*pb.ProcessProcCpuStats, error) {
	strLines, err := GetFileContent1MBAsStringLines(fmt.Sprintf(PidStat, Pid))
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

	strLines, err = GetFileContent1MBAsStringLines(fmt.Sprintf(PidStatus, Pid))
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

	strLines, err = GetFileContent1MBAsStringLines(fmt.Sprintf(PidSchedStat, Pid))
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
	strLines, err := GetFileContent1MBAsStringLines(fmt.Sprintf(PidIO, Pid))
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

	strLines, err = GetFileContent1MBAsStringLines(fmt.Sprintf(PidStat, Pid))
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
	strLines, err := GetFileContent1MBAsStringLines(fmt.Sprintf(PidStatus, Pid))
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

	strLines, err = GetFileContent1MBAsStringLines(ProcMemInfo)
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

func getProcessPidStats(pid int) (*pb.ProcessPidStats, error) {
	// read /proc/{pid}/stat
	ps, err := readProcStat(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to read proc stat of pid %d, error:%v", pid, err)
	}

	// read /proc/{pid}/io
	pio, err := readProcIO(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to read proc io of pid %d, error:%v", pid, err)
	}

	// calculate fd count
	fdCnt, _ := calcProcFdCount(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to get fd count of pid %d, error:%v", pid, err)
	}

	return &pb.ProcessPidStats{
		UserCpu:    uint64(ps.UTime + ps.CUTime),
		SysCpu:     uint64(ps.STime + ps.CSTime),
		Mem:        uint64(ps.MMRss * 4096),
		ReadBytes:  pio.ReadBytes,
		WriteBytes: pio.WriteBytes,
		Fdcnt:      uint64(fdCnt),
		Minflt:     uint64(ps.MinFlt),
		Majflt:     uint64(ps.MajFlt),
		Thread:     uint64(ps.NumThreads),
	}, nil
}

func getProcessCgroupSched(pid int) (*pb.ProcessCgroupSched, error) {
	cgPath, err := getCgroupPathFromPid(pid, "cpuacct")
	if err != nil {
		return nil, fmt.Errorf("failed to get cgroup path from pid %d, error:%v", pid, err)
	}

	CCfsPath := fmt.Sprintf("%s/cpu%s/%s", CfsBaseDir, cgPath, CfsStatics)

	if _, err := os.Stat(CCfsPath); err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("cgroup sched stat file %s not found", CCfsPath)
			return &pb.ProcessCgroupSched{}, nil
		} else {
			return nil, fmt.Errorf("failed to read cgroup sched stat file %s, error:%v", CCfsPath, err)
		}
	}

	line2, line4, err := readCGroupSchedCfsStat(CCfsPath)
	if err != nil {
		return nil, err
	}

	return &pb.ProcessCgroupSched{
		BvtDelay:       line2.BvtDelay,
		NoiseKickDelay: line2.NoiseKickDelay,
		LossTime:       line2.LossTime,
		TasksDelayMs: []uint64{
			line4.TasksDelayMs0,
			line4.TasksDelayMs1,
			line4.TasksDelayMs2,
			line4.TasksDelayMs3,
			line4.TasksDelayMs4,
			line4.TasksDelayMs5,
			line4.TasksDelayMs6,
			line4.TasksDelayMs7,
			line4.TasksDelayMs8,
			line4.TasksDelayMs9,
			line4.TasksDelayMs10,
			line4.TasksDelayMs11,
		},
	}, nil
}

////////////////////////////////////////////////////////////////
// utilities
////////////////////////////////////////////////////////////////
func readProcStat(pid int) (*ProcStat, error) {
	data, err := readSysFile(fmt.Sprintf(PidStat, pid))
	if err != nil {
		return nil, err
	}

	var (
		ignore int

		ps = ProcStat{PID: pid}
		l  = bytes.Index(data, []byte("("))
		r  = bytes.LastIndex(data, []byte(")"))
	)

	if l < 0 || r < 0 {
		return &ps, fmt.Errorf("unexpected format, couldn't extract comm: %s", data)
	}

	ps.TaskName = string(data[l+1 : r])
	_, err = fmt.Fscan(
		bytes.NewBuffer(data[r+2:]),
		&ps.State,
		&ps.PPID,
		&ps.PGID,
		&ps.SID,
		&ps.TtyNr,
		&ps.TtyPgrp,
		&ps.TaskFlags,
		&ps.MinFlt,
		&ps.CMinFlt,
		&ps.MajFlt,
		&ps.CMajFlt,
		&ps.UTime,
		&ps.STime,
		&ps.CUTime,
		&ps.CSTime,
		&ps.Priority,
		&ps.Nice,
		&ps.NumThreads,
		&ignore,
		&ps.StartTime,
		&ps.VSize,
		&ps.MMRss,
	)
	if err != nil {
		return &ps, err
	}

	return &ps, nil
}

func readProcIO(pid int) (*ProcIO, error) {
	pio := ProcIO{}

	data, err := readSysFile(fmt.Sprintf(PidIO, pid))
	if err != nil {
		return &pio, err
	}

	format := "rchar: %d\nwchar: %d\nsyscr: %d\nsyscw: %d\nread_bytes: %d\nwrite_bytes: %d\ncancelled_write_bytes: %d\n"
	_, err = fmt.Sscanf(string(data), format,
		&pio.RChar, &pio.WChar, &pio.SyscR, &pio.SyscW, &pio.ReadBytes, &pio.WriteBytes, &pio.CancelledWriteBytes)

	return &pio, err
}

func readCGroupSchedCfsStat(path string) (*CgroupSchedLine2, *CgroupSchedLine4, error) {
	line2 := &CgroupSchedLine2{}
	line4 := &CgroupSchedLine4{}

	data, err := readSysFile(path)
	if err != nil {
		return line2, line4, err
	}

	logrus.Debugf("%s:\n%s", path, data)

	var ignore int
	for i, line := range strings.Split(string(data), "\n") {
		if i == 1 {
			//parse line2
			_, err = fmt.Fscan(
				bytes.NewBuffer([]byte(line)),
				&line2.BvtDelay,
				&line2.NoiseKickDelay,
				&line2.LossTime,
				&ignore,
			)
		} else if i == 3 {
			//parse line4
			_, err = fmt.Fscan(
				bytes.NewBuffer([]byte(line)),
				&line4.TasksDelayMs0,
				&line4.TasksDelayMs1,
				&line4.TasksDelayMs2,
				&line4.TasksDelayMs3,
				&line4.TasksDelayMs4,
				&line4.TasksDelayMs5,
				&line4.TasksDelayMs6,
				&line4.TasksDelayMs7,
				&line4.TasksDelayMs8,
				&line4.TasksDelayMs9,
				&line4.TasksDelayMs10,
				&line4.TasksDelayMs11,
			)
		}
	}
	return line2, line4, nil
}

func readSysFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		logrus.Errorf("failed to read sysfs %s, error:%v", filename, err)
		return nil, err
	}
	defer f.Close()

	reader := io.LimitReader(f, ReadSysFsBufSize)
	return ioutil.ReadAll(reader)
}

func calcProcFdCount(pid int) (int, error) {
	buf := make([]byte, ReadDirentBufSize)
	dirPath := fmt.Sprintf(PidFd, pid)
	d, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return -1, err
		}
		if strings.Contains(err.Error(), "not a directory") {
			return -1, err
		}
		return -1, err
	}
	defer d.Close()

	fd := int(d.Fd())
	fdCnt := 0
	for {
		nbuf, err := syscall.ReadDirent(fd, buf)
		if err != nil {
			return -1, err
		}
		if nbuf <= 0 {
			break
		}

		for off := 0; off < nbuf; {
			dirent := (*syscall.Dirent)(unsafe.Pointer(&buf[off]))
			off += int(dirent.Reclen)
			if dirent.Type == syscall.DT_LNK {
				fdCnt++
			}
		}
	}

	return fdCnt, nil
}

func getCgroupPathFromPid(pid int, controller string) (string, error) {
	data, err := readSysFile(fmt.Sprintf(PidCgroup, pid))
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		cols := strings.Split(line, ":")
		if len(cols) != 3 {
			continue
		}
		if !strings.Contains(cols[1], controller) {
			continue
		}
		return cols[2], nil
	}
	return "", fmt.Errorf("no cgroup path found")
}

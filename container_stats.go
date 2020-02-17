package main

import (
	"errors"
	pb "github.com/kata-containers/agent/protocols/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MaxCpus = 2048
	MaskLen = MaxCpus / 64
	envAjdkScore = "AJDK_MAX_PROCESSORS_LIMIT"
	// cgroup files to get information
	CgroupFileProcStatV2 = "cpuacct.proc_stat_v2"
	CgroupFileProcStat = "cpuacct.proc_stat"
	CgroupFileCpuacctStat = "cpuacct.stat"
	CgroupFileCpuacctUsage = "cpuacct.usage"
	CgroupFileCpuacctStatics = "cpuacct.sched_cfs_statistics"
	CgroupFileCpuQuato = "cpu.cfs_quato_us"
	CgroupFileCpuPeriod = "cpu.cfs_period_us"
	CgroupFileCpuShares = "cpu.shares"
	CgroupFileCpuStat = "cpu.stat"
	CgroupFileCpuset = "cpuset.cpus"
	CgroupFilePidsCurrent = "pids.current"
	CgroupFileMemoryStat = "memory.stat"
	CgroupFileMemoryKmem = "memory.kmem.usage_in_bytes"
	CgroupFileMemoryFailcnt = "memory.failcnt"
	CgroupFileMemswFailcnt = "memory.memsw.failcnt"
	CgroupFileMemorySwapiness = "memory.swapiness"
	CgroupFileMemswLimit = "memory.memsw.limit_in_bytes"

	// cgroup paths indexes
	cpuacctCgroupIndex = "cpuacct"
	blkioCgroupIndex = "blkio"
	cpuCgroupIndex = "cpu"
	cpusetCgroupIndex ="cpuset"
	devicesCgroupIndex = "devices"
	freezerCgroupIndex = "freezer"
	hugetlbCgroupIndex = "hugetlb"
	memoryCgroupIndex = "memory"
	netclsCgroupIndex = "net_cls"
	netprioCgroupIndex = "net_prio"
	perfeventCgroupIndex = "perf_event"
	pidsCgroupIndex = "pids"

	// mapping indexes of cpu and memory for alikernel
	cpuUser = "user"
	cpuNice = "nice"
	cpuSys = "system"
	cpuIdle = "idle"
	cpuIowait = "iowait"
	cpuIrq = "irq"
	cpuSoftirq = "softirq"
	cpuSteal = "steal"
	cpuGuest = "guest"
	cpuLoad1 = "load average(1min)"
	cpuLoad5 = "load average(5min)"
	cpuLoad15 = "load average(15min)"
	cpuNrRtask = "nr_running"
	cpuNrDtask = "nr_uninterruptible"
	cpuSwitch = "nr_switches"
	cpuVswitch = "nr_vswitches"
	cpuIswitch = "nr_iswitches"
	cpuMigrations = "nr_migrations"
	cpuRLoad1 = "running load average(1min)"
	cpuRLoad5 = "running load average(5min)"
	cpuRLoad15 = "running load average(15min)"
	cpuQuato = "quato"
	cpuPeriod = "period"
	cpuShares = "shares"
	cpuNCpu = "ncpu"
	cpuNrTask = "nr_tasks"
	cpuNrNewTask = "nr_newtask"
	cpuTotalTime = "total_time"
	cpuNrPeriods = "nr_periods"
	cpuNrThrottled = "nr_throttled"
	cpuThrottledTime = "throttled_time"
	// memory
	memMemoryLimit = "hierarchical_memory_limit"
	memMemswLimit = "hierarchical_memsw_limit"
	memCache = "total_cache"
	memRss = "total_rss"
	memRssHuge = "total_rss_huge"
	memMappedFile = "total_mapped_file"
	memDirty = "total_dirty"
	memWriteback = "total_writeback"
	memSwap = "total_swap"
	memWorkingsetRefault = "total_workingset_refault"
	memWorkingsetActive = "total_workingset_active"
	memWorkingsetRestore = "total_workingset_restore"
	memPgpgin = "total_pgpgin"
	memPgpgout = "total_pgpgout"
	memPgFault = "total_pgfault"
	memPgmajFault = "total_pgmajfault"
	memPgOutRun = "total_pgoutrun"
	memAllocStall = "total_allocstall"
	memKswapdSteal = "total_kswapd_steal"
	memPgPgSteal = "total_pg_pgsteal"
	memKswapdPgScan = "total_kswapd_pgscan"
	memPgPgscan = "total_pg_pgscan"
	memPgRefill = "total_pgrefill"
	memInactiveAnon = "total_inactive_anon"
	memActiveAnon = "total_active_anon"
	memInativeFile = "total_inactive_file"
	memActiveFile = "total_active_file"
	memUnevictable = "total_unevictable"
	memKmem = "kmem"
	memFailNo = "failcnt"
	memSwFailNo = "memsw.failcnt"
	memSwLimit = "memsw.limit"
	memSwapiness = "swapiness"
)

var (
)

func readCgroupFileInt64(file string) (int64, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	content := string(data)
	return strconv.ParseInt(strings.Trim(content, " \t\n"), 10, 64)
}

func readCgroupFileUint64(file string) (uint64, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	content := string(data)
	return strconv.ParseUint(strings.Trim(content, " \t\n"), 10, 64)
}

func setCpuMask(start uint64, end uint64, mask []uint64) (error) {
	for i := start; i < end; i++ {
		index = i / 64
		if index >= len(mask) {
			return errors.New("cpu number big than mask length")
		}
		mask[index] |= 1 << (i % 64)
	}

	return nil
}

func parseCpuset(s string) (uint64, []uint64, error) {
	mask := make([]uint64, MaskLen)
	fields := strings.Split(s, ",")
	ncpu := 0
	var start, end uint64

	for _, field := range fields {
		cpus := strings.Split("-")
		if len(cpus) > 2 {
			agentLog.Info("invalid cpu range: " + field)
			continue
		}

		if len(cpus) == 1 {
			cpu, err := strconv.Atoi(strings.TrimSpace(cpus[0]))
			if err != nil {
				agentLog.Info("invalid cpu number: " + cpus[0])
				continue
			}
			// set mask
			start = cpu
			end = cpu
		}

		if len(cpus) == 2 {
			start, err = strconv.Atoi(strings.TrimSpace(cpus[0]))
			if err != nil {
				agentLog.Info("invalid cpu number: " + cpus[0])
				continue
			}

			end, err = strconv.Atoi(strings.TrimSpace(cpus[1]))
			if err != nil {
				agentLog.Info("invalid cpu number: " + cpus[1])
				continue
			}
		}

		err = setCpuMask(cpu, cpu, mask)
		ncpu += end - start + 1
	}

	return ncpu, mask, nil
}

func readCpuset(file string) (uint64, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	content := string(data)
	ncpu, mask, err := parseCpuset(content)

	return ncpu, err
}

func parseContent(lines string) (map[string]string) {
	data := make(map[string]string)
	for _, line := range strings.Split(lines, "\n") {
		fields := strings.Fields(line)
		length = len(fields)
		if length < 2 {
			continue
		}

		index := fields[0]
		if length > 2 {
			// last one in number, join remaing string as index
			index = strings.join(fields[:length - 1], " ")
		}

		number, err := strconv.ParseInt(fields[length - 1], 10, 64)
		if err != nil {
			continue
		}
		data[index] = number
	}
	return data
}

func readAliCgroupCpuData(paths map[string]string) (map[string]uint64, error) {
	acct, ok := paths[cpuacctCgroupIndex]
	if !ok {
		return make(map[string]uint64), status.Error(codes.Unavailable, "No cpuacct cgroup")
	}

	content, err := ioutil.ReadFile(acct + ProcStatV2)
	if err != nil {
		content, err = iouitl.ReadFile(acct + ProcStat)
		if err != nil {
			return make(map[string]uint64), status.Error(codes.NotFound, "No stat file")
		}
	}

	lines := string(content)
	data := parseContent(lines)

	// remaining files
	cpupath, ok := paths[cpuCgroupIndex]
	if !ok {
		agentLog.Info("No cpu cgroup")
		return data, nil
	}

	// cpu.cfs_quato_us
	number, err := readCgroupFileInt64(cpupath + CgroupFileCpuQuato)
	if err == nil {
		data[cpuQuato] = uint64(number)
	}

	if number > 0 {
		//  cpu.cfs_period_us
		period, err := readCgroupFileUint64(cpupath + CgroupFileCpuPeriod)
		if err == nil {
			data[cpuPeriod] = period
		}
	}

	// cpu.shares
	period, err := readCgroupFileUint64(cpupath + CgroupFileCpuShares)
	if err == nil {
		data[cpuShares] = period
	}

	// number of cpus, cpuset.cpus
	cpusetPath, ok := paths[cpusetCgroupIndex]
	if !ok {
		agentLog.Info("No cpuset cgroup")
	} else {
		period, err := readCpuset(cpusetPath + CgroupFileCpuset)
		if err == nil {
			data[cpuNCpu] = period
		}
	}

	// pids.current
	pidpath, ok = paths[pidsCgroupIndex]
	if ok {
		period, err := readCgroupFileUint64(pidpath + CgroupFilePidsCurrent)
		if err == nil {
			data[cpuNrTask] = period
		}
	}
	return data nil
}

func readCgroupCpuData(paths map[string]string) (map[string]uint64, error) {
	cpuacct, ok := paths[cpuacctCgroupIndex]
	if !ok {
		return make(map[string]uint64), status.Error(codes.NotFound, "no cpuacct cgroup")
	}

	rawdata, err := ioutil.ReadFile(cpuacct + CgroupFileCpuacctStat)
	if err != nil {
		return make(map[string]uint64), status.Error(codes.NotFound, "file cpuacct.stat not found")
	}

	content := string(rawdata)
	data := parseContent(content)

	// cpuacct.usage
	number, err := readCgroupFileUint64(cpuacct + CgroupFileCpuacctUsage)
	if err == nil {
		data[cpuTotalTime] = number
	}

	// number of cpus, cpuset.cpus
	cpusetPath, ok := paths[cpusetCgroupIndex]
	if !ok {
		agentLog.Info("No cpuset cgroup")
	} else {
		period, err := readCpuset(cpusetPath + CgroupFileCpuset)
		if err == nil {
			data[cpuNCpu] = period
		}
	}

	cpupath, ok := paths[cpuCgroupIndex]
	if !ok {
		agentLog.Info("No cpu cgroup")
		return data, nil
	}

	// cpu.cfs_quato_us
	number, err := readCgroupFileInt64(cpupath + CgroupFileCpuQuato)
	if err == nil {
		data[cpuQuato] = uint64(number)
	}

	if number > 0 {
		//  cpu.cfs_period_us
		period, err := readCgroupFileUint64(cpupath + CgroupFileCpuPeriod)
		if err == nil {
			data[cpuPeriod] = period
		}
	}

	return data, nil
}

func readAliCgroupMemoryData(cgroups map[string]string) (map[string]uint64, error) {
	memoryPath, ok := cgroups[memoryCgroupIndex]
	if !ok {
		return make(map[string]uint64), status.Error(codes.NotFound, "no memory cgroup")
	}

	rawBytes, err := ioutil.ReadFile(memoryPath + CgroupFileMemoryStat)
	if err != nil {
		return make(map[string]uint64), status.Error(codes.NotFound, "no memory stat file")
	}

	content := string(rawBytes)
	data := parseContent(content)

	// memory.kmem.usage_in_bytes
	number, err := readCgroupFileUint64(memoryPath + CgroupFileMemoryKmem)
	if err == nil {
		data[memKmem] = number
	}

	// memory.failcnt
	number, err = readCgroupFileUint64(memoryPath + CgroupFileMemoryFailcnt)
	if err == nil {
		data[memFailNo] = number
	}

	// memory.memsw.failcnt
	number, err = readCgroupFileUint64(memoryPath + CgroupFileMemswFailcnt)
	if err == nil {
		data[memswFailNo] = number
	}

	// memory.swapiness
	number, err = readCgroupFileUint64(memoryPath + CgroupFileMemorySwapiness)
	if err == nil {
		data[memSwapiness] = number
	}

	// memory.memsw.limit_in_bytes
	if number > 0 {
		number, err = readCgroupFileUint64(memoryPath + CgroupFileMemswLimit)
		if err == nil {
			data[memswLimit] = number
	}

	return data, nil
}

func (c * container) getAjdkScore() (uint64) {
	envs := c.initProcess.process.Env

	for _, env := range envs {
		fields := strings.SplitN(env, "=", 2)
		if len(fields) != 2 {
			continue
		}

		if fields[0] == envAjdkScore {
			score, err := strconv.Atoi(strings.TrimSpace(fields[1]))
			if err != nil {
				continue
			}
			return score
		}
	}

	return 0
}

func getAliCgroupElement(index string, data map[string]string) (uint64) {
	elem, ok := data[index]
	if !ok {
		agentLog.Info("element " + index +" not exists")
		return 0
	}
	return elem
}

func getContainerCgroupCpu(paths: map[string]string) (*pb.ContainerCgroupCpu, error) {
	if aliKernel() {
		data, err := readAliCgroupCpuData(paths)
		if err != nil {
			return nil, err
		}

		User := getAliCgroupElement(cpuUser, data)
		Sys := getAliCgroupElement(cpuSys, data)
		Idle := getAliCgroupElement(cpuIdle, data)
		Iowait := getAliCgroupElement(cpuIowait, data)
		Hirq := getAliCgroupElement(cpuIrq, data)
		Sirq := getAliCgroupElement(cpuSoftirq, data)
		Steal := getAliCgroupElement(cpuSteal, data)
		Nice := getAliCgroupElement(cpuNice, data)
		TotalTime := User + Sys + Idle + Iowait + Hirq +Sirq + Steal
						+ Nice
		return &pb.ContainerCgroupCpu {
			User,
			Sys,
			Nice,
			Iowait,
			Idle,
			Hirq,
			Sirq,
			Steal,
			TotalTime,
			Period: getAliCgroupElement(cpuPeriod, data)
			Quato: getAliCgroupElement(cpuQuato, data)
			Ncpu: getAliCgroupElement(cpuNCpu, data)
			Shares: getAliCgroupElement(cpuShares, data)
		}, nil
	}

	// non-ali kernel
	data, err := readCgroupCpuData(paths)
	if err != nil {
		return nil, err
	}

	return &pb.ContainerCgroupCpu {
		User: getAliCgroupelement(cpuUser, data)
		Sys: getAliCgroupelement(cpuSys, data)
		Ncpu: getAliCgroupelement(cpuNCpu, data)
		TotalTime: getAliCgroupelement(cpuTotalTime, data)
		Quato: getAliCgroupelement(cpuQuato, data)
		Period: getAliCgroupelement(cpuPeriod, data)
	}, nil
}

func getContainerCgroupCpushare(paths map[string]string) (*pb.ContainerCgroupCpushare, error) {
	var result pb.ContainerCgroupCpushare
	cpuPath, ok := paths[cpuCgroupIndex]
	if !ok {
		return nil, status.Error(codes.NotFound, "no cpu cgroup")
	}

	rawdata, err := ioutil.ReadFile(cpuPath + CgroupFileCpuStat)
	if err != nil {
		return nil, err
	}

	content := string(rawdata)
	data := parseContent(content)
	result.NrThrottled = getAliCgroupElement(cpuNrTHrottled, data)
	result.NrPeriod = getAliCgroupElement(cpuNrPeriods, data)
	result.ThrottleTime = getAliCgroupelement(cpuThrottledTime, data)

	PeriodTime, err := readCgroupFileUint64(cpuPath + CgroupFileCpuPeriod)
	if err == nil {
		result.PeriodTime = PeriodTime
	}

	cpuacct, ok := paths[cpuacctCgroupIndex]
	if !ok {
		agentLog.Info("No cpuacct cgroup")
		return &result, nil
	}

	// cpuacct.sched_cfs_statistics
	rawdata, err = ioutil.ReadFile(cpuacct + CgroupFileCpuacctStatics)
	if err != nil {
		agentLog.Info("No cpuacct.sched_cfs_statistics file")
		return &result, nil
	}

	content = string(rawdata)
	lines := strings.Split(content, "\n")

	elems := strings.Fields(lines[0])
	OnServe, err := strconv.ParseUint64(strings.Trim(elems[0]), 10, 64)
	if err == nil {
		result.OnServe = OnServe
	}

	OnCpu, err := strconv.ParseUint64(strings.Trim(elems[1]), 10, 64)
	if err == nil {
		result.OnCpu = OnCpu
	}

	GOther, err := strconv.ParseUint(strings.Trim(elems[2]), 10, 64)
	if err == nil {
		result.GOther = Gother
	}

	GSibling, err := strconv.ParseUint(strings.Trim(elems[3]), 10, 64)
	if err == nil {
		result.GSibling = GSibling
	}

	return &result, nil
}

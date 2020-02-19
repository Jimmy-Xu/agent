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
	CgroupfileCpuacctPercpu = "cpuacct.usage_percpu"
	CgroupFileCpuacctPercpuSys = "cpuacct.usage_percpu_sys"
	CgroupFileCpuacctPercpuUser = "cpuacct.usage_percpu_user"
	CgroupFileCpuacctHistgramEnable = "cpuacct.sched_lat_histgram_enable"
	CgroupFilecpuacctHistgram = "cpuacct.sched_lat_histgram_cfs"
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
	CgroupFileMemoryMemDelay = "memory.memdelay"
	CgroupFileMemoryReclaimDelay = "memory.reclaim.latency"
	CgroupFileMemoryReclaimSched = "memory.direct_reclaim_sched_time"
	CgroupFileMemoryCompactLatency = "memory.direct_compaction_latency"
	CgroupFileMemoryReclaimLatency = "memory.direct_reclaim_latency"
	CgroupFileMemoryThrottleLatency = "memory.dirty_throttle_latency"

	// blkio files
	CgroupFileBlkioServiced = "blkio.throttle.io_serviced"
	CgroupFileBlkioServiceBytes = "blkio.throttle.io_service_bytes"
	CgroupFileBlkioWaitTime = "blkio.throttle.io_wait_time"
	CgroupFileBlkioServiceTime = "blkio.throttle.io_service_time"
	CgroupFileBlkioQueued = "blkio.throttle.total_io_queued"
	CgroupFileBlkioCompleted = "blkio.throttle.io_completed"
	CgroupFileBlkioMetaWrite = "blkio.throttle.meta_write_serviced"
	CgroupFileBlkioStateChange = "blkio.throttle.state_change_counter"

	// freezer files
	CgroupFileFreezerTasks = "tasks"

	// Proc files
	ProcFileMemInfo = "/proc/meminfo"
	ProcFileVersion = "/proc/version"
	ProcFileMemDelayEnable = "/proc/memdelay_enable"
	ProcSysfsPidMax = "/proc/sys/kernel/pid_max"
	ProcFileSnmpStat = "/proc/%d/net/snmp"
	ProcFileTcpxStat = "/proc/%d/net/netstat"
	ProcFileNetDev = "/proc/%d/net/dev"
	ProcFileLocalMount = "/proc/%d/mounts"
	ProcFileGlobalMount = "/proc/mounts"

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

	// meminfo indexes
	meminfoMemTotal = "MemTotal"
	meminfoMemFree = "MemFree"
	meminfoMemAvailable = "MemAvailable"
	meminfoBuffers = "Buffers"
	meminfoCached = "Cached"
	meminfoSwapCached = "SwapCached"
	meminfoActive = "Active"
	meminfoInactive = "Inactive"
	meminfoActive(anon) = "Active(anon)"
	meminfoInactive(anon) = "Inactive(anon)"
	meminfoActive(file) = "Active(file)"
	meminfoInactive(file) = "Inactive(file)"
	meminfoUnevictable = "Unevictable"
	meminfoMlocked = "Mlocked"
	meminfoSwapTotal = "SwapTotal"
	meminfoSwapFree = "SwapFree"
	meminfoDirty = "Dirty"
	meminfoWriteback = "Writeback"
	meminfoAnonPages = "AnonPages"
	meminfoMapped = "Mapped"
	meminfoShmem = "Shmem"
	meminfoKReclaimable = "KReclaimable"
	meminfoSlab = "Slab"
	meminfoSReclaimable = "SReclaimable"
	meminfoSUnreclaim = "SUnreclaim"
	meminfoKernelStack = "KernelStack"
	meminfoPageTables = "PageTables"
	meminfoNFS_Unstable = "NFS_Unstable"
	meminfoBounce = "Bounce"
	meminfoWritebackTmp = "WritebackTmp"
	meminfoCommitLimit = "CommitLimit"
	meminfoCommitted_AS = "Committed_AS"
	meminfoVmallocTotal = "VmallocTotal"
	meminfoVmallocUsed = "VmallocUsed"
	meminfoVmallocChunk = "VmallocChunk"
	meminfoPercpu = "Percpu"
	meminfoHardwareCorrupted = "HardwareCorrupted"
	meminfoAnonHugePages = "AnonHugePages"
	meminfoShmemHugePages = "ShmemHugePages"
	meminfoShmemPmdMapped = "ShmemPmdMapped"
	meminfoFileHugePages = "FileHugePages"
	meminfoFilePmdMapped = "FilePmdMapped"
	meminfoHugePages_Total = "HugePages_Total"
	meminfoHugePages_Free = "HugePages_Free"
	meminfoHugePages_Rsvd = "HugePages_Rsvd"
	meminfoHugePages_Surp = "HugePages_Surp"
	meminfoHugepagesize = "Hugepagesize"
	meminfoHugetlb = "Hugetlb"
	meminfoDirectMap4k = "DirectMap4k"
	meminfoDirectMap2M = "DirectMap2M"
	meminfoDirectMap1G = "DirectMap1G"

	// tcp Indexes
	tcpRtoAlgorithm = "RtoAlgorithm"
	tcpRtoMin = "RtoMin"
	tcpRtoMax = "RtoMax"
	tcpMaxConn = "MaxConn"
	tcpActiveOpens = "ActiveOpens"
	tcpPassiveOpens = "PassiveOpens"
	tcpAttemptFails = "AttemptFails"
	tcpEstabResets = "EstabResets"
	tcpCurrEstab = "CurrEstab"
	tcpInSegs = "InSegs"
	tcpOutSegs = "OutSegs"
	tcpRetransSegs = "RetransSegs"
	tcpInErrs = "InErrs"
	tcpOutRsts = "OutRsts"
	tcpInCsumErrors = "InCsumErrors"

	//tcpx Indexes
	tcpxSyncookiesSent = "SyncookiesSent"
	tcpxSyncookiesRecv = "SyncookiesRecv"
	tcpxSyncookiesFailed = "SyncookiesFailed"
	tcpxEmbryonicRsts = "EmbryonicRsts"
	tcpxPruneCalled = "PruneCalled"
	tcpxRcvPruned = "RcvPruned"
	tcpxOfoPruned = "OfoPruned"
	tcpxOutOfWindowIcmps = "OutOfWindowIcmps"
	tcpxLockDroppedIcmps = "LockDroppedIcmps"
	tcpxArpFilter = "ArpFilter"
	tcpxTW = "TW"
	tcpxTWRecycled = "TWRecycled"
	tcpxTWKilled = "TWKilled"
	tcpxPAWSActive = "PAWSActive"
	tcpxPAWSEstab = "PAWSEstab"
	tcpxDelayedACKs = "DelayedACKs"
	tcpxDelayedACKLocked = "DelayedACKLocked"
	tcpxDelayedACKLost = "DelayedACKLost"
	tcpxListenOverflows = "ListenOverflows"
	tcpxListenDrops = "ListenDrops"
	tcpxTCPHPHits = "TCPHPHits"
	tcpxTCPPureAcks = "TCPPureAcks"
	tcpxTCPHPAcks = "TCPHPAcks"
	tcpxTCPRenoRecovery = "TCPRenoRecovery"
	tcpxTCPSackRecovery = "TCPSackRecovery"
	tcpxTCPSACKReneging = "TCPSACKReneging"
	tcpxTCPSACKReorder = "TCPSACKReorder"
	tcpxTCPRenoReorder = "TCPRenoReorder"
	tcpxTCPTSReorder = "TCPTSReorder"
	tcpxTCPFullUndo = "TCPFullUndo"
	tcpxTCPPartialUndo = "TCPPartialUndo"
	tcpxTCPDSACKUndo = "TCPDSACKUndo"
	tcpxTCPLossUndo = "TCPLossUndo"
	tcpxTCPLostRetransmit = "TCPLostRetransmit"
	tcpxTCPRenoFailures = "TCPRenoFailures"
	tcpxTCPSackFailures = "TCPSackFailures"
	tcpxTCPLossFailures = "TCPLossFailures"
	tcpxTCPFastRetrans = "TCPFastRetrans"
	tcpxTCPSlowStartRetrans = "TCPSlowStartRetrans"
	tcpxTCPTimeouts = "TCPTimeouts"
	tcpxTCPLossProbes = "TCPLossProbes"
	tcpxTCPLossProbeRecovery = "TCPLossProbeRecovery"
	tcpxTCPRenoRecoveryFail = "TCPRenoRecoveryFail"
	tcpxTCPSackRecoveryFail = "TCPSackRecoveryFail"
	tcpxTCPRcvCollapsed = "TCPRcvCollapsed"
	tcpxTCPBacklogCoalesce = "TCPBacklogCoalesce"
	tcpxTCPDSACKOldSent = "TCPDSACKOldSent"
	tcpxTCPDSACKOfoSent = "TCPDSACKOfoSent"
	tcpxTCPDSACKRecv = "TCPDSACKRecv"
	tcpxTCPDSACKOfoRecv = "TCPDSACKOfoRecv"
	tcpxTCPAbortOnData = "TCPAbortOnData"
	tcpxTCPAbortOnClose = "TCPAbortOnClose"
	tcpxTCPAbortOnMemory = "TCPAbortOnMemory"
	tcpxTCPAbortOnTimeout = "TCPAbortOnTimeout"
	tcpxTCPAbortOnLinger = "TCPAbortOnLinger"
	tcpxTCPAbortFailed = "TCPAbortFailed"
	tcpxTCPMemoryPressures = "TCPMemoryPressures"
	tcpxTCPMemoryPressuresChrono = "TCPMemoryPressuresChrono"
	tcpxTCPSACKDiscard = "TCPSACKDiscard"
	tcpxTCPDSACKIgnoredOld = "TCPDSACKIgnoredOld"
	tcpxTCPDSACKIgnoredNoUndo = "TCPDSACKIgnoredNoUndo"
	tcpxTCPSpuriousRTOs = "TCPSpuriousRTOs"
	tcpxTCPMD5NotFound = "TCPMD5NotFound"
	tcpxTCPMD5Unexpected = "TCPMD5Unexpected"
	tcpxTCPMD5Failure = "TCPMD5Failure"
	tcpxTCPSackShifted = "TCPSackShifted"
	tcpxTCPSackMerged = "TCPSackMerged"
	tcpxTCPSackShiftFallback = "TCPSackShiftFallback"
	tcpxTCPBacklogDrop = "TCPBacklogDrop"
	tcpxPFMemallocDrop = "PFMemallocDrop"
	tcpxTCPMinTTLDrop = "TCPMinTTLDrop"
	tcpxTCPDeferAcceptDrop = "TCPDeferAcceptDrop"
	tcpxIPReversePathFilter = "IPReversePathFilter"
	tcpxTCPTimeWaitOverflow = "TCPTimeWaitOverflow"
	tcpxTCPReqQFullDoCookies = "TCPReqQFullDoCookies"
	tcpxTCPReqQFullDrop = "TCPReqQFullDrop"
	tcpxTCPRetransFail = "TCPRetransFail"
	tcpxTCPRcvCoalesce = "TCPRcvCoalesce"
	tcpxTCPOFOQueue = "TCPOFOQueue"
	tcpxTCPOFODrop = "TCPOFODrop"
	tcpxTCPOFOMerge = "TCPOFOMerge"
	tcpxTCPChallengeACK = "TCPChallengeACK"
	tcpxTCPSYNChallenge = "TCPSYNChallenge"
	tcpxTCPFastOpenActive = "TCPFastOpenActive"
	tcpxTCPFastOpenActiveFail = "TCPFastOpenActiveFail"
	tcpxTCPFastOpenPassive = "TCPFastOpenPassive"
	tcpxTCPFastOpenPassiveFail = "TCPFastOpenPassiveFail"
	tcpxTCPFastOpenListenOverflow = "TCPFastOpenListenOverflow"
	tcpxTCPFastOpenCookieReqd = "TCPFastOpenCookieReqd"
	tcpxTCPFastOpenBlackhole = "TCPFastOpenBlackhole"
	tcpxTCPSpuriousRtxHostQueues = "TCPSpuriousRtxHostQueues"
	tcpxBusyPollRxPackets = "BusyPollRxPackets"
	tcpxTCPAutoCorking = "TCPAutoCorking"
	tcpxTCPFromZeroWindowAdv = "TCPFromZeroWindowAdv"
	tcpxTCPToZeroWindowAdv = "TCPToZeroWindowAdv"
	tcpxTCPWantZeroWindowAdv = "TCPWantZeroWindowAdv"
	tcpxTCPSynRetrans = "TCPSynRetrans"
	tcpxTCPOrigDataSent = "TCPOrigDataSent"
	tcpxTCPHystartTrainDetect = "TCPHystartTrainDetect"
	tcpxTCPHystartTrainCwnd = "TCPHystartTrainCwnd"
	tcpxTCPHystartDelayDetect = "TCPHystartDelayDetect"
	tcpxTCPHystartDelayCwnd = "TCPHystartDelayCwnd"
	tcpxTCPACKSkippedSynRecv = "TCPACKSkippedSynRecv"
	tcpxTCPACKSkippedPAWS = "TCPACKSkippedPAWS"
	tcpxTCPACKSkippedSeq = "TCPACKSkippedSeq"
	tcpxTCPACKSkippedFinWait2 = "TCPACKSkippedFinWait2"
	tcpxTCPACKSkippedTimeWait = "TCPACKSkippedTimeWait"
	tcpxTCPACKSkippedChallenge = "TCPACKSkippedChallenge"
	tcpxTCPWinProbe = "TCPWinProbe"
	tcpxTCPKeepAlive = "TCPKeepAlive"
	tcpxTCPMTUPFail = "TCPMTUPFail"
	tcpxTCPMTUPSuccess = "TCPMTUPSuccess"
	tcpxTCPDelivered = "TCPDelivered"
	tcpxTCPDeliveredCE = "TCPDeliveredCE"
	tcpxTCPAckCompressed = "TCPAckCompressed"
	tcpxTCPZeroWindowDrop = "TCPZeroWindowDrop"
	tcpxTCPRcvQDrop = "TCPRcvQDrop"
	tcpxTCPWqueueTooBig = "TCPWqueueTooBig"
	tcpxTCPFastOpenPassiveAltKey = "TCPFastOpenPassiveAltKey"
	tcpxTCPForwardRetrans = "TCPForwardRetrans"

	//udp Indexes
	udpInDatagrams = "InDatagrams"
	udpNoPorts = "NoPorts"
	udpInErrors = "InErrors"
	udpOutDatagrams = "OutDatagrams"
	udpRcvbufErrors = "RcvbufErrors"
	udpSndbufErrors = "SndbufErrors"
	udpInCsumErrors = "InCsumErrors"
	udpIgnoredMulti = "IgnoredMulti"

	// kB KiB mB MB gB GB
	factor map[string]uint64 = {
		"kB": 1024,
		"kiB": 1000,
		"KB": 1024,
		"KiB": 1000,
		"mB": 1024 * 1024,
		"MB": 1024 * 1024,
		"miB": 1000 * 1000,
		"MiB": 1000 * 1000,
		"gB": 1024 * 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"giB": 1000 * 1000 * 1000,
		"GiB": 1000 * 1000 * 1000,
		"tB": 1024 * 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
		"tiB": 1000 * 1000 * 1000 * 1000,
		"TiB": 1000 * 1000 * 1000 * 1000,
	}

	// delayIndex
	delayIndex map[string]int = {
		"0-50us:" : 0,
		"50-100us:" : 1,
		"100us-200us:" : 2,
		"200-500us:" : 3,
		"500-1000us:" : 4,
		"1-5ms:" : 5,
		"5-100ms:" : 6,
		"10-100ms:" : 7,
		"100ms-bigger:" : 8,
	}

	// delayIndex2
	delayIndex2 map[string]int = {
		"0-50us:" : 0,
		"50-100us:" : 1,
		"100us-200us:" : 2,
		"200-500us:" : 3,
		"500-1000us:" : 4,
		"1-5ms:" : 5,
		"5-100ms:" : 6,
		"10-100ms:" : 7,
		">=100ms" : 8,
	}
)

var (
	kernelVersion int32 = 0
	memDelay int32 = 0
)

func GetKernelVersion() (int32, error) {
	if kernelVersion > 0 {
		return kernelVersion, nil
	}

	rawdata, err := ioutil.ReadFile(ProcFileVersion)
	if err != nil {
		agentLog.Error("Read /proc/version failed")
		return 0, err
	}

	content := string(rawdata)
	if !strings.HasPrefix(content, "Linux version") {
		return 0, errors.New("kernel version not found")
	}

	fields := strings.Fields(content)
	// version in fields[2]
	versions := strings.Split(fields[2], ".")
	major, err := strconv.Atoi(strings.TrimSpace(versions[0]))
	if err != nil {
		return 0, errors.New("Invalid major version number")
	}

	minor, err := strconv.Atoi(strings.TrimSpace(versions[1]))
	if err != nil {
		return 0, errors.New("Invalid minor number")
	}

	kernelVersion = major * 100 + minor

	return kernelVersion, nil
}

func memDelayEnabled() (bool) {
	if memDelay > 0 {
		return true
	}

	if memDelay < 0 {
		return false
	}

	number, err := readCgroupFileInt64(ProcFileMemDelayEnable)
	if err != nil {
		if os.IsNotExist(err) {
			memDelay = -1
			return false
		}
		agentLog.Error("Cannot read /proc/memdelay_enable")
		return false
	}

	if number != 0 {
		memDelay = 1
		return true
	} else {
		memDelay = -1
		return false
	}
}

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

type blkioResult struct {
	read uint64
	write uint64
	async uint64
	sync uint64
	up uint64
	down uint64
	meta uint64
}

func parseBlkioCgroupFile(file string) (map[string]blkioResult, error) {
	result := make(map[string]blkioResult)
	rawdata, err := ioutil.ReadFile(file)
	if err != nil {
		return result, nil
	}

	content := string(rawdata)
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}

		number, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			continue
		}

		if _, ok := result[fields[0]]; !ok {
			result[fields[0]] = blkioResult {
				write: 0,
				read: 0,
				async: 0,
				sync: 0,
				up: 0,
				down: 0,
				meta: 0,
			}
		}

		if strings.HasPrefix(fields[1], "R") {
			result[fields[0]].read += number
		}

		if strings.HasPrefix(fields[1], "W") {
			result[fields[0]].write += number
		}

		if strings.HasPrefix(fields[1], "A") {
			result[fields[0]].async += number
		}

		if strings.HasPrefix(fields[1], "S") {
			result[fields[0]].sync += number
		}

		if strings.HasPrefix(fields[1], "U") {
			result[fields[0]].up += number
		}

		if strings.HasPrefix(fields[1], "D") {
			result[fields[0]].down += number
		}

		result[fields[0]].meta += number
	}

	return result, nil
}

func getContainerCgroupIo(paths map[string]string) ([]*pb.ContainerCgroupIO, error) {
	blkio, ok := paths[blkioCgroupIndex]
	if !ok {
		return []*pb.ContainerCgroupIO{}, status.Error(codes.Unavailable, "No blkio cgroup")
	}

	ioServiced, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioServiced)
	if err != nil {
		return []*pb.ContainerCgroupIO{}, err
	}

	serviceBytes, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioServiceBytes)
	if err != nil {
		return []*pb.ContainerCgroupIO{}, err
	}

	serviceTime, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioServiceTime)
	if err != nil {
		return []*pb.ContainerCgroupIO{}, err
	}

	waitTime, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioWaitTime)
	if err != nil {
		return []*pb.ContainerCgroupIO{}, err
	}

	queued, err := parseBlkiCgroupFile(blkio + CgroupFileBlkioQueued)
	if err != nil {
		return []*pb.ContainerCgroupIO{}, err
	}

	var result []*pb.ContainerCgroupIO
	for key, _ := range ioServiced {
		fields := strings.Split(key, ":")
		if len(fields) != 2 {
			continue
		}

		Major, err := strconv.ParseInt(fields[0], 10, 32)
		if err != nil {
			continue
		}

		Minor, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			continue
		}

		Rio := 0
		Wio := 0
		if elem, ok := ioServiced[key]; ok {
			Rio = elem.read
			Wio = elem.write
		}

		Rbytes := 0
		Wbytes := 0
		if elem, ok = serviceBytes[key]; ok {
			Rbytes = elem.read
			Wbytes = elem.write
		}

		IoQueued := 0
		if elem, ok = queued[key]; ok {
			IoQueued = elem.read + elem.write
		}

		WaitTime : = 0
		if elem, ok = waitTime[key]; ok {
			WaitTime = elem.read + elem.write
		}

		ServiceTime := 0
		if elem, ok = serviceTime[key]; ok {
			ServiceTime = elem.read + elem.write
		}

		result = append(result, &pb.ContainerCgroupIO {
			Rio,
			Wio,
			Rbyte,
			Wbyte,
			IoQueued,
			WaitTime,
			ServiceTime,
			Major,
			Minor,
		})
	}
	return result, nil
}

func getContainerCgroupIoTime(paths map[string]string) ([]*pb.ContainerCgroupIOTime, error) {
	blkio, ok := paths[blkioCgroupIndex]
	if !ok {
		return []*pb.ContainerCgroupIOTime{}, status.Error(codes.Unavailable, "No blkio cgroup")
	}

	serviceTime, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioServiceTime)
	if err != nil {
		return []*pb.ContainerCgroupIOTime{}, err
	}

	waitTime, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioWaitTime)
	if err != nil {
		return []*pb.ContainerCgroupIOTime{}, err
	}

	complete, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioCompleted)
	if err != nil {
		return []*pb.ContainerCgroupIOTime{}, err
	}

	var result []*pb.ContainerCgroupIOTime
	for key, _ := range serviceTime {
		fields := strings.Split(key, ":")
		if len(fields) != 2 {
			continue
		}

		Major, err := strconv.ParseInt(fields[0], 10, 32)
		if err != nil {
			continue
		}

		Minor, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			continue
		}

		ReadServiceTime := 0
		WriteServiceTime := 0
		AsyncServiceTime := 0
		SyncServiceTime : = 0
		if elem, ok := serviceTime[key]; ok {
			ReadServiceTime = elem.read
			WriteServiceTime = elem.write
			AsyncServiceTime = elem.async
			SyncServiceTime = elem.sync
		}

		ReadWaitTime := 0
		WriteWaitTime := 0
		AsyncWaitTime := 0
		SyncWaitTime := 0
		if elem, ok = waitTime[key]; ok {
			ReadWaitTime = elem.read
			WriteWaitTime = elem.write
			AsyncWaitTime = elem.async
			SyncWaitTime = elem.sync
		}

		IoReadComplete := 0
		IoWriteComplete := 0
		IoAsyncComplete := 0
		IoSyncComplete := 0
		if elem, ok = complete[key]; ok {
			IoReadComplete = elem.read
			IoWriteComplete = elem.write
			IoAsyncComplete = elem.async
			IoSyncCompltete = elem.sync
		}
		result = append(result, &pn.ContainerCgroupIOTime {
			ReadServiceTime,
			WrtieServiceTime,
			AsyncServiceTime,
			SyncServiceTime,
			ReadWaitTime,
			WriteWaitTime,
			AsyncWaitTime,
			SyncWaitTime,
			IoReadComplete,
			IoWriteComplete,
			IoAsyncComplete,
			IoSyncComplete,
			Major,
			Minor,
		})
	}
	return result, nil
}

func getContianerCgroupIox(paths map[string]string) ([]*pb.ContianerCgroupIOX, error) {
	blkio, ok := paths[blkioCgroupIndex]
	if !ok {
		return []*pb.ContainerCgroupIOX{}, status.Error(codes.Unavailable, "No blkio cgroup")
	}

	ioServiced, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioServiced)
	if err != nil {
		return []*pb.ContainerCgroupIOX{}, err
	}

	serviceBytes, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioServiceBytes)
	if err != nil {
		return []*pb.ContainerCgroupIOX{}, err
	}

	metaWrite, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioMetaWrite)
	if err != nil {
		return []*pb.ContainerCgroupIOX{}, err
	}

	stateChange, err := parseBlkioCgroupFile(blkio + CgroupFileBlkioStateChange)
	if err != nil {
		return []*pb.ContainerCgroupIOX{}, err
	}

	var result []*pb.ContainerCgroupIOX
	for key, _ := range ioServiced {
		fields := strings.Split(key, ":")
		if len(fields) != 2 {
			continue
		}

		Major, err := strconv.ParseInt(fields[0], 10, 32)
		if err != nil {
			continue
		}

		Minor, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			continue
		}

		Sio := 0
		Asio := 0
		if elem, ok := ioServiced[key]; ok {
			Sio = elem.sync
			Asio = elem.async
		}

		SyncBytes := 0
		AsyncBytes := 0
		if elem, ok := serviceBytes[key]; ok {
			SyncBytes = elem.sync
			AsyncBytes = elem.async
		}

		MetaWrites := 0
		if elem, ok := metaWrite[key]; ok {
			MetaWrites = elem.meta
		}

		IoUp := 0
		IoDown := 0
		if elem, ok = stateChange[key]; ok {
			IoUp = elem.up
			IoDown = elem.down
		}
		result = append(result, &pb.ContainerCgroupIOX {
			Sio,
			Asio,
			Syncbytes: SyncBytes,
			Asyncbytes: AsyncBytes,
			MetaWrites,
			IoUp,
			IoDown,
			Major,
			Minor,
		})
	}
	return result, nil
}

func getContianerCgroupLoad(paths map[string]string) (*pb.ContianerCgroupLoad, error) {
	var data map[string]uint64
	if !aliKernel() {
		data = readAliCgroupCpuData(paths)
	} else {
		data = readCgroupCpuData(paths)
	}

	return &pb.ContainerCgroupLoad{
		Load1: getAliCgroupElement(cpuLoad1, data),
		Load5: getAliCgroupElement(cpuLoad5, data),
		Load15: getAliCgroupElement(cpuLoad15, data),
		RLoad1: getAliCgroupElement(cpuRLoad1, data),
		RLoad5: getAliCgroupElement(cpuRLoad5, data),
		RLoad15: getAliCgroupElement(cpuRLoad15, data),
	}, nil
}

// read meminfo
func readProcMemInfo() (map[string]uint64, error) {
	rawdata, err := ioutil.ReadFile(ProcFileMemInfo)
	if err != nil {
		return make(map[string]uint64), err
	}

	data := make(map[string]uint64)
	content := string(rawdata)
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)

		key := strings.Trim(fields[0], ":")
		nf := len(fields)
		if nf != 3 && nf != 2 {
			continue
		}

		value, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
		if err != nil {
			continue
		}

		if len(fields) == 3 {
			if elem, ok := factor[strings.Trim(fields[2])]; !ok {
				agentLog.Info("Unknown factor!")
				continue
			}
			data[key] = value * elem
		}

		if len(fields) == 2 {
			data[key] = value
		}
	}

	return data, nil
}

func getContainerCgroupMem(paths map[string]string) (*pb.ContainerCgroupMem, error) {
	data, err := readAliCgroupMemoryData(paths)
	if err != nil {
		return nil, err
	}

	meminfo, err := readProcMemInfo()
	if err != nil {
		return nil, err
	}

	tmpTotal := getAliCgroupElement(memMemoryLimit, data)
	Rss := getAliCgroupElement(memRss, data)
	Cache := getAliCgroupElement(memCache, data)
	activeFile := getAliCgroupElement(memActiveFile, data)
	inactiveFile := getAliCgroupElement(memInactiveFile, data)
	Slab := getAliCgroupElement(memSlab, data)
	Swap := 0
	Map := getAliCgroupElement(memMappedFile, data)
	swapiness := getAliCgroupElement(memSwapiness, data)
	if swapiness > 0 {
		Swap = getAliCgroupElement(memSwap, data)
		tmpTotal = getAliCgroupElement(memMemswLimit, data)
	}

	sysTotal := getAliCgroupElement(meminfoMemTotal, meminfo)
	if tmpTotal > sysTotal {
		tmpTotal = sysTotal
	}

	Total := tmpTotal
	Avail := Total - Rss - Slab - Swap - Cache + activeFile + inactiveFile

	return &pb.ContainerCgroupMem {
		Total,
		Rss,
		Cache,
		Avail,
		Slab,
		Swap,
		Map,
	}, nil
}

func parseOneMdsLoad(s string) (uint64, eror) {
	fields := strings.Split(s, ".")
	if len(fields) != 2 {
		return 0, errors.New("corrupted string format")
	}

	major, err := strconv.ParseUint(strings.TrimSpace(fields[0]), 10, 64)
	if err != nil {
		return 0, err
	}

	minor, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
	if err != nil {
		return 0, err
	}

	return major * 100 + minor, nil
}

func parseMdsLoad(s string) (uint64, uint64, uint64, error) {
	fields := strings.Fields(s)
	if len(fields) != 3 {
		return 0, 0, 0, errors.New("corrupted string format")
	}

	Load1, err := parseOneMdsLoad(fields[0])
	if err != nil {
		return 0, 0, 0, err
	}

	Load5, err := parseOneMdsLoad(fields[1])
	if err != nil {
		return 0, 0, 0, err
	}

	Load15, err := parseOneMdsLoad(fields[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return Load1, Load5, Load15, err
}

func getContainerCgroupMemDelay409(cgroup string) (*pb.ContainerCgroupMemDelay, error) {
	rawdata, err := ioutil.ReadFile(cgroup + CgroupFileMemoryMemDelay)
	if err != nil {
		return nil, err
	}

	content := string(rawdata)
	lines := strings.Split(content, "\n")

	fields := strings.Fields(lines[0])
	if len(fields) != 3 {
		return nil, errors.New("Unknown data format")
	}

	AggregateTotal, err := strconv.ParseUint(strings.TrimSpace(fields[0]), 10, 64)
	if err != nil {
		return nil, err
	}

	AggregateDirect, err : = strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
	if err != nil {
		return nil, err
	}

	AggregateBackgroud, err := strconv.ParseUint(strings.TrimSpace(fields[2]), 10, 64)
	if err != nil {
		return nil, err
	}

	MdsSomeLoad1, MdsSomeLoad5, MdsSomeLoad15, err := parseMdsLoad(lines[1])
	if err != nil {
		return nil, err
	}

	MdsFullLoad1, MdsFullLoad5, MdsFullLoad15, err := parseMdsLoad(lines[2])
	if err != nil {
		return nil, err
	}

	return &pb.ContainerCgroupMemDelay {
		AggregateTotal,
		AggregateDirect,
		AggregateBackgroud,
		MdsSomeLoad1,
		MdsSomeLoad5,
		MdsSomeLoad15,
		MdsFullLoad1,
		MdsFullLoad5,
		MdsFullLoad15,
	}, nil
}

func parseReclaimDelay(s string) ([]uint64, error) {
	var data []uint64
	rawdata, err := ioutil.ReadFile(s + CgroupFileMemoryReclaimDelay)
	if err != nil {
		return data, nil
	}

	length := len(delayIndex)
	data = make([]uint64, length)

	content := string(rawdata)
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		value, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
		if err != nil {
			agentLog.Error("cannot get memdelay's value")
			continue
		}

		if index, ok := delayIndex[fields[0]]; !ok {
			agentLog.Error("corrupted data, not in delayIndex")
			continue
		}

		data[index] = value
	}

	return data, nil
}

func parseOneDirectReclaimSched(s string) (uint64, error) {
	fields := strings.Fields(s)
	if len(fields) != 6 {
		return 0, errors.New("corrupted string!")
	}

	return strconv.ParseUint(strings.TrimSpace(fields[4]), 10, 64)
}

func parseDirectReclaimSched(s string) (uint64, uint64, uint64, error) {
	rawdata, err := ioutil.ReadFile(s + CgroupFileMemoryReclaimSched)
	if err != nil {
		return 0, 0, 0, err
	}

	content := string(rawdata)
	lines := strings.Split(content, "\n")
	WaitIf, err := parseOneDirectReclaimSched(lines[1])
	if err != nil {
		return 0, 0, 0, err
	}

	WbWait, err := parseOneDirectReclaimSched(lines[2])
	if err != nil {
		return 0, 0, 0, err
	}

	WaitOnPageWb, err := parseOneDirectRecaimSched(lines[3])
	if err != nil {
		return 0, 0, 0, err
	}

	return WaitIf, WbWait, WaitOnPageWb, nil
}

func getContainerCgroupMemDelayOther(cgroup string) (*pb.ContainerCgroupMemDelay, error) {
	DelayData, err := parseReclaimDelay(cgroup)
	if err != nil {
		return nil, err
	}

	WaitIf, WbWait, WaitOnPageWb, err := parseDirectReclaimSched(cgroup)
	if err != nil {
		return nil ,err
	}

	return &pb.ContianerCgroupMemDelay {
		DelayData,
		WaitIf,
		WbWait,
		WaitOnPageWb,
	}, nil
}

func getContianerCgroupMemDelay(paths map[string]string) (*pb.ContianerCgroupMemDelay, error) {
	memory, ok := paths[memoryCgroupIndex]
	if !ok {
		return nil, errors.New("No memory cgroup")
	}

	if GetKernelVersion() == 409 {
		if memDelayEnabled() {
			return getContainerCgroupMemDelay409(memory)
		}

		return nil, nil
	}

	return getContainerCgroupMemDelayOther(memory)
}

func getLatencyInfoFromFile(f, name string) (*pb.ContainerCgroupLatencyStat, error) {
	rawdata, err := ioutil.ReadFile(f)
	if err != nil {
		 return nil, err
	}

	content := string(rawdata)
	data := make([]uint64, len(delayIndex2))
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		length := len(fields)

		if length != 3 && length != 2 {
			agentLog.Info("mal-formatted line")
			continue
		}

		index, ok := delayIndex2[strings.TrimSpace(fields[0])]
		if !ok {
			agentLog.Error("mal-formed data? no Index in delayIndex2")
			continue
		}

		value1, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
		if err != nil {
			continue
		}

		value2 := 0
		if length == 3 {
			value2, err := strconv.ParseUint(strings.TrimSpace(fields[2]), 10, 64)
			if err != nil {
			continue
			}
		}

		data[index] = value1 + value2
	}

	return &pb.ContainerCgroupLatencyStats {
		Name: name,
		Ticks: data,
	}, nil
}

func getContainerCgroupLatencyStats(paths map[string]string) ([]*pb.ContianerCgroupLatencyStat, error) {
	memory, ok := paths[memoryCgroupIndex]
	if !ok {
		return []*pb.ContainerCgroupLatencyStats{}, errors.New("no memory cgroup")
	}

	var result []*pb.ContainerCgroupLatencyStats
	compact, err := getLatencyInfoFromFile(memory + CgroupFileMmeoryCompactLatency, "compact")
	if err == nil {
		result = append(result, compact)
	}

	reclaim, err := getLatencyInfoFromFile(memory + CgroupFileMemoryReclaimLatency, "reclaim")
	if err == nil {
		result = append(result, reclaim)
	}

	throttle, err := getLatencyInfoFromFile(memory + CgroupFileMemoryThrottleLatency, "throttle")
	if err == nil {
		result = append(result, throttle)
	}

	return result, nil
}

func getContainerCgroupMemx(paths map[string]string) (*pb.ContainerCgroupMemx, error) {
	data, err := readAliCgroupMemoryData(paths)
	if err != nil {
		return nil, err
	}

	return &pb.ContainerCgroupMemx {
		Aanon: getAliCgroupElement(memActiveAnon,data),
		Ianon: getAliCgroupElement(memInactiveAnon, data),
		Afile: getAliCgroupElement(memActiveFile, data),
		Ifile: getAliCgroupElement(memInactiveFile, data),
		Mpfile: getAliCgroupElement(memMappedFile, data),
		Dirty: getAliCgroupElement(memDirty, data),
		Wback: getAliCgroupElement(memWriteback, data),
		Slab: getAliCgroupElement(memKmem, data),
		Lock: getAliCgroupelement(memUnevictable),
		Rss: getAliCgroupElement(memRss, data),
		Huge: getAliCgroupElement(memRssHuge, data),
		Swap: getAliCgroupElement(memSwap, data),
	}, nil
}

func parseCgroupFreezerTasks(cgroup string) (uint64, []uint64, error) {
	rawdata, err := ioutil.ReadFile(cgroup + CgroupFileFreezerTasks)
	if err != nil {
		return 0, nil, err
	}

	content := string(rawdata)
	var result []uint64
	var ntasks uint64 = 0

	for _, line := range strings.Split(content, "\n") {
		pid, err := strconv.ParseUint(strings.TrimSpace(line), 10, 64)
		if err != nil {
			agentLog.Error("Cannot parse pid!")
			continue
		}

		ntasks++
		result = append(result, pid)
	}

	return ntasks, result, nil
}

func getContianerCgroupPcsw(paths map[string]string) (*pb.ContainerCgroupPcsw, error) {
	data, err := readAliCgroupCpuData(paths)
	if err != nil {
		return nil, err
	}

	maxPid, err := readCgroupFileUint64(ProcSysfsPidMax)
	if err != nil {
		return nil, err
	}

	NTasks := getAliCgroupElement(cpuNrTask, data)
	NRunning := getaliCgroupElement(cpuNrRtask, data)
	NDrunning := getAliCgroupElement(cpuNrDtask, data)
	NSwch := getAliCgroupElement(cpuSwitch, data)
	NewTask := getAliCgroupElement(cpuNrNewTask, data)
	MaxPid := maxPid
	NrVolSwch := getAliCgroupElement(cpuVswitch, data)
	NrInvolSwch := getAliCgrou[pElement(cpuIswitch, data)

	if NTasks == 0 {
		if freezer, ok := paths[freezerCgroupIndex]; ok {
			tasks, _, err := parseCgroupFreezerTasks(freezer)
			if err == nil {
				NTasks = tasks
			}
		}
	}

	return &pb.ContainerCgroupPcsw {
		NTasks,
		NRunning,
		NDrunning,
		NSwch,
		NewTask,
		MaxPid,
		NrVolSwch,
		NrInvolSwch,
	}, nil
}

func parseCpuacctPercpuStats(file string) ([]uint64, error) {
	rawdata, err := ioutil.ReadFile(file)
	if err != nil {
		return nil ,err
	}

	var result []uint64
	content := string(rawdata)
	for _, field := range strings.Fields(content) {
		number, err := strconv.ParseUint(strings.TrimSpace(field), 10, 64)
		if err != nil {
			agentLog.WithField("file", file).Error("Cannot parse cpuacct percpu stats")
			continue
		}

		result = append(result, number)
	}
}

const (
	CpuCpuset = "/sys/fs/cgroup/cpu/cpuset.cpus"
	CpusetCpuset = "/sys/fs/cgroup/cpuset/cpuset.cpus"
)

func getTotalCpus() (int32, error) {
	rawdata, err := ioutil.ReadFile(cpuCpuset)
	if err != nil {
		rawdata, err = ioutil.ReadFile(CpusetCpuset)
		if err != nil {
			return 0, err
		}
	}

	content := string(rawdata)
	fields := strings.Split(content, "-")
	if len(fields) != 2 {
		return 0, errors.New("Cannot parse cpu.cpuset")
	}

	start, err := strconv.Atoi(strings.TrimSpace(fields[0]))
	if err != nil {
		return 0, err
	}

	end, err := strconv.Atoi(strings.TrimSpace(fields[1]))
	if err != nil {
		return 0, err
	}

	return end - start + 1, nil
}

func getContianerCgroupPercpu(paths map[string]string) (*pb.ContianerCgroupPercpu, error) {
	cpuacct, ok := paths[cpuacctCgroupIndex]
	if !ok {
		return nil, errors.New("No cpuacct cgroup")
	}

	usage, err := parseCpuacctPercpuStats(cpuacct + CgroupFileCpuacctPercpu)
	if err != nil {
		return nil, err
	}

	sys, err := parseCpuacctPercpuStats(cpuacct + CgroupFileCpuacctPercpuSys)
	if err != nil {
		return nil, err
	}

	user, err := parseCpuacctPercpustats(cpuacct + CgroupFileCpuacctPercpuUser)
	if err != nil {
		return nil, err
	}

	var result []*pb.ContainerPercpuStats
	for i, _ := range usage {
		CpuUser := user[i]
		CpuSys := sys[i]
		CpuUsage := usage[i]
		CpuName := fmt.Sprintf("cpu%d", i)

		result = append(result, &pb.ContainerPercpuStats {
			CpuUser,
			CpuSys,
			CpuUsage,
			CpuName,
		})
	}

	Ncpu := 0
	ncpu, err := getTotalCpus()
	if err == nil {
		Ncpu = ncpu
	}

	return &pb.ContainerCgroupPercpu {
		Ncpu,
		PercpuStats: result,
	}, nil
}

func stringToArrayUint64(s string) ([]uint64, error) {
	var result []uint64
	for _, field := range strings.Fields(s) {
		number, err := strconv.ParseUint(strings.TrimSpace(field), 10, 64)
		if err != nil {
			agentLog.Error("Cannot parse string")
			return nil, err
		}

		result = append(result, number)
	}

	return result, nil
}

func getContianerCgroupSchedbvt(paths map[string]string) (*pb.ContianerCgroupSchedbvt, error) {
	kver, err:= GetKernelVersion()
	if err != nil {
		return nil, errors.New("Cannot get kernel version")
	}

	if kver < 409 {
		return nil, nil
	}

	cpuacct, ok := paths[cpuacctCgroupIndex]
	if !ok {
		return nil, errors.New("No cpuacct cgroup")
	}

	rawdata, err := ioutil.ReadFile(cpuacct + CgroupFileCpuacctStatics)
	if err != nil {
		return nil, err
	}

	content := string(rawdata)
	lines := strings.Split(content, "\n")

	if len(lines) < 2 {
		return nil, errore.New("No data")
	}

	data, err := stringToArrayUint64(lines[1])
	if err != nil {
		return nil, err
	}

	if len(data) < 3 {
		return nil, errors.New("No enough data")
	}

	Ncpu := 0
	cpuset, ok := paths[cpusetCgroupIndex]
	if ok {
		ncpu, err := readCpuset(cpuset + CgroupFileCpuset)
		if err == nil {
			Ncpu = ncpu
		}
	}

	return &pb.ContainerCgroupSchedbvt{
		BatchDelayStrict: data[0],
		StrictNoiseCleanBatch: data[1],
		BatchNoiseCleanDelay: data[2],
		Ncpu,
	}, nil
}

func getContainerCgroupSchedcg(paths map[string]string) *pb.ContianerCgroupSchedcg, error) {
	kver, err:= GetKernelVersion()
	if err != nil {
		return nil, errors.New("Cannot get kernel version")
	}

	if kver < 409 {
		return nil, nil
	}

	cpuacct, ok := paths[cpuacctCgroupIndex]
	if !ok {
		return nil, errors.New("No cpuacct cgroup")
	}

	rawdata, err := ioutil.ReadFile(cpuacct + CgroupFileCpuacctStatics)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	if len(fields) < 3) {
		return nil, errors.New("No data")
	}

	data, err := stringToArrayUint64(lines[2])
	if err != nil {
		return nil, err
	}

	if len(data) < 12 {
		return nil, errors.New("No enough data")
	}

	return &pb.ContainerCgroupSchedcg {
		CgDelayMs: data[:12],
	}, nil
}

func testAndEnableCfsSched(cgroup string) (error){
	rawdata, err := ioutil.ReadFile(cgroup + CgroupFileCpuacctHistgramEnable)
	if err != nil {
		return err
	}

	content := string(rawdata)

	if strings.HasPrefix(content, "1") {
		return nil
	}

	// enable it
	return ioutil.WriteFile(cgroup + CgroupFileCpuacctHistgramEnable, []byte("1\n"))
}

func getCgroupSchedTask409(cgroup string) ([]uint64, error) {
	rawdata,err := ioutil.ReadFile(cgroup + CgroupFileCpuacctStatics)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(rawdata), "\n")
	if len(lines) < 4 {
		return nil, errors.New("No enough data")
	}

	return stringToArrayUint64(lines[3])
}

const (
	DefaultLength = 12
)

func getCgroupSchedTaskOther(cgroup string) ([]uint64, error) {
	rawdata, err := ioutil.ReadFile(cgroup + CgroupFileCpuacctHistgram)
	if err != nil {
		return nil, err
	}

	content := string(rawdata)
	var data = make([]uint64, DefaultLength)
	for _, line := range strings.Split(content, "\n") {
		var _start, end int
		var value uint64

		n, err := fmt.Sscanf(line, "%dms-%dms %d", &_start, &end, &value)
		if err != nil || n != 3 {
			agentLod.Error("parse histgram failed")
			continue
		}

		index := end / 10
		for ; index > len(data); {
			data = append(data, make([]uint64], DefaultLength))
		}

		data[index] += value
	}
	return data, nil
}

func getContianerCgroupTask(patsh map[string]string) (*pb.ContainerCgroupSchedTask, error) {
	kver, err := GetKernelVersion()
	if err != nil {
		return nil, err
	}

	cpuacct, ok := paths[cpuacctCgroupIndex]
	if !ok {
		return nil, errors.New("No cpuacct gcgroup")

	if kver != 409 {
		err = testAndEnableCfsSched(cpuacct)
		if err != nil {
			return nil, err
		}
	}

	var data []uint64
	if kver == 409 {
		data, err = getCgroupSchedTask409(cpuacct)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = getCgroupSchedTaskOther(cpuacct)
		if err != nil {
			return nil, err
		}
	}

	if len(data) < 12 {
		return nil, errors.New("No enough data")
	}

	return &pb.ContainerCgroupSchedTask {
		TasksDelayMs: data[:12]
	}, nil
}

func stringsToMaps(index, data, prefix string) (map[string]uint64) {
	idx := strings.Fields(strings.TrimPrefix(index, prefix))
	nums := strings.Fields(strings.TrimPrefix(data, prefix))
	result := make(map[string]uint64)

	for i := 0; i < len(idx); i++ {
		num, err := strconv.ParseUint(strings.TrimSpace(nums[i]), 10, 64)
		if err != nil {
			continue
		}

		result[strings.TrimSpace(idx[i])] = num
	}

	return result
}

func readTcpUdpStatsFromSnmp(pid int32) (map[string]uint64, map[string]uint64, error) {
	var tcp, udp map[string]uint64
	tcp = make(map[string]uint64)
	udp = make(map[string]uint64)

	file := fmt.Sprintf(ProcFileSnmpstat, pid)

	rawdata, err := ioutil.ReadFile(file)
	if err != nil {
		return tcp, udp, err
	}

	lines := strings.Split(string(rawdata), "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "Tcp:") {
			tcp = stringsToMaps(lines[i], lines[i + 1], "Tcp:")
			i++
		}

		if strings.HasPrefix(lines[i], "Udp:") {
			udp = stringsToMaps(lines[i], lines[i + 1], "Udp:")
			i++
		}
	}

	return tcp, udp, nil
}

func getContainerCgroupTcpStats(paths map[string]string, pid int32) (*pb.ContainerCgroupTcpStats, error) {
	tcp, _, err := readTcpUdpStatsFromSnmp(pid)
	if err != nil {
		return nil, err
	}
	return &pb.ContianerCgroupTcpStats {
		ActiveOpens: getAliCgroupElement(tcpActiveOpens, tcp),
		PassiveOpens: getAliCgroupElement(tcpPassiveOpens, tcp),
		InSegs: getAliCgroupElement(tcpInSegs, tcp),
		OutSegs: getAliCgroupElement(tcpOutSegs, tcp),
		AttemptFails: getAliCgroupElement(tcpAttemtFails, tcp),
		EstabResets: getAliCgroupElement(tcpEstabResets, tcp),
		CurrEstab: getAliCgroupElement(tcpCurrEstab, tcp),
		RetransSegs: getAliCgroupElement(tcpRetansSegs, tcp),
		InErrs: getAliCgroupElement(tcpInErrs, tcp),
		OutRsts: getAliCgroupelement(tcpOutRsts, tcp),
	}, nil
}

func readTcpxStats(pid int32) (map[string]uint64, error) {
	file := fmt.Sprintf(ProcFileTcpxStat, pid)
	rawdata, err := ioutil.ReadFile(file)
	if err != nil {
		return map[string]uint64 {}, err
	}

	lines := strings.Split(string(rawdata), "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "TcpExt:") {
			tcpx := stringsToMaps(lines[i], lines[i + 1], "TcpExt:")
			return tcpx, nil
		}
	}
}

func getContianerCgroupTcpxStats(paths map[string]string, pid int32) (*pb.ContainerCgroupTcpxStats, error) {
	tcpx, err := readTcpxStats(pid)
	if err != nil {
		return nil, err
	}

	return &pb.ContainerCgroupTcpxStats {
		TcpLostRetrans: getAliCgroupElement(tcpxTCPLostRetransmit, tcpx),
		TcpSynRetrans: getAliCgroupElement(tcpxTCPSynRetrans, tcpx),
		TcpFastRetrans: getAliCgroupElement(tcpxTCPFastRetrans, tcpx),
		TcpForwardRetrans: getAliCgroupElement(tcpxTCPForwardRetrans, tcpx),
		TcpSlowstartRetrans: getAliCgroupElement(tcpxTCPSlowStartRetrans, tcpx),
		TcpListenOverFlow: getAliCgroupelement(tcpxTCPListenOverflows, tcpx),
		TcpListenDrop: getAliCgroupElement(tcpxTCPListenDrops, tcpx),
		TcpAbortData: getAliCgroupElement(tcpxTCPAbortOnData, tcpx),
		TcpAbortClose: getAliCgroupElement(tcpxTCPAbortOnClose, tcpx),
		TcpAbortMem: getAliCgroupElement(tcpxTCPAbortOnMemory, tcpx),
		TcpAbortTimeout: getAliCgroupElement(tcpxTCPAbortOnTimeout, tcpx),
		TcpAbortLinger: getAliCgroupElement(tcpxTCPAbortOnLinger, tcpx),
	}, nil
}

func parseTrafficData(content string) (*pb.ContainerCgroupTrafficStats, error) {
	Bytein := 0
	byteout := 0
	Pktin : =0
	Pktout := 0
	Pkterrin := 0
	Pktdrpin := 0
	Pkterrout := 0
	Pktdrpout := 0

	for _, line := range strings.Split(string(rawdata), "\n") {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}

		// fields[0] is ifname, fields[1] are statistics
		ifname = strings.TrimSpace(fields[0])
		if strings.Contains(ifname, "eth") ||
			strings.Contains(ifname, "em") ||
			strings.Contains(ifname, "venet") ||
			strings.Contains(ifname, "t1") ||
			strings.Contains(ifname, "t2") {
			if strings.Contains(ifname, "veth") {
				continue
			}

			data, err := stringToArrayUint64(fields[1])
			if err != nil {
				continue
			}

			if len(data) < 16 {
				continue
			}

			Bytein += data[0]
			Pktin += data[1]
			Pkterrin += data[2]
			Pktdrpin += data[3]
			Byteout += data[8]
			Pktout += data[9]
			Pkterrout += data[10]
			Pktdrpout += data[11]
		}
	}

	return &pb.ContainerCgroupTrafficStats {
		Bytein,
		Pktin,
		Pkterrin,
		Pktdrpin,
		Byteout,
		Pktout,
		Pkterrout,
		Pktdrpout,
	}, nil
}

func getContianerCgroupTrafficStats(paths map[string]string, pid int32) (*pb.ContianerCgroupTrafficStats, error) {
	rawdata, err := ioutil.ReadFile(fmt.Sprintf(ProcFileNetDev, pid))
	if err != nil {
		return nil, err
	}

	return parseTrafficData(string(rawdata))

}

func getContianerCgroupUdpStats(paths map[string]string, pid int32) (*pb.ContainerCgroupUdpStats, error) {
	_, udp, err := readTcpUdpStatsFromSnmp(pid)
	if err != nil {
		return nil, err
	}

	return &pb.ContainerCgroupUdpStats {
		InDatagrams: getAliCgroupElement(udpInDatagrams, udp),
		OutDatagrams: getAliCgroupElement(udpOutDatagrams, udp),
		NoPorts: getAliCgroupElement(udpNoPorts, udp),
		InErrors: getAliCgroupElement(udpInErrors, udp),
	}, nil
}

func getContainerCgroupVm(paths map[string]string) (*pb.ContainerCgroupVm, error) {
	data, err := readAliCgroupMemoryData(paths)
	if err != nil {
		return nil, err
	}

	Dsteal := getAliCgroupelement(memPgPgSteal, data)
	return &pb.ContainerCgroupVm {
		PageIn: getAliCgroupElement(memPgpgin, data),
		PageOut: getAliCgroupElemnt(memPgpgout, data),
		PageFault: getAliCgroupelement(memPgFault, data) + getAliCgroupElement(memPgmajFault, data),
		Steal: Dsteal + getAliCgroupElement(memKswapdSteal, data),
		Dsteal,
		Scan: getAliCgroupElement(memPgPgScan, data) + getAliCgroupElement(memKswapdScan, data),
		Dirty: getAliCgroupElement(memDirty, data),
		Failcnt: getAliCgroupElement(memFailNo, data) + getAliCgroupElement(memSwFailNo, data),
		Allocstall: getAliCgroupElement(memAllocstall, data),
	}, nil
}

func parseMountEntries(content string) ([]*pb.ContainerPartitionStats, error) {
	var result []*pb.ContainerPartitionStats
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		mnt := strings.TrimSpace(fields[1])
		if strings.HasPrefix(mnt, "/proc") ||
			strings.HasPrefix(mnt, "/dev") ||
			strings.HasPrefix(mnt, "/run") ||
			strings.HasPrefix(mnt, "/sys") {
			continue
		}

		var buf syscall.Statfs_t
		err := syscall.Statfs(mnt, &buf)
		if err != nil {
			continue
		}

		result = append(result, &pb.ContainerPartitionStats{
			Bsize: buf.Bsize,
			Blocks: buf.Blocks,
			Bfree: buf.Bfree,
			Bavail: buf.Bavail,
			Itotal: buf.Files,
			Ifree: buf.Ffree,
		})
	}

	return result, nil
}

// Looks like have to enter mnt ns to correctly get partition
// information. what if read /proc/<main_pid>/mounts, and the 
// stat some entry to get the information?
func getContainerPartitionStats(paths map[string]string, pid int32) ([]*pb.ContainerPartitionStats, error) {
	rawdata, err := ioutil.ReadFile(fmt.Sprintf(ProcFileLocalMount))
	if err != nil {
		return nil, err
	}

	return parseMountEntries(string(rawdata))
}

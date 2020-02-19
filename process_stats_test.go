package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

const pid = 1

/*
// example
$ go test -test.run TestProcessGetPidStats
get pid by pid name: pid=1
split pidof into array pid
pid: 1
get all pid's info
read values from /proc/1/stat
strcpy_safe container_id
filename: /proc/1/stat
get io info from /proc/1/io
get_fd_count: root=/proc/1/fd/
get_fd_count: total=61
store data to asar
[go]received:   dst=root_1=1132456,243136,4931584,1,4630327808,36499460096,61,23000,68;
ProcessPidStats: {
 "user_cpu": 1132456,
 "sys_cpu": 243136,
 "mem": 4931584,
 "read_bytes": 1,
 "write_bytes": 4630327808,
 "fdcnt": 36499460096,
 "minflt": 61,
 "majflt": 23000,
 "thread": 68
}
PASS
ok  	github.com/kata-containers/agent	0.006s
*/
func TestProcessGetPidStats(*testing.T) {
	rlt, _ := getProcessPidStats(pid)
	fmt.Printf("ProcessPidStats: %v\n", toJson(rlt))
}

/*
// example
$ go test -test.run TestProcessGetCgroupSched
filename: /proc/68206/cgroup
[line] 12:cpuset:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 11:pids:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 10:memory:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 9:perf_event:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 8:rdma:/
[line] 7:cpu,cpuacct:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
clear \n at end
cgroup_path=/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794 pid=68206
path: /sys/fs/cgroup/cpu/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794/cpuacct.sched_cfs_statistics
path: /tmp/cpuacct.sched_cfs_statistics
ENTER_CGROUP_COLLECT cpu: cgroup_path=/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794 incron=0
filename: /proc/68206/cgroup
[line] 12:cpuset:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 11:pids:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 10:memory:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 9:perf_event:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 8:rdma:/
[line] 7:cpu,cpuacct:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
[line] 6:freezer:/docker/94fe916cc502b772414522486c6d162e60b1b5194b650ba9c9dc2e7c96f0a794
clear \n at end
open path: /tmp/cpuacct.sched_cfs_statistics
[line 1] 429206713052 380561407843 48645305209 0 100341758
[line 2] 0 0 0 0
[line 3] 2090 0 0 0 0 7 0 0 0 0 0 0
[line 4] 5818 0 0 0 0 11 0 0 0 0 0 0
EXIT_CGROUP_COLLECT cpu
write data
FINISH_CGROUP_COLLECT
[go]received:   dst=94fe916cc502=5818,0,0,0,0,11,0,0,0,0,0,0,0,0,0,0,0,0,0;
ProcessCgroupSched: {
 "tasks_delay_ms": [
  5818,
  0,
  0,
  0,
  0,
  11,
  0,
  0,
  0,
  0,
  0
 ]
}
PASS
ok  	github.com/kata-containers/agent	0.017s
*/
func TestProcessGetCgroupSched(*testing.T) {
	rlt, _ := getProcessCgroupSched(pid)
	fmt.Printf("ProcessCgroupSched: %v\n", toJson(rlt))
}

/*
// example
$ go test -test.run TestProcessGetCpuStats
read /proc/stat
line: cpu  2035949 92 1107585 64583469 4409 0 140473 0 0 0
line: cpu0 733508 27 430206 21498512 922 0 13078 0 0 0
line: cpu1 626281 25 287115 21489670 1399 0 112471 0 0 0
line: cpu2 676159 39 390263 21595286 2086 0 14923 0 0 0
line: intr 279744386 15 14 0 0 0 0 0 0 1 0 0 0 15 0 0 222055 1034863 771305 139 126808637 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 35692 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0get all pid's info
read values from /proc/pid/stat
get container_id
get number of threads
get run on rq & wait on rq
[go]received:   dst=root_1=1132456,243137,22623992,7,6981564151,689072405,1;
ProcessProcCpuStats: {
 "user_cpu": 1132456,
 "sys_cpu": 243137,
 "blkio_delay": 22623992,
 "run_on_cpu": 7,
 "wait_on_rq": 6981564151,
 "n_threads": 689072405
}
PASS
ok  	github.com/kata-containers/agent	0.007s
*/
func TestProcessGetCpuStats(*testing.T) {
	rlt, _ := getProcessProcCpuStats(pid)
	fmt.Printf("ProcessProcCpuStats: %v\n", toJson(rlt))
}

/*
// example
$ go test -test.run TestProcessGetIoStats
get all pid's info
collect data
[go]received:   dst=root_1=343909606088,64996517409,4630327808,36499460096,7,66389914,13556641;
ProcessProcIOStats: {
 "rchar": 343909606088,
 "wchar": 64996517409,
 "rbytes": 4630327808,
 "wbytes": 36499460096,
 "iolat": 7,
 "syscr": 66389914,
 "syscw": 13556641
}
PASS
ok  	github.com/kata-containers/agent	0.008s
*/
func TestProcessGetIoStats(*testing.T) {
	rlt, _ := getProcessProcIoStats(pid)
	fmt.Printf("ProcessProcIOStats: %v\n", toJson(rlt))
}

/*
// example
$ go test -test.run TestProcessGetMemStats
get all pid's info
collect data
get total swap mem
[go]received:   dst=root_1=4816,3360,1456,0,0,48,2097148;
ProcessProcMemStats: {
 "usage": 4816,
 "cache": 3360,
 "anon": 1456,
 "swp": 48,
 "total_swp": 2097148
}
PASS
ok  	github.com/kata-containers/agent	0.007s
*/
func TestProcessGetMemStats(*testing.T) {
	rlt, _ := getProcessProcMemStats(pid)
	fmt.Printf("ProcessProcMemStats: %v\n", toJson(rlt))
}

func toJson(data interface{}) string {
	buf, _ := json.MarshalIndent(data, "", " ")
	return string(buf)
}

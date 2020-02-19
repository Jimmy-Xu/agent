package main

import (
	"fmt"
	"testing"
)

const pid = 1

/*
// example
$ go test -test.run TestGetProcessPidStats
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
[go]received:   dst=root_1=1132347,242799,4902912,1,4630282752,36499460096,61,21715,68;
{1132347 242799 4902912 1 4630282752 36499460096 61 21715 68 {} [] 0}
PASS
ok  	github.com/kata-containers/agent	0.008s
*/
func TestGetProcessPidStats(*testing.T) {
	rlt, _ := getProcessPidStats(pid)
	fmt.Printf("%v\n", *rlt)
}

/*
$ go test -test.run TestGetProcessCgroupSched
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
{[5818 0 0 0 0 11 0 0 0 0 0] 0 0 0 {} [] 0}
PASS
ok  	github.com/kata-containers/agent	0.008s
*/
func TestGetProcessCgroupSched(*testing.T) {
	rlt, _ := getProcessCgroupSched(pid)
	fmt.Printf("%v\n", *rlt)
}

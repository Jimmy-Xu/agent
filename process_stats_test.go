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

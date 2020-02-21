package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
)

const CMD = "docker"

var pid int

type ContainerInspect struct {
	State ContainerState
}
type ContainerState struct {
	Pid int
}

// setup
func setup() (int, string, error) {
	fmt.Printf("> prepare container\n")

	// run container
	cmd := exec.Command(CMD, "run", "-d", "busybox", "top", "-b")
	cidBuf, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run %s container, error:%v", CMD, err)
	}
	cid := strings.Split(string(cidBuf), "\n")[0]

	// inspect container
	cmdStr := fmt.Sprintf(`%s inspect %s`, CMD, cid)
	cmd = exec.Command("sh", "-c", cmdStr)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("failed to inspect container %s, error:%v", cid, err)
	} else {
		// get Pid
		var inspect []ContainerInspect
		err := json.Unmarshal(out, &inspect)
		if err != nil {
			log.Fatalf("failed to parse container inspect, error:%v", err)
		}
		if len(inspect) == 0 {
			log.Fatalf("failed to find container %s", cid)
		}
		fmt.Printf("Pid of container %s is %d\n\n", string(cid), inspect[0].State.Pid)
		return inspect[0].State.Pid, string(cid), nil
	}
	return -1, "", fmt.Errorf("failed to prepare container")
}

// teardown
func tearDown(cid string) {
	fmt.Printf("> clean container %s\n", cid)
	cmd := exec.Command("docker", "rm", "-fv", cid)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to clean container %s, error:%v", cid, err)
	}
	log.Printf("> clean container %s ok\n", cid)
}

// main
func TestMain(m *testing.M) {
	var (
		err error
		cid string
	)
	pid, cid, err = setup()
	if err != nil {
		log.Fatal(err)
	}
	defer tearDown(cid)

	m.Run()
}

//////////////////////////////////
// test cases
//////////////////////////////////
func TestProcessGetPidStats(t *testing.T) {
	showProcFile()
	rlt, _ := getProcessPidStats(pid)
	fmt.Printf("ProcessPidStats : %v\n", toJson(rlt))
}

func TestProcessGetCgroupSched(t *testing.T) {
	rlt, _ := getProcessCgroupSched(pid)
	fmt.Printf("ProcessCgroupSched : %v\n", toJson(rlt))
}

//////////////////////////////////
// utilities
//////////////////////////////////
func showProcFile() {
	// show proc/{pid}/stat
	pidStat := fmt.Sprintf(PidStat, pid)
	buf, err := readSysFile(pidStat)
	if err != nil {
		log.Fatalf("failed to read %s", pidStat)
	}
	fmt.Printf("[%s]\n%s\n", pidStat, buf)

	// show proc/{pid}/ioo
	pidIO := fmt.Sprintf(PidIO, pid)
	buf, err = readSysFile(pidIO)
	if err != nil {
		log.Fatalf("failed to read %s", pidStat)
	}
	fmt.Printf("[%s]\n%s\n", pidIO, buf)
}

func toJson(data interface{}) string {
	buf, _ := json.MarshalIndent(data, "", " ")
	return string(buf)
}

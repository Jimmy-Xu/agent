package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
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
		fmt.Printf("setup> container %s is created, pid is %d\n", string(cid), inspect[0].State.Pid)
		return inspect[0].State.Pid, string(cid), nil
	}
	return -1, "", fmt.Errorf("failed to prepare container")
}

// teardown
func tearDown(cid string) {
	cmd := exec.Command("docker", "rm", "-fv", cid)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to clean container %s, error:%v", cid, err)
	}
	fmt.Printf("teardown> container %s is removed\n", cid)
}

// main
func TestMain(m *testing.M) {
	var (
		err error
		cid string
	)

	for _, item := range os.Args {
		if item == "-test.v=true" {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}

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
func TestProcessGetProcCpuStats(t *testing.T) {
	showSysFile(fmt.Sprintf(PidStat, pid))
	showSysFile(fmt.Sprintf(PidStatus, pid))
	showSysFile(fmt.Sprintf(PidSchedStat, pid))
	rlt, err := getProcessProcCpuStats(pid)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("ProcessProcCpuStats : %v\n", toJson(rlt))
}

func TestProcessGetProcMemStats(t *testing.T) {
	showSysFile(fmt.Sprintf(PidStatus, pid))
	showSysFile(ProcMemInfo)
	rlt, err := getProcessProcMemStats(pid)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("ProcessProcMemStats : %v\n", toJson(rlt))
}

func TestProcessGetProcIoStats(t *testing.T) {
	showSysFile(fmt.Sprintf(PidIO, pid))
	showSysFile(fmt.Sprintf(PidStat, pid))
	rlt, err := getProcessProcIoStats(pid)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("ProcessProcIOStats : %v\n", toJson(rlt))
}

func TestProcessGetPidStats(t *testing.T) {
	showSysFile(fmt.Sprintf(PidStat, pid))
	showSysFile(fmt.Sprintf(PidIO, pid))
	listFd(fmt.Sprintf(PidFd, pid))
	rlt, err := getProcessPidStats(pid)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("ProcessPidStats : %v\n", toJson(rlt))
}

func TestProcessGetCgroupSched(t *testing.T) {
	rlt, err := getProcessCgroupSched(pid)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("ProcessCgroupSched : %v\n", toJson(rlt))
}

//////////////////////////////////
// utilities
//////////////////////////////////
func showSysFile(filename string) {
	buf, err := readSysFile(filename)
	if err != nil {
		log.Fatalf("failed to read %s", filename)
	}
	logrus.Debugf(`
================================
file: [%s] begin
================================
%s--------------------------------
file: [%s] end
--------------------------------

`, filename, buf, filename)
}

func listFd(path string) {

	cmd := exec.Command("ls", "-l", path)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to list file in %s, error:%v", path, err)
	}
	logrus.Debugf(`
================================
dir: [%s]
================================
%s
`, path, buf)
}

func toJson(data interface{}) string {
	var buf []byte
	if logrus.GetLevel() == logrus.DebugLevel {
		buf, _ = json.MarshalIndent(data, "", " ")
	} else {
		buf, _ = json.Marshal(data)
	}
	return string(buf)
}

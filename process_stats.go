package main

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	pb "github.com/kata-containers/agent/protocols/grpc"
)

/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <dirent.h>

// mod_pid.c
struct stats_pid {
    unsigned long long user_cpu;
    unsigned long long sys_cpu;
    unsigned long long mem;
    unsigned long long read_bytes;
    unsigned long long write_bytes;
    unsigned long long fdcnt;
    unsigned long long minflt;
    unsigned long long majflt;
    unsigned long long thread;
    char container_id[64];
};

#define PID_IO          "/proc/%u/io"
#define PID_FD          "/proc/%u/fd/"
#define PID_STAT        "/proc/%u/stat"
#define PID_STATUS      "/proc/%u/status"

#define ITEM_SPLIT  ";"
#define MAX_PROC_COLLECT (512)

#define SIZE_1K (1<<10)
#define SIZE_128K (128*SIZE_1K)
#define PROC_BUFSIZE SIZE_128K

__thread char *strtok_saveptr = NULL;
static inline char *strtok_safe(char *str, const char *delim)
{
    return strtok_r(str, delim, &strtok_saveptr);
}

static inline
char *strcpy_safe(char *dst, const char *src, unsigned int buf_size)
{
    strncpy(dst, src, buf_size - 1);
    dst[buf_size - 1] = 0;
    return dst;
}


int
store_single_pid(char *buf, int pid, struct stats_pid *st_pid, int size)
{
    printf("store data to asar\n");
    int len = 0;
    len = snprintf(buf + len, size, "%.12s_%d=%lld,%lld,%lld,%lld,%lld,%lld,%lld,%lld,%lld" ITEM_SPLIT,
            st_pid->container_id,
            pid,
            st_pid->user_cpu,
            st_pid->sys_cpu,
            st_pid->mem,
            st_pid->thread,
            st_pid->read_bytes,
            st_pid->write_bytes,
            st_pid->fdcnt,
            st_pid->minflt,
            st_pid->majflt
            );
    //if (pos < LEN_4096) buf[pos] = '\0';
    return len;
}

int get_fd_count(char *root)
{
    printf("get_fd_count: root=%s\n", root);
    DIR *dir;
    struct dirent * ptr;
    int total = 0;
    dir = opendir(root);
    if(dir == NULL) {
        return 0;
    }
    int errno = 0;
    while((ptr = readdir(dir)) != NULL) {
        if(ptr->d_type == DT_LNK)
            total++;
        else
            continue;
    }
    if(errno != 0) {
        return 0;
    }
    closedir(dir);
    printf("get_fd_count: total=%d\n", total);
    return total;
}

static void
read_pid_stats(char *parameter, char* buf)
{
    int     pid[MAX_PROC_COLLECT], nb = 0, i, pos;
    //char    buf[PROC_BUFSIZE];
    char    filename[128], line[256];
    FILE   *fp;
    struct  stats_pid st_pid;

    printf("get pid by pid name: pid=%s\n", parameter);
    if (strlen(parameter) > 200) {
        return;
    }
    char ppp[256];
    strcpy(ppp, parameter);
    printf("split pidof into array pid\n");
    char *p;
    p = strtok_safe(ppp, " ");
    while (p) {
        pid[nb] = atoi(p);
        printf("pid: %d\n", pid[nb]);
        if (pid[nb++] <= 0) {
            return;
        }
        if (nb >= 16) {
            return;
        }
        p = strtok_safe(NULL, " ");
    }

    printf("get all pid's info\n");
    pos = 0;
    for(i = 0; i < nb; ++i) {
        printf("read values from /proc/%d/stat\n", pid[i]);
        memset(&st_pid, 0, sizeof(st_pid));
#if IN_7U
        printf("get_container_id_from_pid: %d \n", pid[i]);
        if(get_container_id_from_pid(st_pid.container_id, 64, pid[i]) < 0) {
            continue;
        }
#else
        printf("strcpy_safe container_id\n");
        strcpy_safe(st_pid.container_id, "root", 64);
#endif
        sprintf(filename, PID_STAT, pid[i]);
        printf("filename: %s\n", filename);

        if ((fp = fopen(filename, "re")) == NULL) {
            printf("!!!! Should do something notify\n");
            continue;
        }
        unsigned long long cpudata[4];
        if (fgets(line, 256, fp) == NULL) {
            fclose(fp);
            continue;
        }

        if ((p = strstr(line, ")")) == NULL) {
            fclose(fp);
            continue;
        }
        if (sscanf(p, "%*s %*c %*d %*d %*d %*d %*d %*u %llu %*u %llu %*u %llu %llu %llu %llu %*u %*u %llu %*u %*u %*u %llu ",
               &st_pid.minflt, &st_pid.majflt, &cpudata[0], &cpudata[1],
               &cpudata[2], &cpudata[3], &st_pid.thread, &st_pid.mem) == EOF) {
            fclose(fp);
            continue;
        }
        st_pid.user_cpu = cpudata[0] + cpudata[2];
        st_pid.sys_cpu += cpudata[1] + cpudata[3];
        st_pid.mem = st_pid.mem * 4096;
        fclose(fp);

        printf("get io info from /proc/%d/io\n", pid[i]);
        sprintf(filename, PID_IO, pid[i]);
        if ((fp = fopen(filename, "re")) != NULL) {
            while (fgets(line, 256, fp) != NULL) {
                if (!strncmp(line, "read_bytes:", 11))
                    sscanf(line + 11, "%llu", &st_pid.read_bytes);
                else if (!strncmp(line, "write_bytes:", 12))
                    sscanf(line + 12, "%llu", &st_pid.write_bytes);
            }
            fclose(fp);
        }

        sprintf(filename, PID_FD, pid[i]);
        st_pid.fdcnt = get_fd_count(filename);
        pos += store_single_pid(buf + pos, pid[i], &st_pid, PROC_BUFSIZE - pos);
    }

}
*/
import "C"

const (
	ITEM_SPLIT = ";"
	ITEM_SPSTART = "="
	DATA_SPLIT = ","
)


func getProcessPidStats(pid int) (*pb.ProcessPidStats, error) {
	var (
		userCpu    uint64
		sysCpu     uint64
		mem        uint64
		readBytes  uint64
		writeBytes uint64
		fdCnt      uint64
		minFlt     uint64
		majFlt     uint64
		thread     uint64
	)

	// call read_pid_stats
	dst := ""
	c_src := C.CString("1")
	defer C.free(unsafe.Pointer(c_src))
	c_dst := C.CString(dst)
	//defer C.free(unsafe.Pointer(c_dst))
	C.read_pid_stats(c_src, c_dst)

	// show result
	received := C.GoString(c_dst)
	fmt.Printf("[go]received:   dst=%s\n", received)

	// convert result
	content := strings.Split(received, ITEM_SPSTART)[1]
	item := strings.Split(content,ITEM_SPLIT)[0]
	cols := strings.Split(item, DATA_SPLIT)

	userCpu, _ = strconv.ParseUint(cols[0],10,64)
	sysCpu, _ = strconv.ParseUint(cols[1],10,64)
	mem, _ = strconv.ParseUint(cols[2],10,64)
	readBytes, _ = strconv.ParseUint(cols[3],10,64)
	writeBytes, _ = strconv.ParseUint(cols[4],10,64)
	fdCnt, _ = strconv.ParseUint(cols[5],10,64)
	minFlt, _ = strconv.ParseUint(cols[6],10,64)
	majFlt, _ = strconv.ParseUint(cols[7],10,64)
	thread, _ = strconv.ParseUint(cols[8],10,64)

	return &pb.ProcessPidStats{
		UserCpu:    userCpu,
		SysCpu:     sysCpu,
		Mem:        mem,
		ReadBytes:  readBytes,
		WriteBytes: writeBytes,
		Fdcnt:      fdCnt,
		Minflt:     minFlt,
		Majflt:     majFlt,
		Thread:     thread,
	}, nil
}

func getProcessCgroupSched(pid int) (*pb.ProcessCgroupSched, error) {
	var (
		tasksDelayMs   []uint64
		bvtDelay       uint64
		noiseKickDelay uint64
		lossTime       uint64
	)
	return &pb.ProcessCgroupSched{
		TasksDelayMs:   tasksDelayMs,
		BvtDelay:       bvtDelay,
		NoiseKickDelay: noiseKickDelay,
		LossTime:       lossTime,
	}, nil
}

func getProcessProcCpuStats(pid int) (*pb.ProcessProcCpuStats, error) {
	var (
		userCpu    uint64
		sysCpu     uint64
		blkioDelay uint64
		runOnCpu   uint64
		waitOnRq   uint64
		nThreads   uint64
	)
	return &pb.ProcessProcCpuStats{
		UserCpu:    userCpu,
		SysCpu:     sysCpu,
		BlkioDelay: blkioDelay,
		RunOnCpu:   runOnCpu,
		WaitOnRq:   waitOnRq,
		NThreads:   nThreads,
	}, nil
}

func getProcessProcIoStats(pid int) (*pb.ProcessProcIOStats, error) {
	var (
		rchar  uint64
		wchar  uint64
		rbytes uint64
		wbytes uint64
		iolat  uint64
		syscr  uint64
		syscw  uint64
	)
	return &pb.ProcessProcIOStats{
		Rchar:  rchar,
		Wchar:  wchar,
		Rbytes: rbytes,
		Wbytes: wbytes,
		Iolat:  iolat,
		Syscr:  syscr,
		Syscw:  syscw,
	}, nil
}

func getProcessProcMemStats(pid int) (*pb.ProcessProcMemStats, error) {
	var (
		usage    uint64
		cache    uint64
		anon     uint64
		shmem    uint64
		lock     uint64
		swp      uint64
		totalSwp uint64
	)
	return &pb.ProcessProcMemStats{
		Usage:    usage,
		Cache:    cache,
		Anon:     anon,
		Shmem:    shmem,
		Lock:     lock,
		Swp:      swp,
		TotalSwp: totalSwp,
	}, nil
}

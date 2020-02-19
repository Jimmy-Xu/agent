package main

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	pb "github.com/kata-containers/agent/protocols/grpc"
)

/*
//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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


//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// mod_cgroup_sched.c

#include <stdint.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

//////////////////////////////////////////////////////////////////////
#ifndef _CGROUP_H_
#define _CGROUP_H_

#include <time.h>

int attach_cgroup(const char *subsystem, const char *path);
int detach_cgroup(const char *subsystem);

#ifdef CGROUP_ATTACH

#define DEFINE_INCRON \
        int incron = (conf.running_mode == RUN_CRON && cache_active())

#define ENTER_CGROUP_COLLECT(cg_subsys, cg_path, do_fork) \
    { \
        int pid, status; \
        if(do_fork) { \
            if((pid = fork()) < 0){ \
                continue; \
            } \
            if(pid != 0) { \
                waitpid(pid, &status, 0); \
                if(status != 0)  \
                    ALOG_ERR("collect subsystem: %s path: %s failed\n", cg_subsys, cg_path); \
                continue; \
            }  \
        } \
        if(attach_cgroup(cg_subsys, cg_path) < 0) { \
            if((do_fork))    \
                exit(-1); \
            else  \
                continue; \
        } \
    }

#define EXIT_CGROUP_COLLECT(cg_subsys) \
    { \
        if(detach_cgroup(cg_subsys) < 0) \
            exit(-1); \
    }

#define ABORT_THIS_CGROUP_COLLECT(do_fork, cg_subsys) \
        if(do_fork){ \
            exit(-1);   \
        } else {   \
            if(detach_cgroup(cg_subsys) < 0)  \
                exit(-1); \
            continue;   \
        }

#define FINISH_CGROUP_COLLECT(_mod, _buf, do_fork, n_bytes) \
        if((do_fork)) { \
            if((n_bytes) > 0) { \
                set_mod_record_v2((_mod), (_buf), 1); \
            } \
            exit(0); \
        }

#else

#define DEFINE_INCRON \
        int incron = 0; \
        (void)incron;

#define ENTER_CGROUP_COLLECT(cg_subsys, cg_path, do_fork)

#define ABORT_THIS_CGROUP_COLLECT(do_fork, cg_subsys) \
        continue

#define EXIT_CGROUP_COLLECT(cg_subsys)

#define FINISH_CGROUP_COLLECT(_mod, _buf, do_fork, n_bytes)

#endif

enum cgpath_type {
    CG_CPUACCT = 0,
    CG_CPUSET,
    CG_CPU,
    CG_MEM,
    CG_BLKIO,
    CG_FREEZER,
    CG_PIDS,
    CG_MAX
};

#define N_CGROUP_DATA_SIZE (40)

struct cgroup_data {
    uint64_t data[N_CGROUP_DATA_SIZE];
    time_t timestamp;
    void *priv;
};

static inline uint64_t
get_cgdata(struct cgroup_data *ca, int idx)
{
    return ca->data[idx];
}

static inline uint64_t
get_cgdata_diff(struct cgroup_data *ca1, struct cgroup_data *ca2, int idx)
{
    return ca1->data[idx] > ca2->data[idx] ? ca1->data[idx] - ca2->data[idx] : 0UL;
}

int
get_cgcpu_data(struct cgroup_data *ca, char *cgpath[CG_MAX]);

enum e_cpu_stat {
    // CPUACCT.PROC_STAT || CPUACCT.STAT
    CPU_USER = 0,
    CPU_NICE,
    CPU_SYS,
    CPU_IDLE,
    CPU_IOWAIT,
    CPU_IRQ,
    CPU_SOFTIRQ,
    CPU_STEAL,
    CPU_GUEST,
    CPU_LOAD1,
    CPU_LOAD5,
    CPU_LOAD15,
    CPU_RLOAD1,
    CPU_RLOAD5,
    CPU_RLOAD15,
    CPU_NR_RTSK,
    CPU_NR_DTSK,
    CPU_SWITCHS,
    CPU_VSWITCHS,
    CPU_ISWITCHS,
    CPU_MIGRATIONS,
	// CPU.QUOTA
    CPU_QUOTA,
	// CPU.PERIOD
    CPU_PERIOD,
	// CPU.SHARES
    CPU_SHARES,
	// CPUSET.CPUS
    CPU_NCPU,
    // NOT COLLECT FROM CPU CGROUP, FROM PIDS & FREEZER
    CPU_NR_TSK,
    // COLLECT FROM NAMESPACE ONLY
    CPU_NR_NEWTSK,
};

enum e_mem_stat {
    // MEMORY.STAT
    MEM_memory_limit = 0,
    MEM_memsw_limit,
    MEM_cache,
    MEM_rss,
    MEM_rss_huge,
    MEM_mapped_file,
    MEM_dirty,
    MEM_writeback,
    MEM_swap,
    MEM_workingset_refault,
    MEM_workingset_activate,
    MEM_workingset_restore,
    MEM_pgpgin,
    MEM_pgpgout,
    MEM_pgfault,
    MEM_pgmajfault,
    MEM_pgoutrun,
    MEM_allocstall,
    MEM_kswapd_steal,
    MEM_pg_pgsteal,
    MEM_kswapd_pgscan,
    MEM_pg_pgscan,
    MEM_pgrefill,
    MEM_inactive_anon,
    MEM_active_anon,
    MEM_inactive_file,
    MEM_active_file,
    MEM_unevictable,
    // MEMORY.KMEM.USAGE_IN_BYTES
    MEM_KMEM,
    // MEMORY.FAILNO
    MEMORY_FAILNO,
    MEMSW_FAILNO,
    MEMSW_LIMIT,
    MEMSW_SWAPPINESS,
};

int
get_cgmem_data(struct cgroup_data *ca, char *cgpath[CG_MAX]);

#endif


//////////////////////////////////////////////////////////////////////


#define MAX_GROUP 64
#define MAX_TASK 1024
#define MAX_NAME_LENGTH 128
#define CFS_STATICS "cpuacct.sched_cfs_statistics"
//#define CFS_HISTGRAM "cpuacct.sched_lat_histgram_cfs"
//#define CFS_HISTGRAM_ENABLE "cpuacct.sched_lat_histgram_enable"
#define MAX_NUM_OF_CONTAINERS 48
#define DATA_FORMAT "%.12s=%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,%lu,0,0,0,0;"

#define ISDOT(a)        (a[0] == '.' && (!a[1] || (a[1] == '.' && !a[2])))

#define MAX_PHY_DEVICES (32)
#define CG_PER_RECORD_SIZE (512) // 32 + 20 * 20 (64 bit number long)
#define MAX_CONTAINER_COLLECT (512)
#define CG_BUFSIZE (MAX_CONTAINER_COLLECT * CG_PER_RECORD_SIZE * MAX_PHY_DEVICES)

#define MAX_CGPATH_LEN (4096)

#define TABLE_SIZE (512)



struct st_cg_sched {
    uint64_t tasks_delay_ms[12];
    uint64_t bvt_delay;
    uint64_t noise_kick_delay;
    uint64_t loss_time;
    char container_id[64];
};

int
split_string(char *s, char *result[], unsigned int n_result, char c)
{
	unsigned int n = 0;
	char *p = s, *start = s;

	while(*p){
    	if(*p == c){
        	*p = 0;
        	result[n++] = start;
        	if(n == n_result){
        		break;
        	}
        	start = p + 1;
 		}
    	p++;
	}
    // last
	if(*start && n != n_result){
		result[n++] = start;
	}
	return n;
}

int
get_cgroup_path_from_pid(char *path, unsigned int size, unsigned int pid, const char *type)
{
    #define PROC_CGPATH "/proc/%u/cgroup"
    FILE *fp;
    char line[MAX_CGPATH_LEN + 128], filename[128], *result[64];
    int ns;

    sprintf(filename, PROC_CGPATH, pid);

    printf("filename: %s\n", filename);

    if ((fp = fopen(filename, "re")) == NULL)
        return -1;
    while (fgets(line, sizeof(line) - 1, fp) != NULL){
        printf("[line] %s", line);
        ns = split_string(line, result, 64, ':');
        if(ns != 3)
            continue;
        if(strstr(result[1], type) != NULL){
            strcpy_safe(path, result[2], size);
            goto found;
        }
    }
    fclose(fp);
    return -1;

found:
    printf("clear \\n at end\n");
    ns = strlen(path) - 1;
    while(path[ns] == '\n' && ns >= 0){
        path[ns--] = 0;
    }
    fclose(fp);
    return 0;
}


int
get_container_id_from_pid(char *id, unsigned int size, unsigned int pid)
{
    char buf[MAX_CGPATH_LEN];
    char *result[64];
    int n_s, r;

    r = get_cgroup_path_from_pid(buf, MAX_CGPATH_LEN, pid, "freezer");
    if (r < 0)
        return -1;
    if (buf[0] == '/' && buf[1] == 0){
        // root cpuacct
        strcpy_safe(id, "root", size);
        return 1;
    }
    n_s = split_string(buf, result, 64, '/');
    strcpy_safe(id, result[n_s - 1], size);

    return 0;
}


static void
read_cgroup_sched_stats(char *parameter, char *buf)
{
    char *wp;
    char path[256], line[1024];
    char cgroup_path[128];
    FILE *fp;
    int pid[MAX_NUM_OF_CONTAINERS];
    struct st_cg_sched st_cgsched;
    int n, i, r, kv;
    //char buf[CG_BUFSIZE];
    DEFINE_INCRON;

    //kv = kernel_version();
    kv = 409;

    //n = get_container_pids(pid, MAX_NUM_OF_CONTAINERS);
    n = 1;
    pid[0]=68206;

    wp = buf;
    for(i = 0; i < n; ++i) {
        r = get_cgroup_path_from_pid(cgroup_path, 128, pid[i], "cpuacct");
        if(r < 0) {
            continue;
        }
        printf("cgroup_path=%s pid=%d\n", cgroup_path, pid[i]);
        if(kv == 409) {
            snprintf(path, 256, "/sys/fs/cgroup/cpu%s/" CFS_STATICS, cgroup_path);
            printf("path: %s\n", path);
            sprintf(path, "/tmp/cpuacct.sched_cfs_statistics"); //for test
            printf("path: %s\n", path);
            if(access(path, R_OK) != 0) {
                continue;
            }
        }
        printf("ENTER_CGROUP_COLLECT cpu: cgroup_path=%s incron=%d\n", cgroup_path, incron );
        ENTER_CGROUP_COLLECT("cpu", cgroup_path, incron);
        memset(&st_cgsched, 0, sizeof(struct st_cg_sched));
        if(get_container_id_from_pid(st_cgsched.container_id, 64, pid[i]) < 0){
            ABORT_THIS_CGROUP_COLLECT(incron, "cpu");
        }

        printf("open path: %s\n", path);
        fp = fopen(path, "r");
        if(fp == NULL) {
            ABORT_THIS_CGROUP_COLLECT(incron, "cpu");
        }
        r = 0;
        if(kv == 409) {
            while(fgets(line, 1024, fp) != NULL) {
                r++;
                printf("[line %d] %s", r, line);
                if(r == 2) {
                    sscanf(line, "%lu %lu %lu %*u",
                                 &st_cgsched.bvt_delay,
                                 &st_cgsched.noise_kick_delay,
                                 &st_cgsched.loss_time);
                } else if(r == 4) {
                    sscanf(line, "%lu %lu %lu %lu %lu %lu %lu %lu %lu %lu %lu %lu",
                                 &st_cgsched.tasks_delay_ms[0],
                                 &st_cgsched.tasks_delay_ms[1],
                                 &st_cgsched.tasks_delay_ms[2],
                                 &st_cgsched.tasks_delay_ms[3],
                                 &st_cgsched.tasks_delay_ms[4],
                                 &st_cgsched.tasks_delay_ms[5],
                                 &st_cgsched.tasks_delay_ms[6],
                                 &st_cgsched.tasks_delay_ms[7],
                                 &st_cgsched.tasks_delay_ms[8],
                                 &st_cgsched.tasks_delay_ms[9],
                                 &st_cgsched.tasks_delay_ms[10],
                                 &st_cgsched.tasks_delay_ms[11]);
                    break;
                }
            }
        }
        fclose(fp);
        printf("EXIT_CGROUP_COLLECT cpu\n");
        EXIT_CGROUP_COLLECT("cpu");
        printf("write data\n");
        r = snprintf(wp, CG_BUFSIZE - (wp - buf), DATA_FORMAT,
                     st_cgsched.container_id,
                     st_cgsched.tasks_delay_ms[0],
                     st_cgsched.tasks_delay_ms[1],
                     st_cgsched.tasks_delay_ms[2],
                     st_cgsched.tasks_delay_ms[3],
                     st_cgsched.tasks_delay_ms[4],
                     st_cgsched.tasks_delay_ms[5],
                     st_cgsched.tasks_delay_ms[6],
                     st_cgsched.tasks_delay_ms[7],
                     st_cgsched.tasks_delay_ms[8],
                     st_cgsched.tasks_delay_ms[9],
                     st_cgsched.tasks_delay_ms[10],
                     st_cgsched.tasks_delay_ms[11],
                     st_cgsched.bvt_delay,
                     st_cgsched.noise_kick_delay,
                     st_cgsched.loss_time);
        printf("FINISH_CGROUP_COLLECT\n");
        FINISH_CGROUP_COLLECT(mod, buf, incron, r);
        if(r + wp >= buf + CG_BUFSIZE){
            break;
        }
        wp += r;
    }

}
//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

*/
import "C"

const (
	ITEM_SPLIT   = ";"
	ITEM_SPSTART = "="
	DATA_SPLIT   = ","
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
	item := strings.Split(content, ITEM_SPLIT)[0]
	cols := strings.Split(item, DATA_SPLIT)

	userCpu, _ = strconv.ParseUint(cols[0], 10, 64)
	sysCpu, _ = strconv.ParseUint(cols[1], 10, 64)
	mem, _ = strconv.ParseUint(cols[2], 10, 64)
	readBytes, _ = strconv.ParseUint(cols[3], 10, 64)
	writeBytes, _ = strconv.ParseUint(cols[4], 10, 64)
	fdCnt, _ = strconv.ParseUint(cols[5], 10, 64)
	minFlt, _ = strconv.ParseUint(cols[6], 10, 64)
	majFlt, _ = strconv.ParseUint(cols[7], 10, 64)
	thread, _ = strconv.ParseUint(cols[8], 10, 64)

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

	// call read_pid_stats
	dst := ""
	c_src := C.CString("1")
	defer C.free(unsafe.Pointer(c_src))
	c_dst := C.CString(dst)
	//defer C.free(unsafe.Pointer(c_dst))
	C.read_cgroup_sched_stats(c_src, c_dst)

	// show result
	received := C.GoString(c_dst)
	fmt.Printf("[go]received:   dst=%s\n", received)

	// convert result
	content := strings.Split(received, ITEM_SPSTART)[1]
	item := strings.Split(content, ITEM_SPLIT)[0]
	cols := strings.Split(item, DATA_SPLIT)
	for i := 0; i < 11; i++ {
		n, _ := strconv.ParseUint(cols[i], 10, 64)
		tasksDelayMs = append(tasksDelayMs, n)
	}
	bvtDelay, _ = strconv.ParseUint(cols[12], 10, 64)
	noiseKickDelay, _ = strconv.ParseUint(cols[13], 10, 64)
	lossTime, _ = strconv.ParseUint(cols[14], 10, 64)

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

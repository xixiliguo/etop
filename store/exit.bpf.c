//go:build ignore

#include "vmlinux.h"
#include "bpf_helpers.h"
#include "bpf_core_read.h"
#include "bpf_endian.h"
#include "bpf_tracing.h"

#define TASK_COMM_LEN 16
#define TASK_CMDLINE_LEN 32

char __license[] SEC("license") = "Dual MIT/GPL";

struct
{
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(int));
    __uint(value_size, sizeof(int));
} events SEC(".maps");

struct event
{
    int pid;
    int ppid;
    long exit_code;
    u8 comm[TASK_COMM_LEN];
    u8 cmdline[TASK_CMDLINE_LEN];
    u64 utime;
    u64 stime;
    u64 start_time;
    u64 end_time;
    int num_threads;
    u64 on_cpu;
    u64 priority;
    u64 nice;
    u64 delayacct_blkio_ticks;
    u64 min_flt;
    u64 maj_flt;
    u64 vss_pages;
    u64 rss_pages;
    u64 rchar;
    u64 wchar;
    u64 syscr;
    u64 syscw;
    u64 io_read_bytes;
    u64 io_write_bytes;
    u64 cancelled_write_bytes;
};

struct event *unused __attribute__((unused));

struct mm_rss_stat___pre62 {
	atomic_long_t count[4];
} __attribute__((preserve_access_index));

struct mm_struct___pre62 {
	struct mm_rss_stat___pre62 rss_stat;
} __attribute__((preserve_access_index));

struct mm_struct___post62 {
  struct percpu_counter rss_stat[NR_MM_COUNTERS];
} __attribute__((preserve_access_index));

// Same as __percpu_counter_read_positive in kernel
s64 percpu_counter_read_positive(struct percpu_counter *c) {
  s64 ret;
  ret = c->count;
  if (ret >= 0)
    return ret;
  return 0;
}

SEC("kprobe/tty_audit_exit")
int handle_exit(struct pt_regs *ctx)
{
    struct event e;
    __builtin_memset(&e, 0, sizeof(e));
    int pid;

    u64 id = bpf_get_current_pid_tgid();
    pid = id >> 32;

    u64 now = bpf_ktime_get_ns();
    struct task_struct *task = (struct task_struct *)bpf_get_current_task();

    e.pid = pid;
    e.ppid = BPF_CORE_READ(task, real_parent, tgid);
    
    e.exit_code = BPF_CORE_READ(task, signal, pacct.ac_exitcode);
    bpf_get_current_comm(&e.comm, sizeof(e.comm));

    e.utime = BPF_CORE_READ(task, signal, pacct.ac_utime) / 10000000;
    e.stime = BPF_CORE_READ(task, signal, pacct.ac_stime) / 10000000;

    e.start_time = BPF_CORE_READ(task, group_leader, start_time) / 10000000;
    e.end_time = now / 10000000;
    e.num_threads = BPF_CORE_READ(task, signal, nr_threads);

    e.on_cpu = bpf_get_smp_processor_id();
    e.priority = BPF_CORE_READ(task, prio) - 100;
    e.nice = BPF_CORE_READ(task, static_prio) - (100 + (19 - (-20) + 1) / 2);
    e.delayacct_blkio_ticks = BPF_CORE_READ(task, delays, blkio_delay) + BPF_CORE_READ(task, delays, swapin_delay);
    e.delayacct_blkio_ticks = e.delayacct_blkio_ticks / 10000000;

    e.min_flt = BPF_CORE_READ(task, signal, pacct.ac_minflt);
    e.maj_flt = BPF_CORE_READ(task, signal, pacct.ac_majflt);

    const struct mm_struct *mm = BPF_CORE_READ(task, mm);
    if (mm)
    {
        u64 arg_start = BPF_CORE_READ(mm, arg_start);
        u64 arg_end = BPF_CORE_READ(mm, arg_end);
        e.cmdline[0] = arg_end - arg_start;
        bpf_probe_read_user(&e.cmdline[1], TASK_CMDLINE_LEN-1, (void *)arg_start);

        e.vss_pages = BPF_CORE_READ(mm, total_vm);

        u64 file_pages = 0;
        u64 anon_pages = 0;
        u64 shmem_pages = 0;

        if (bpf_core_type_matches(struct mm_struct___pre62)) {
            const struct mm_struct___pre62 *mms = (struct mm_struct___pre62 *)mm;
            file_pages = BPF_CORE_READ(mms, rss_stat.count[MM_FILEPAGES].counter);
            anon_pages = BPF_CORE_READ(mms, rss_stat.count[MM_ANONPAGES].counter);
            shmem_pages = BPF_CORE_READ(mms, rss_stat.count[MM_SHMEMPAGES].counter);
        } else if (bpf_core_type_matches(struct mm_struct___post62)) {
            const struct mm_struct___post62 *mms = (struct mm_struct___post62 *)mm;
            struct percpu_counter file_percpu = BPF_CORE_READ(mms, rss_stat[MM_FILEPAGES]);
            struct percpu_counter anon_percpu = BPF_CORE_READ(mms, rss_stat[MM_ANONPAGES]);
            struct percpu_counter shmem_percpu = BPF_CORE_READ(mms, rss_stat[MM_SHMEMPAGES]);
            file_pages = percpu_counter_read_positive(&file_percpu);
            anon_pages = percpu_counter_read_positive(&anon_percpu);
            shmem_pages = percpu_counter_read_positive(&shmem_percpu);
        }
        e.rss_pages = file_pages + anon_pages + shmem_pages;
    }
    else
    {
        e.vss_pages = 0;
        e.rss_pages = 0;
    }

    e.rchar = BPF_CORE_READ(task, ioac.rchar);
    e.rchar += BPF_CORE_READ(task, signal, ioac.rchar);
    e.wchar = BPF_CORE_READ(task, ioac.wchar);
    e.wchar += BPF_CORE_READ(task, signal, ioac.wchar);
    e.syscr = BPF_CORE_READ(task, ioac.syscr);
    e.syscr += BPF_CORE_READ(task, signal, ioac.syscr);
    e.syscw = BPF_CORE_READ(task, ioac.syscw);
    e.syscw += BPF_CORE_READ(task, signal, ioac.syscw);
    e.io_read_bytes = BPF_CORE_READ(task, ioac.read_bytes);
    e.io_read_bytes += BPF_CORE_READ(task, signal, ioac.read_bytes);
    e.io_write_bytes = BPF_CORE_READ(task, ioac.write_bytes);
    e.io_write_bytes += BPF_CORE_READ(task, signal, ioac.write_bytes);
    e.cancelled_write_bytes = BPF_CORE_READ(task, ioac.cancelled_write_bytes);
    e.cancelled_write_bytes += BPF_CORE_READ(task, signal, ioac.cancelled_write_bytes);

    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &e, sizeof(e));

    return 0;
}
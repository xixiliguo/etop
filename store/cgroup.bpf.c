//go:build ignore

#include "vmlinux.h"
#include "bpf_helpers.h"
#include "bpf_core_read.h"
#include "bpf_endian.h"
#include "bpf_tracing.h"


char __license[] SEC("license") = "Dual MIT/GPL";


struct cgroup_net_stat {
	u64 rx_packet;
	u64 rx_byte;
	u64 tx_packet;
	u64 tx_byte;
};

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__type(key, u64);
	__type(value, struct cgroup_net_stat);
	__uint(max_entries, 1024);
} cgroup_net_stats SEC(".maps");


#define LOOPBACK  1

SEC("cgroup_skb/ingress")
int count_ingress_packets(struct __sk_buff *skb) {

	if (skb->ifindex == LOOPBACK ) {
		return 1;
	}

	struct cgroup_net_stat init_val = {
		.rx_packet = 1,
    	.rx_byte = skb->len,
    	.tx_packet = 0,
		.tx_byte = 0,
	};

	u64 cgroup_id = bpf_skb_cgroup_id(skb);

	struct cgroup_net_stat *stat = bpf_map_lookup_elem(&cgroup_net_stats, &cgroup_id);
	if (!stat) {
		bpf_map_update_elem(&cgroup_net_stats, &cgroup_id, &init_val, BPF_ANY);
		return 1;
	}
	__sync_fetch_and_add(&stat->rx_packet, 1);
	__sync_fetch_and_add(&stat->rx_byte, skb->len);

	return 1;
}


SEC("cgroup_skb/egress")
int count_egress_packets(struct __sk_buff *skb) {


	if (skb->ifindex == LOOPBACK ) {
		return 1;
	}

	struct cgroup_net_stat init_val = {
		.rx_packet = 0,
    	.rx_byte = 0,
    	.tx_packet = 1,
		.tx_byte = skb->len,
	};

	u64 cgroup_id = bpf_skb_cgroup_id(skb);

	struct cgroup_net_stat *stat = bpf_map_lookup_elem(&cgroup_net_stats, &cgroup_id);
	if (!stat) {
		bpf_map_update_elem(&cgroup_net_stats, &cgroup_id, &init_val, BPF_ANY);
		return 1;
	}
	__sync_fetch_and_add(&stat->tx_packet, 1);
	__sync_fetch_and_add(&stat->tx_byte, skb->len);

	return 1;
}

SEC("tracepoint/cgroup/cgroup_rmdir")
int cgroup_rmdir(struct trace_event_raw_cgroup *ctx) {
	u64 id = ctx->id;
	bpf_map_delete_elem(&cgroup_net_stats, &id);
	return 1;
}
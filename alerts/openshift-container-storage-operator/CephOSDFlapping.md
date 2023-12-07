# CephOSDFlapping

## Meaning

Ceph storage OSD is flapping, indicating that a daemon has restarted 5 times \
in the last 5 minutes.

## Impact

This may affect Ceph storage stability and reliability.

## Diagnosis

Check pod events or Ceph status to identify the cause of OSD flapping.

## Mitigating

### Recommended Network Configuration

The upstream Ceph community traditionally suggests having separate public
(front-end) and private (cluster/back-end/replication) networks, offering the
following benefits:

1. **Segregation of Traffic:**
   - Heartbeat traffic and replication/recovery traffic (private) are separated
     from traffic between clients and OSDs/monitors (public).
   - Prevents one stream of traffic from DoS-ing the other, avoiding cascading failures.

2. **Additional Throughput:**
   - Enhances throughput for both public and private traffic

### Halting Flapping

If OSDs repeatedly flap (marked down and then up again), force monitors to halt
the flapping by temporarily freezing their states:

```bash
ceph osd set noup      # prevent OSDs from getting marked up
ceph osd set nodown    # prevent OSDs from getting marked down
```

These flags are recorded in the osdmap:

```bash
ceph osd dump | grep flags
```

Two other flags, noin and noout, prevent booting OSDs from being marked in or
out, respectively. Clear these flags with:

```bash
ceph osd unset noup
ceph osd unset nodown
```

Two additional flags, noin and noout, prevent booting OSDs from being marked in
or protect OSDs from eventually being marked out, regardless of the current
value of mon_osd_down_out_interval.

> Note: noup, noout, and nodown are temporary; after clearing the flags, the
> blocked action becomes possible shortly thereafter. However, the noin flag
> prevents OSDs from being marked in on boot, and daemons that started while the
> flag was set will remain that way.
---
> Note: Causes and effects of flapping can be mitigated to some extent by making
> careful adjustments to mon_osd_down_out_subtree_limit,
> mon_osd_reporter_subtree_level, and mon_osd_min_down_reporters. The optimal
> settings depend on cluster size, topology, and the Ceph release in use. The
> interaction of all these factors is subtle and beyond the scope of this
> document.

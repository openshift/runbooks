# Troubleshooting Ceph

Ceph commands

Some common commands to troubleshoot a Ceph cluster:

* ceph status
* ceph osd status
* cepd osd df
* ceph osd utilization
* ceph osd pool stats
* ceph osd tree
* ceph pg stat

The first two status commands provide the overall cluster health.
The normal state for cluster operations is `HEALTH\_OK`, but will still function
 when the state is in a `HEALTH\_WARN` state. If you are in a `WARN` state, then
 the cluster is in a condition that it may enter the `HEALTH\_ERROR` state at
 which point all disk I/O operations are halted. If a `HEALTH\_WARN` state is
 observed, then one should take action to prevent the cluster from halting
 when it enters the `HEALTH\_ERROR` state.

Problem 1

Ceph status shows that the OSD is full .Example Ceph OSD-FULL error

```bash
  ceph status
    cluster:
      id:     62661e0d-417c-485e-b01f-562e9493f121
      health: HEALTH\_ERR
              3 full osd(s)
              3 pool(s) full

    services:
      mon: 3 daemons, quorum a,b,c (age 3h)
      mgr: a(active, since 3h)
      mds: ocs-storagecluster-cephfilesystem:1 {0=ocs-storagecluster-cephfilesystem-a=up:active} 1 up:standby-replay
      osd: 3 osds: 3 up (since 3h), 3 in (since 3h)

    data:
      pools:   3 pools, 192 pgs
      objects: 223.01k objects, 870 GiB
      usage:   2.6 TiB used, 460 GiB / 3 TiB avail
      pgs:     192 active+clean

    io:
      client:   853 B/s rd, 1 op/s rd, 0 op/s wr
```

1) Check the alert manager for readonly alert

```bash
  curl -k -H "Authorization: Bearer $(oc -n openshift-monitoring sa get-token prometheus-k8s)"  https://${MYALERTMANAGER}/api/v1/alerts | jq '.data[] | select( .labels.alertname) | { ALERT: .labels.alertname, STATE: .status.state}'
```

2) If CephClusterReadOnly alert is listed from the above curl command, then see:

    [CephClusterReadOnly alert](../CephClusterReadOnly.md)

Problem 2

Ceph status shows an issue with osd, as see in example below

Example Ceph OSD error

```bash
  cluster:
      id:     263935ae-deb3-47e0-9355-d4a5c935aaf5
      health: HEALTH\_ERR
              1 MDSs report slow metadata IOs
              2 osds down
              2 hosts (2 osds) down
              1 nearfull osd(s)
              3 pool(s) nearfull
              11/2142 objects unfound (0.514%)
              Reduced data availability: 237 pgs inactive, 237 pgs down
              Possible data damage: 8 pgs recovery\_unfound
              Degraded data redundancy: 833/6426 objects degraded (12.963%), 24 pgs degraded, 63 pgs undersized

    services:
      mon: 3 daemons, quorum a,b,c (age 115m)
      mgr: a(active, since 112m)
      mds: myfs:1 {0=myfs-b=up:active} 1 up:standby-replay
      osd: 3 osds: 1 up (since 2m), 3 in (since 113m)
```

Take a look to solving [common osd errors](https://access.redhat.com/documentation/en-us/red_hat_ceph_storage/4/html/troubleshooting_guide/troubleshooting-ceph-osds#most-common-ceph-osd-errors)

Problem 3

Issues seen with PG, example

Example Ceph PG error

```bash
  cluster:
      id: 0a1a6dcb-2146-42f7-9e6f-8b933614c45f
      health: HEALTH\_ERR Degraded data redundancy
              126/78009 objects degraded (0.162%)
              7 pgs degraded Degraded data redundancy (low space)
              1 pg backfill\_toofull

      data:
      pools:   10 pools, 80 pgs
      objects: 26.00k objects, 100 GiB
      usage:   306 GiB used, 5.7 TiB / 6.0 TiB avail
      pgs:     126/78009 objects degraded (0.162%)
              35510/78009 objects misplaced (45.520%)
              55 active+clean
              12 active+remapped+backfill\_wait
              4  active+recovery\_wait+undersized+degraded+remapped
              3  active+recovery\_wait+degraded
              2  active+recovery\_wait
              2  active+recovering+undersized+remapped
              1  active+recovering
              1  active+remapped+backfill\_toofull
```

Review [Solving pg error](https://access.redhat.com/documentation/en-us/red_hat_ceph_storage/4/html/troubleshooting_guide/troubleshooting-ceph-placement-groups#most-common-ceph-placement-group-errors)



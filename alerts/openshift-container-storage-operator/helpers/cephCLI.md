# How to execute ceph commands from Openshift

Locate your rook-ceph-operator pod and connect into it

```bash
  oc rsh -n openshift-storage $(oc get pods -n openshift-storage -o name -l app=rook-ceph-operator)
```

Set your CEPH_ARGS environment variable

```bash
  sh-4.4$ export CEPH_ARGS='-c /var/lib/rook/openshift-storage/openshift-storage.config'
```

One can now run Ceph commands

```bash
  sh-5.1$ ceph -s
    cluster:
      id:     050dce42-7e96-4e9f-abff-74a14891376a
      health: HEALTH_OK

    services:
      mon: 3 daemons, quorum a,b,c (age 28m)
      mgr: a(active, since 27m)
      mds: 1/1 daemons up, 1 hot standby
      osd: 3 osds: 3 up (since 27m), 3 in (since 27m)

    data:
      volumes: 1/1 healthy
      pools:   4 pools, 113 pgs
      objects: 91 objects, 129 MiB
      usage:   285 MiB used, 6.0 TiB / 6 TiB avail
      pgs:     113 active+clean

    io:
      client:   853 B/s rd, 18 KiB/s wr, 1 op/s rd, 1 op/s wr
```

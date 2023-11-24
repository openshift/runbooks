# Check Ceph Network Connectivity
## Check OCP Cluster Connectivity
### Option 1 (Preferred)

First, get the IP addresses and ports of the Ceph mons:

```bash
    oc get cm rook-ceph-mon-endpoints -ojson -n openshift-storage | jq -r '.data."data"'
```

Now, start a debug pod on any node in the cluster

```bash
    oc debug node/<any-node-name> -n openshift-storage
```

Node names can be found with `oc get nodes`.
Once you are in the debug container, start a `toolbox-root` container with the
RHEL support tools:

```bash
    sh-4.4# chroot /host
    sh-4.4# toolbox
```

Now you are able to use Netcat and can check the connectivity to the mons using:

```bash
    nc -vz <mon-ip> <mon-port>
```

If you cannot connect to the mons, that is an indication that there is a
network issue.

### Option 2
The first option is preferred, only use this option if you do not have
permission to create a debug pod.

First, count the number of Ceph mon pods:

```bash
    oc get pods -n openshift-storage | grep rook-ceph-mon
```

Next, get the ceph tools pod:

```bash
    oc get pods -n openshift-storage | grep rook-ceph-tools
```

Then, create a remote shell session in that pod:

```bash
    oc rsh -n openshift-storage <rook-ceph-tools pod name from previous step>
```

Finally, check the status of the Ceph cluster:

```bash
    ceph -s
```

It is an indication that there might be a network connectivity issue if you get
an output which has a cluster status of `HEALTH_WARN` and which does not have
all the mon daemons in the service section, e.g.

```bash
    $ oc rsh -n openshift-storage rook-ceph-tools-79ccc8ddc5-77brq
    sh-4.4$ ceph -s
        cluster:
        id:     70f6e8bf-3ee1-494e-b6a0-c89ac75582d1
        health: HEALTH_WARN
                1 MDSs report slow metadata IOs

        services:
        mon: 1 daemons, quorum a (age 116m)
        mgr: no daemons active
        mds: 1/1 daemons up
        osd: 0 osds: 0 up, 0 in

        data:
        volumes: 1/1 healthy
        pools:   0 pools, 0 pgs
        objects: 0 objects, 0 B
        usage:   0 B used, 0 B / 0 B avail
        pgs:
```

## Check Consumer to Provider Connectivity
To start, get the provider IP and port:

```bash
    oc get storagecluster -n openshift-storage -o json | grep storageProviderEndpoint
```

Now, create a debug pod in the `openshift-storage` namespace:

```bash
    oc debug -n openshift-storage
```

First, check that you can ping the provider:

```bash
    oc exec bash -- ping -c 5 <provider ip from first step>
```

Second, check that you can connect to the port:

```bash
    nc -vz <provider ip from first step> <provider port from first step>
```

## If it is a Network Issue
We cannot do much about the possible causes of network issues e.g. misconfigured
 AWS security group. Escalate to the ODF/openshift support.

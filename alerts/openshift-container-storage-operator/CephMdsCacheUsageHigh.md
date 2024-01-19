# MDSCacheUsageHigh

## Meaning

Ceph MDS cache usage for the MDS daemon has exceeded above 95% of the `mds_cache_memory_limit`.

## Impact

If the MDS cannot keep its cache usage under the target threshold, that is,
`mds_health_cache_threshold` (150%) of the cache limit, that is,
`mds_cache_memory_limit`, the MDS will send a health alert to the Monitors
indicating the cache is too large.

As a result the MDS related operations, like, caps revocation, will become slow.

## Diagnosis

Check the usage of `ceph_mds_mem_rss` metric and ensure that it is under 95% of the
cache limit set in `mds_cache_memory_limit`.

The MDS tries to stay under a reservation of the `mds_cache_memory_limit` by
trimming unused metadata in its cache and recalling cached items in the client
caches. It is possible for the MDS to exceed this limit due to slow recall from
clients as result of multiple clients accesing the files.

Read more about ceph MDS cache configuration [here](https://docs.ceph.com/en/latest/cephfs/cache-configuration/?highlight=mds%20cache%20configuration#mds-cache-configuration)

## Mitigation

Make sure we have enough memory provisioned for MDS cache. Default is 4GB, but recomended
is minimum 8GB.

Memory resources for the MDS pods should be updated in the ocs-storageCluster in
order to increase the `mds_cache_memory_limit`. For example, run the following command
to set the memory of MDS pods to 8GB

```bash
oc patch -n openshift-storage storagecluster ocs-storagecluster \
    --type merge \
    --patch '{"spec": {"resources": {"mds": {"limits": {"memory": "8Gi"},"requests": {"memory": "8Gi"}}}}}' 
```

**Note**: ODF sets `mds_cache_memory_limit` to half of the MDS pod memory request/limit.
So if the memory is set to 8GB using above command, then the operator will set the
mds cache memory limit to 4GB

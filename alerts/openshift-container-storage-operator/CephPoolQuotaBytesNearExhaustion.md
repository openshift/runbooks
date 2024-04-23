# CephPoolQuotaBytesNearExhaustion

## Meaning

Storage pool quota usage has crossed 70%. One or more pools is approaching a
configured fullness threshold.
One threshold that can trigger this warning condition is the
`mon_pool_quota_warn_threshold` configuration option.

## Impact

Due the quota configured the pool will become readonly when the quota will be
exhausted completelly

## Diagnosis

[Execute the following Ceph command](helpers/cephCLI.md) to have information
about pool status in the cluster:

```bash
    sh-5.1$ rados df
```

## Mitigation

Pool quotas can be adjusted up or down (or removed) with
[Ceph CLI](helpers/cephCLI.md)

```bash
    ceph osd pool set-quota <pool> max_bytes <bytes>
    ceph osd pool set-quota <pool> max_objects <objects>
```

Setting the quota value to 0 will disable the quota.

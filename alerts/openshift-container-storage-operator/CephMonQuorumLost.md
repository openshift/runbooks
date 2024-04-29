# CephMonQuorumLost

## Meaning

The number of monitors in the storage cluster are not enough.
Multiple mons work together to provide redundancy by each keeping a copy
of the metadata. Cluster is deployed with 3 or 5 mons, and require 2 or more mons
to be up and running for quorum and for the storage operations to run.

This alert indicates that there is only 1 monitor pod running or even none.

## Impact

If quorum is lost and it is beyond recovery now. Any data lose is permanent
at this point.

## Diagnosis

Set logging to files true,
`# ceph config set global log_to_file true`
`# ceph config set global mon_cluster_log_to_file true`
Then check the corresponding Ceph Monitor logs in /var/log/ceph/<cluster-id> location

## Mitigation

[Restore Ceph Mon Quorum Lost](https://access.redhat.com/solutions/5898541)
[Troubleshooting Monitor](https://docs.ceph.com/en/latest/rados/troubleshooting/troubleshooting-mon/)



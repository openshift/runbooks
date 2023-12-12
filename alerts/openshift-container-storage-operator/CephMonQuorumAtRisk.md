# CephMonQuorumAtRisk

## Meaning

Storage cluster quorum is low.
Multiple mons work together to provide redundancy by each keeping a copy
of the metadata. Cluster is deployed with 3 mons, and require 2 or more mons
to be up and running for quorum and for the storage operations to run.

## Impact

If quorum is lost, access to data is at risk.

## Diagnosis

Run following command for each Monitor in the cluster.
`ceph tell mon.ID mon_status`

For more on this command’s output, see [Understanding mon_status](https://docs.ceph.com/en/latest/rados/troubleshooting/troubleshooting-mon/#rados-troubleshoting-troubleshooting-mon-understanding-mon-status).

## Mitigation

[Restore Ceph Mon Quorum Lost](https://access.redhat.com/solutions/5898541)
[Troubleshooting Monitor](https://docs.ceph.com/en/latest/rados/troubleshooting/troubleshooting-mon/)


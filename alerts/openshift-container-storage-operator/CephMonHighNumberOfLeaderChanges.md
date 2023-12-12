# CephMonHighNumberOfLeaderChanges

## Meaning

In a Ceph cluster there is a redundant set of monitors that store critical
information about the storage cluster. Monitors synchronize periodically to
obtain information about the storage cluster. The first monitor to get the
most updated information become leader and other monitors will start their
synchronization process asking the leader.

This alert indicates a high frequent Ceph Monitor leader change per minute.

## Impact

An unusual change of leader is usually produced by problems in network
connection, or another kind of problem in one or more monitor pods.
This situation can affect negatively to the storage cluster performance.

## Diagnosis

The alert should indicate the monitor pod with the problem:

    Ceph Monitor <rook-ceph-mon-X... pod> on host <hostX> has seen <X> leader
    changes per minute recently.

Check the affected monitor's logs. More information on the cause can be seen
from these logs.

## Mitigation

[pod debug](helpers/podDebug.md) [gather_logs](helpers/gatherLogs.md)


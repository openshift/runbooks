# CephOSDDown

## Meaning

CephOSDDown indicates that one or more ceph-osd daemons are not running as expected.

## Impact

Ceph detects that the OSDs are down and automatically starts the recovery
process by moving the data to other available OSDs. But if the OSDs having
the copies of the data also fail during this recovery, then there is a
chance of permanent data loss.

## Diagnosis

The alert is triggered when Ceph OSD(s) is/are down,
please check the ceph-osd daemons and take corrective measures.

## Mitigation

### Recommended Actions

1. In an LSO cluster, if the disk failed, the OSD may need to be replaced.
Please ref:
[Instructions for safely replacing operational or failed devices]:
  https://docs.redhat.com/en/documentation/red_hat_openshift_data_foundation/4.17/html-single/replacing_devices/index

2. Investigate why one or more OSDs are marked down. The ceph-osd daemon(s)
or their host(s) may have crashed or been stopped, or peer OSDs might be
unable to reach the OSD over the public or private network. Common causes
include a stopped or crashed daemon, a “down” host, or a network failure.

Verify that the host is healthy, if not switch to the node/host where the
failed OSDs are running and check the logs at
(/var/lib/rook/openshift-storage/*) may contain troubleshooting information.

Unable to resolve the failed OSDs, please connect with Red Hat Support.

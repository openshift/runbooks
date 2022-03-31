# etcdMembersDown

## Meaning

This alert fires etcd members report divergent revisions.

## Impact

Members may be stuck or, in extreme circumstances, split brained.

Stuck members can cause quorum loss, or bring the etcd cluster closer
to quorum loss, which breaks writes and [probably breaks
reads][default-quorum-read].

Split-brained members may allow writes to continue, but etcd responses
(and thus Kubernetes API responses, and other responses which depend
on etcd state) may no longer be consistent.

## Diagnosis

If the alert goes off, you're in trouble.

## Mitigation

For stuck members, you may need to [replace the lagging
member][replace-member].

For split-brained members, see [rhbz#2068601][rbhz-2068601] and [the
associated KCS Solution][kcs-split-brained].

[docs]: https://docs.openshift.com/container-platform/4.10/backup_and_restore/disaster_recovery/about-disaster-recovery.html
[default-quorum-read]: https://github.com/kubernetes/kubernetes/pull/53717
[kcs-split-brained]: https://access.redhat.com/solutions/6849521
[rbhz-2068601]:https://bugzilla.redhat.com/show_bug.cgi?id=2068601
[replace-member]: https://docs.openshift.com/container-platform/4.10/backup_and_restore/control_plane_backup_and_restore/replacing-unhealthy-etcd-member.html

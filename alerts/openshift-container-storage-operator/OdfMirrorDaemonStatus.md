# OdfMirrorDaemonStatus

## Meaning

Mirror daemon is in unhealthy status for more than 1 minute. Mirroring on this
cluster is not working as expected. Disaster recovery is failing for the entire
cluster.

## Impact

Critical.

The mirroring operations are stopped. No images will be synced until RBD mirror
daemon will be started and running properly.
Disaster recovery functionality will not work properly until RBD mirror daemons
will be running an all images synced properly.

## Diagnosis

The RBD mirror pod in the Openshift storage namespace is in error state.

## Mitigation

Review [RBD mirrorpod](helpers/podDebug.md) to find out more information about
the problem.

[gather_logs](helpers/gather_logs.md) to provide more information to support
teams.

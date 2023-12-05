# OdfPoolMirroringImageHealth

## Meaning

Mirroring image(s) (PV) a pool are in Unknown state for more than 1 minute.
Mirroring might not work as expected.

## Impact

Critical.
Image duplication to the remote cluster is not working. Disaster recovery
cluster will be affected because in case of failover the information will not
be synced with the source cluster properly

## Diagnosis

Verify images status using the [Ceph CLI](helpers/cephCLI.md):

```bash
    rbd mirror image status POOL_NAME/IMAGE_NAME
```

## Mitigation

Review [RBD mirrorpod](helpers/podDebug.md) to find out more information about
the problem.

[gather_logs](helpers/gather_logs.md) to provide more information to support
teams.

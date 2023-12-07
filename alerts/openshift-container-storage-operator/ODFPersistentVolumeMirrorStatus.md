# ODFPersistentVolumeMirrorStatus

## Meaning

The alert 'ODFPersistentVolumeMirrorStatus' indicates the mirroring status of
persistent volumes (PVs) in a Ceph pool. The two specific alert instances are
defined as follows:

1. **Critical Alert:**
   - The mirroring image(s) (PV) in a pool are not mirrored properly to the peer
     site for more than the specified alert time.
   - RBD image, CephBlockPool, and affected persistent volume details are
     provided in the alert message.

2. **Warning Alert:**
   - The status is unknown for the mirroring of persistent volumes (PVs) to the
     peer site for more than the specified alert time.
   - RBD image, CephBlockPool, and affected persistent volume details are
     included in the alert message.

## Impact

- **Critical Alert:**
  - Critical severity indicates that the mirroring of the PVs is not working
    correctly.
  - Disaster recovery may be affected, leading to potential data inconsistencies
    during failover.

- **Warning Alert:**
  - Warning severity implies an unknown status in mirroring, posing a potential
    risk to data consistency during failover.

## Diagnosis

Verify the mirroring status of the affected PVs using the
[Ceph CLI](helpers/cephCLI.md):

```bash
rbd mirror image status POOL_NAME/IMAGE_NAME
```

## Mitigation

Review [RBD mirrorpod](helpers/podDebug.md) to find out more information about
the problem.

[gather_logs](helpers/gatherLogs.md) to provide more information to support
teams.

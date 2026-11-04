# OSDKeyRotationJobFailed

## Meaning

Rook runs a CronJob `rook-ceph-osd-key-rotation-<osd-id>` for each encrypted
OSD. Each run (default `@weekly`) rotates the LUKS encryption key of the OSD's
device and stores the new key in the KMS (the Kubernetes Secret
`rook-ceph-osd-encryption-key-<pvc-name>` or an external KMS such as Vault).
Key rotation is enabled by default when encryption is used without an external
KMS and can be set with `spec.encryption.keyRotation` on the `StorageCluster`.

This alert fires when a key rotation Job has failed, i.e. it used up all its
retries without rotating the key.

## Impact

**Severity:** Critical
**Potential Customer Impact:** High

The OSD's encryption key was not rotated. If the key stored in the KMS no
longer unlocks the device, the running OSD keeps working, but it will fail to
start after its pod or node restarts.

## Diagnosis

See [general diagnosis document](helpers/diagnosis.md)

1. Find the failed Job (the alert's `job_name` label) and check its pod logs.
   The OSD ID is the number at the end of the CronJob name.

    ```bash
    oc get jobs -n openshift-storage | grep rook-ceph-osd-key-rotation
    oc get pods -n openshift-storage -l job-name=<job_name>
    oc logs <key_rotation_pod> -n openshift-storage --previous
    ```

    If the pod is gone, re-run the Job (Mitigation step 2) and follow its logs.

2. Match the log error to its cause:

    - `failed to get secret`, `failed to update secret`, or
      `failed to verify the new key in the KMS`: KMS problem, such as an
      unreachable KMS, an expired Vault token, wrong permissions, or a missing
      Secret.
    - `failed to add the current key to slot "1"`: the key in the KMS no longer
      unlocks the device.
    - Any other `slot` error: cryptsetup failed on the device (device busy,
      timeout).

3. If the pod never ran (`Pending`, `ImagePullBackOff`), follow
   [Pod debugging](helpers/podDebug.md). The pod must run on the same node as
   its OSD, so check that the node is `Ready` and has enough resources.

4. For KMS errors:
    - Without an external KMS, check that the OSD's Secret
      `rook-ceph-osd-encryption-key-<pvc-name>` exists and is not empty.
    - With an external KMS, check the `ocs-kms-connection-details` ConfigMap
      and whether [KMSServerConnectionAlert](KMSServerConnectionAlert.md) is
      also firing.

## Mitigation

1. Fix the cause found during diagnosis.

2. Re-run the rotation and confirm the log ends with
   `Successfully rotated the key`. Re-running is safe.

    ```bash
    oc create job --from=cronjob/rook-ceph-osd-key-rotation-<osd-id> \
      rook-ceph-osd-key-rotation-<osd-id>-manual -n openshift-storage
    oc logs -f -n openshift-storage -l job-name=rook-ceph-osd-key-rotation-<osd-id>-manual
    ```

3. Delete the failed Job. The alert keeps firing while it exists, even after a
   later run succeeds.

    ```bash
    oc delete job <job_name> -n openshift-storage
    ```

4. If the rotation still fails, [gather logs](helpers/gatherLogs.md) and
   contact Red Hat Support. Do not delete the failed Job or disable key
   rotation, since that only clears the alert. If the error is
   `failed to add the current key to slot "1"`, **do not restart the OSD pod,
   or drain or reboot its node**, because the OSD may not come back up.

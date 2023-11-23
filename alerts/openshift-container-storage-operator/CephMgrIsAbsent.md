# CephMgrIsAbsent

## Meaning

Ceph Manager cannot be found or accessed from Prometheus target discovery.

## Impact

Storage metrics collector won't be available and users won't be able to see
the storage metrics anymore. Not having a Ceph manager running impacts the
monitoring of the cluster, PVC creation and deletion requests and should be
resolved as soon as possible. Impact is critical here.

## Diagnosis

### Check if it is a network issue
Check the [network connectivity](helpers/networkConnectivity.md).
We cannot do much about the possible causes of network issues
e.g. misconfigured AWS security group.
Therefore, if it is a network issue, escalate to the ODF team by following
the steps [here](helpers/sre-to-engineering-escalation.md#procedure).

## Mitigation

Verify the rook-ceph-mgr pod is failing and restart if necessary.
If the ceph mgr pod restart fails, use general basic pod troubleshooting to resolve.

Verify the ceph mgr pod is failing:

    oc get pods -n openshift-storage | grep mgr

Describe the ceph mgr pod for more detail:

    oc describe -n openshift-storage pods/<rook-ceph-mgr pod name from previous step>

Analyze errors (i.e. resource issues?)

Try deleting the pod and watch for a successful restart:

    oc get pods -n openshift-storage | grep mgr

If above fails, follow general pod troubleshooting procedures.

[pod debug](helpers/podDebug.md) [gather_logs](helpers/gatherLogs.md)


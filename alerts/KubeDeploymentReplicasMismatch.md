# KubeDeploymentReplicasMismatch

## Meaning

This alert is fired in case of a discrepancy between the desired number of replicas
to the actual number of running instances was observed for deployment over a
certain period.

## Impact

This strength of the impact very much differs depending on the discrepancy.

## Diagnosis

The alert should note where the discrepancy occurred under the `deployment`
label.

### elasticsearch-*

The notification details should list the deployment with the mismatch (and the
deployment's namespace). For example:

```console
 - alertname = KubeDeploymentReplicasMismatch
...
 - deployment = elasticsearch-cdm-u1gqqbu6-2
...
 - namespace = openshift-logging
...
```

Start by checking the status of the deployment:

```console
$ oc get deploy -n $NAMESPACE $DEPLOYMENT
```

Possibilities include (but are not limited to) a pod stuck in
`ContainerCreating` or `CrashLoopBackoff`.

### etcd-quorum-guard-*"

Review the current deployment using the details available in the alert.
Review the following in the target namespace to ascertain the reason behind this:
* events
* pods
* nodes (check scheduling status)

## Mitigation

The specific instance that was resolved using the above required deleting all 3
of the `etcd-quorum-guard` pods.


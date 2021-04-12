# KubeDeploymentReplicasMismatch

## Meaning

This alert is fired when a discrepancy between the desired number of replicas to
the actual number of running instances for deployment was observed for a certain
period.

## Impact

The impact very much differs depending on the discrepancy.

## Diagnosis

The alert should note where the discrepancy occurred under the `deployment`
label:

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

Review the current deployment using the details available in the alert. Review
the following in the target namespace to ascertain the reason behind this. The
events:

```console
$ oc get events -n $NAMESPACE
```

Further, check the states of the Pods that the deployment manages:

```console
$ oc get pods -n $NAMESPACE --selector=app=$DEPLOYMENT
```

Possibilities include (but are not limited to) a pod stuck in
`ContainerCreating` or `CrashLoopBackoff`. The events may list this case
information about possible failed actions of a pod. Application and startup
failures should be visible with:

```console
$ oc describe pod $POD
```

If Pods are stuck in `Pending`, it means that insufficient resources prevent the
pod from being scheduled. Check the health of the nodes.

```console
$ oc get nodes
```

It is further possible that the CPU and Memory of the host are exhausted.

```console
$ oc adm top nodes
```

## Mitigation

Resolve the problems discovered during the diagnosis according to the
documentation. It is safe to delete the pods since they are managed by the
deployment. However, it may also be required to add more nodes in case of
insufficient resources.

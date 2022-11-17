# KubeDeploymentReplicasMismatch

## Meaning

Th `KubeDeploymentReplicasMismatch` alert triggers when, over a certain time
period, a discrepancy occurs between the desired number of pod replicas for
deployment and the actual number of running instances.

## Impact

The impact differs depending on the size of the discrepancy.

## Diagnosis

The alert message under the `deployment` label describes where the discrepancy
occurred:

```console
 - alertname = KubeDeploymentReplicasMismatch
...
 - deployment = elasticsearch-cdm-u1gqqbu6-2
...
 - namespace = openshift-logging
...
```

Review the current deployment details by examining the items available in
the alert.

* Start by reviewing the status of the deployment:

    ```console
    $ oc get deploy -n $NAMESPACE $DEPLOYMENT
    ```

* Run the following command in the target namespace to review the
events:

    ```console
    $ oc get events -n $NAMESPACE
    ```

* Review the status of the pods that the deployment manages:

    ```console
    $ oc get pods -n $NAMESPACE --selector=app=$DEPLOYMENT
    ```

    Possible problems include a pod stuck in a `ContainerCreating` or
    `CrashLoopBackoff` state.

* The events might also list information about possible failed actions of a
pod. You can view application and start-up failures by running:

    ```console
    $ oc describe pod $POD
    ```

* For pods stuck in a `Pending` state, insufficient resources are
preventing the pod from being scheduled. Check the health of the nodes:

    ```console
    $ oc get nodes
    ```

* Verify whether or not the host has sufficient CPU and memory resources:

    ```console
    $ oc adm top nodes
    ```

## Mitigation

After you diagnose the issue, refer to the OpenShift documentation to learn how
to resolve the problems. You can safely delete the pods because they are
managed by the deployment. However, you might also need to add more nodes if
your diagnostic steps showed that the host had insufficient resources.

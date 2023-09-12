# KubeDaemonSetMisScheduled

## Meaning

The `KubeDaemonSetMisScheduled` alert triggers when, over a defined time interval,
an excessive number of DaemonSet pods are discovered to be running on the nodes
where they should not be running. This condition typically occurs when a pod is
scheduled on a node with a taint for which the pod does not have a corresponding
toleration.

## Impact

DaemonSet pods are running on the nodes where they are not supposed to run, which
results in excessive resource.

## Diagnosis

- Check DaemonSet status:

  ```console
  oc -n $NAMESPACE describe daemonset $NAME
  ```

  Replace `$NAME` with the name of the DaemonSet having the issue.
  The namespace name is included along with alert labels, so replace
  `$NAMESPACE` with the namespace name from the received alert.

- Check the status of the pods belonging to the DaemonSet and examine the `Node`,
  `Tolerations`, and `NodeSelector` fields:

  ```console
  oc -n $NAMESPACE describe pod $NAME
  ```

  Replace $NAME with the name of the pod having the issue.

  - Look for the `Node`, `Node-Selectors` and `Tolerations` fields in the output.

  - Check the labels for the node and verify that they match
    the `NodeSelector` value for the pod.

    ```console
    oc get node $NAME --show-labels
    ```

    Replace `$NAME` with the node name of the pod having the issue.

  - Verify that the `Node-Selectors` field for the pod has one of the
    same labels as the node

    ```console
    oc -n $NAMESPACE describe pod $NAME | grep Node-Selectors
    ```

  - Check the `Taints` value for the node and verify that it matches the
    `Tolerations` value for the pod:

    - Check the taints for the node:

      ```console
      oc describe node $NAME
      ```

      Replace $NAME with the name of the node identified in the previous step.

    - Check the tolerations for the pod:

      ```console
      oc-n $NAMESPACE describe pod $NAME | grep Tolerations
      ```

     Replace `$NAMESPACE` with the namespace name and `$NAME` with the name of
     the pod having the issue.

    - Verify that the identified `Taints` are tolerated, that is, that they are
      identified in the `Tolerations` field in the previous step.
  
## Mitigation

- Verify node health [1]
  - If nodes are unhealthy or unsuitable for running pods, consider draining or
    evicting pods from those nodes to allow them to recover or be replaced.

- Check pod template parameters such as the following:
  - `pod priority`:  A setting for this parameter can cause the pod to be
     evicted by higher priority pods.
  - `affinity rules`: A setting for this parameter might restrict the number of
     nodes, which will make it impossible for other pods to be scheduled.

- Check taints and tolerations:

  - Verify that the nodes where mis-scheduled pods are running have the
    expected taints.
  - Confirm that the pods have tolerations matching those taints.

- Inspect the Kubernetes event logs for any relevant events that might explain
  why pods were mis-scheduled.


[1]: https://docs.openshift.com/container-platform/latest/support/troubleshooting/verifying-node-health.html

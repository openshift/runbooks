# StorageAutoScalingFailed

## Meaning

During automatic storage scaling, failures may occur due to various factors. In
such instances, the `StorageAutoScalingFailed` alert is fired.

## Impact

The storage cannot be scaled at the moment.

## Diagnosis

There are multiple stages where the storage configurations might not have matched
the requested storage parameters.

1. OSD pod is in `Pending` state:
    If it is in `Pending` state with the following logs, there is a high chance
    it is due to the exhaustion of the resources. Describe the pod in `Pending`
    state and check for the events.
    ```bash
    oc describe pod <osd_pod> -n <namespace>
    ```

    If the output is similar to the example below, it indicates that resources
    are exhausted:
    ```bash
    Warning  FailedScheduling  3m53s  default-scheduler  0/6 nodes are available: 1 Insufficient cpu, 2 node(s) didn't match Pod's node affinity/selector, 3 node(s) had untolerated taint {node-role.kubernetes.io/master: }. preemption: 0/6 nodes are available: 1 Insufficient cpu, 5 Preemption is not helpful for scheduling.
    ```

    If there are no errors, retrieve the pvc of the pod in `Pending` state.
    ```bash
    pvcName=$(oc get pod <osd_pod> -n <namespace> -o jsonpath='{.spec.volumes[*].persistentVolumeClaim.claimName}')
    oc get pvc $pvcName -n <namespace> -o yaml
    ```
    Check for the health and events of the PVC for detailed information.

2. PVC is in `Pending` state:

    ```bash
    oc get pvc <pvc_name> -n <namespace> -o yaml
    ```

    Inspect the events of the pvc. Additionally, retrieve the pv of the pvc if
    it is still in `pending` state.
    ```bash
    oc get pvc <pvc_name> -n <namespace> -o jsonpath='{.spec.volumeName}'
    ```

3. PV is in `Pending` state:
    ```bash
    oc get pv <pv_name> -n <namespace> -o yaml
    ```
    If the pv is in `Pending` state, there might be a storage issue. Check the
    events of the pv for more information.

## Mitigation

1. If resources are depleted and pods start failing, allocate additional
resources to the storage cluster. Possible solutions include adding more nodes
or scaling up CPU and memory on existing nodes.
- For instructions on adding nodes, please refer to this document: [Adding a node](https://docs.redhat.com/en/documentation/red_hat_openshift_data_foundation/latest/html-single/scaling_storage/index#adding_a_node)

2. Check upon the Ceph Storage: In case the OSDs are `down` or `out`, restart them
and try again.

3. During vertical scaling, if the pv takes time to resize, then the osds would not
resize correctly. In this case, the osd pods can be restarted.

4. Additionally, one could try to restart the rook operator if it is stuck reconciling:
    ```bash
    oc rollout restart deployment rook-ceph-operator -n <namespace>
    ```

    After restarting, look for errors or warnings in the logs of rook-ceph-operator:
    ```bash
    oc logs -l app=rook-ceph-operator -n <namespace>
    ```
    Monitor the status of the pods and cephcluster post restart.
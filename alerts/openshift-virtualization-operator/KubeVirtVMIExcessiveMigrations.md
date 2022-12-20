# KubeVirtVMIExcessiveMigrations
<!-- Edited by apinnick, Nov 2022-->

## Meaning

This alert fires when a virtual machine instance (VMI) live migrates more than
12 times over a period of 24 hours.

This migration rate is abnormally high, even during an upgrade. This alert might
indicate a problem in the cluster infrastructure, such as network disruptions
or insufficient resources.

## Impact

A virtual machine (VM) that migrates too frequently might experience degraded
performance because memory page faults occur during the transition.

## Diagnosis

1. Verify that the worker node has sufficient resources:

   ```bash
   $ oc get nodes -l node-role.kubernetes.io/worker= -o json | \
     jq .items[].status.allocatable
   ```

   Example output:

   ```json
   {
     "cpu": "3500m",
     "devices.kubevirt.io/kvm": "1k",
     "devices.kubevirt.io/sev": "0",
     "devices.kubevirt.io/tun": "1k",
     "devices.kubevirt.io/vhost-net": "1k",
     "ephemeral-storage": "38161122446",
     "hugepages-1Gi": "0",
     "hugepages-2Mi": "0",
     "memory": "7000128Ki",
     "pods": "250"
   }
   ```

2. Check the status of the worker node:

   ```bash
   $ oc get nodes -l node-role.kubernetes.io/worker= -o json | \
     jq .items[].status.conditions
   ```

   Example output:

   ```json
   {
     "lastHeartbeatTime": "2022-05-26T07:36:01Z",
     "lastTransitionTime": "2022-05-23T08:12:02Z",
     "message": "kubelet has sufficient memory available",
     "reason": "KubeletHasSufficientMemory",
     "status": "False",
     "type": "MemoryPressure"
   },
   {
     "lastHeartbeatTime": "2022-05-26T07:36:01Z",
     "lastTransitionTime": "2022-05-23T08:12:02Z",
     "message": "kubelet has no disk pressure",
     "reason": "KubeletHasNoDiskPressure",
     "status": "False",
     "type": "DiskPressure"
   },
   {
     "lastHeartbeatTime": "2022-05-26T07:36:01Z",
     "lastTransitionTime": "2022-05-23T08:12:02Z",
     "message": "kubelet has sufficient PID available",
     "reason": "KubeletHasSufficientPID",
     "status": "False",
     "type": "PIDPressure"
   },
   {
     "lastHeartbeatTime": "2022-05-26T07:36:01Z",
     "lastTransitionTime": "2022-05-23T08:24:15Z",
     "message": "kubelet is posting ready status",
     "reason": "KubeletReady",
     "status": "True",
     "type": "Ready"
   }
   ```

3. Log in to the worker node and verify that the `kubelet` service is running:

   ```bash
   $ systemctl status kubelet
   ```

4. Check the `kubelet` journal log for error messages:

   ```bash
   $ journalctl -r -u kubelet
   ```

## Mitigation

Ensure that the worker nodes have sufficient resources (CPU, memory, disk) to
run VM workloads without interruption.

If the problem persists, try to identify the root cause and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.

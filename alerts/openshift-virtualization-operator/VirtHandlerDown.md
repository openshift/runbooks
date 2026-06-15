# VirtHandlerDown

## Meaning

No running `virt-handler` pod has been detected for 10 minutes.

The alert expression evaluates
`cluster:kubevirt_virt_handler_pods_running:count == 0` with a `for`
duration of 10 minutes. The recording rule counts pods in `Running`
phase matching `virt-handler-.*`. Because the recording rule
aggregates across all pods, the alert carries only static labels
(`severity`, `operator_health_impact`) and does not include
per-pod or per-node dimensions.

The `virt-handler` runs as a DaemonSet on every node that can
schedule VMIs. It is responsible for domain lifecycle, network
configuration, and other node-level operations for virtual machine
instances.

In newer versions of OpenShift Virtualization, the alert expression is reworked
to
surface additional diagnostic labels (`pod`, `node`, `reason`) when
a container waiting reason is available. If your alert includes
these labels, see step 1 of the diagnosis below.

## Impact

Virtual machine instances (VMIs) on affected nodes cannot be managed
properly. New VMIs may not start on nodes without a running
`virt-handler`, and existing VMIs may not receive updates or clean
shutdowns.

## Diagnosis

1. **Check the alert labels**:

   If the alert includes a `reason` label (for example,
   `CrashLoopBackOff`, `ErrImagePull`, `ImagePullBackOff`), it
   directly identifies why `virt-handler` is down. The `pod` and
   `node` labels identify the affected pod and node. Skip to
   [Mitigation](#mitigation) for the matching root cause. If these
   labels are not present, continue with the steps below.

2. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A \
       -o custom-columns="":.metadata.namespace)"
   ```

3. Check the status of the `virt-handler` DaemonSet and pods:

   ```bash
   $ oc -n $NAMESPACE get daemonset virt-handler -o yaml
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-handler
   ```

4. Check DaemonSet events and pod status:

   ```bash
   $ oc -n $NAMESPACE describe daemonset virt-handler
   $ oc -n $NAMESPACE describe pod -l kubevirt.io=virt-handler
   ```

5. Check for node issues (for example, nodes not ready or taints):

   ```bash
   $ oc get nodes
   ```

6. If any `virt-handler` pod exists, review its logs:

   ```bash
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-handler \
       --previous
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-handler
   ```

## Mitigation

Try to identify the root cause and resolve the issue. Common
causes include:

- **CrashLoopBackOff**: The `virt-handler` container is crashing
  repeatedly. Check the pod logs for the root cause (panic, OOM,
  misconfiguration).
- **ErrImagePull / ImagePullBackOff**: The container image cannot
  be pulled. Verify the image reference, registry availability,
  and pull secrets.
- **Pods absent**: No `virt-handler` pods exist. Check whether the
  DaemonSet has been deleted, is not scheduling due to node taints,
  or is blocked by resource constraints.
- **Node issues**: Nodes may be in `NotReady` state, under resource
  pressure, or have taints that prevent the DaemonSet from
  scheduling.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
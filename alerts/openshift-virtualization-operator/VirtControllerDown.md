# VirtControllerDown

## Meaning

No running `virt-controller` pod has been detected for 10 minutes.

The alert expression evaluates
`cluster:kubevirt_virt_controller_pods_running:count == 0` with a
`for` duration of 10 minutes. The recording rule counts pods in
`Running` phase matching `virt-controller-.*`.

In newer versions of KubeVirt, the alert expression is reworked to
surface additional diagnostic labels (`pod`, `reason`) when a
container waiting reason is available. If your alert includes these
labels, see step 1 of the diagnosis below.

## Impact

Any actions related to virtual machine (VM) lifecycle management
fail. This notably includes launching a new virtual machine instance
(VMI) or shutting down an existing VMI.

## Diagnosis

1. **Check the alert labels**:

   If the alert includes a `reason` label (for example,
   `CrashLoopBackOff`, `ErrImagePull`, `ImagePullBackOff`), it
   directly identifies why `virt-controller` is down. The `pod`
   label identifies the affected pod. Skip to
   [Mitigation](#mitigation) for the matching root cause. If these
   labels are not present, continue with the steps below.

2. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A \
       -o custom-columns="":.metadata.namespace)"
   ```

3. Check the status of the `virt-controller` deployment:

   ```bash
   $ oc -n $NAMESPACE get deploy virt-controller -o yaml
   ```

4. Check the `virt-controller` deployment details for issues such
   as crashing pods or image pull failures:

   ```bash
   $ oc -n $NAMESPACE describe deploy virt-controller
   ```

5. Check the status of the `virt-controller` pods:

   ```bash
   $ oc -n $NAMESPACE get pods \
       -l kubevirt.io=virt-controller
   ```

6. Review the logs of the `virt-controller` pods:

   ```bash
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-controller \
       --previous
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-controller
   ```

7. Check for issues such as nodes in a `NotReady` state:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Try to identify the root cause and resolve the issue. Common
causes include:

- **CrashLoopBackOff**: The `virt-controller` container is
  crashing repeatedly. Check the pod logs for the root cause
  (panic, OOM, misconfiguration).
- **ErrImagePull / ImagePullBackOff**: The container image cannot
  be pulled. Verify the image reference, registry availability,
  and pull secrets.
- **Pods absent**: No `virt-controller` pods exist. Check whether
  the deployment has been scaled to zero, deleted, or blocked by
  resource constraints.
- **Node resource exhaustion**: Not enough memory or CPU on the
  cluster to schedule the pods.
- **Node issues**: Nodes may be in `NotReady` state or under
  resource pressure.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
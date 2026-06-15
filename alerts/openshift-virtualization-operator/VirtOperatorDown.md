# VirtOperatorDown

## Meaning

This alert fires when
`cluster:kubevirt_virt_operator_pods_running:count == 0` for 10
minutes, meaning no `virt-operator` pod in `Running` phase has been
detected.

The `virt-operator` is the first Operator to start in a cluster. Its
primary responsibilities include the following:

- Installing, live-updating, and live-upgrading a cluster
- Monitoring the life cycle of top-level controllers, such as
  `virt-controller`, `virt-handler`, `virt-launcher`, and managing
  their reconciliation
- Certain cluster-wide tasks, such as certificate rotation and
  infrastructure management

In newer versions of KubeVirt, the alert expression is reworked to
surface additional diagnostic labels (`pod`, `reason`) when a
container waiting reason is available. If your alert includes these
labels, see step 1 of the diagnosis below.

## Impact

This alert indicates a failure at the level of the cluster. Critical
cluster-wide management functionalities, such as certification
rotation, upgrade, and reconciliation of controllers, might not be
available.

The `virt-operator` is not directly responsible for virtual machines
(VMs) in the cluster. Therefore, its temporary unavailability does
not significantly affect VM workloads.

## Diagnosis

1. **Check the alert labels**:

   If the alert includes a `reason` label (for example,
   `CrashLoopBackOff`, `ErrImagePull`, `ImagePullBackOff`), it
   directly identifies why `virt-operator` is down. The `pod` label
   identifies the affected pod. Skip to
   [Mitigation](#mitigation) for the matching root cause. If these
   labels are not present, continue with the steps below.

2. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A \
       -o custom-columns="":.metadata.namespace)"
   ```

3. Check the status of the `virt-operator` deployment:

   ```bash
   $ oc -n $NAMESPACE get deploy virt-operator -o yaml
   ```

4. Obtain the details of the `virt-operator` deployment:

   ```bash
   $ oc -n $NAMESPACE describe deploy virt-operator
   ```

5. Check the status of the `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE get pods \
       -l kubevirt.io=virt-operator
   ```

6. Review the logs of the `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-operator \
       --previous
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-operator
   ```

7. Check for node issues, such as a `NotReady` state:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Try to identify the root cause and resolve the issue. Common
causes include:

- **CrashLoopBackOff**: The `virt-operator` container is crashing
  repeatedly. Check the pod logs for the root cause (panic, OOM,
  misconfiguration).
- **ErrImagePull / ImagePullBackOff**: The container image cannot
  be pulled. Verify the image reference, registry availability,
  and pull secrets.
- **Pods absent**: No `virt-operator` pods exist. Check whether the
  deployment has been scaled to zero, deleted, or blocked by
  resource constraints.
- **Node issues**: Nodes may be in `NotReady` state, under resource
  pressure, or have scheduling constraints that prevent the pods
  from running.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
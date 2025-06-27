# NoReadyVirtOperator

## Meaning

This alert fires when no `virt-operator` pod in a `Ready` state has been
detected for 10 minutes.

The `virt-operator` is the first Operator to start in a cluster. Its primary
responsibilities include the following:

- Installing, live-updating, and live-upgrading a cluster
- Monitoring the life cycle of top-level controllers, such as `virt-controller`,
`virt-handler`, `virt-launcher`, and managing their reconciliation
- Certain cluster-wide tasks, such as certificate rotation and infrastructure
management

The default deployment is two `virt-operator` pods.

## Impact

This alert indicates a cluster-level failure. Critical cluster management
functionalities, such as certification rotation, upgrade, and reconciliation of
controllers, might not be not available.

The `virt-operator` is not directly responsible for virtual machines in the
cluster. Therefore, its temporary unavailability does not significantly affect
custom workloads.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Obtain the name of the `virt-operator` deployment:

   ```bash
   $ oc -n $NAMESPACE get deploy virt-operator -o yaml
   ```

3. Generate the description of the `virt-operator` deployment:

   ```bash
   $ oc -n $NAMESPACE describe deploy virt-operator
   ```

4. Check for node issues, such as a `NotReady` state:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
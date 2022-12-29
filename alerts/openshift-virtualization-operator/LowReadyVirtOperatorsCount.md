# LowReadyVirtOperatorsCount

## Meaning

This alert fires when one or more `virt-operator` pods are running, but
none of these pods has been in a `Ready` state for the last 10 minutes.

The `virt-operator` is the first Operator to start in a cluster. The `virt-operator`
deployment has a default replica of two `virt-operator` pods.

Its primary responsibilities include the following:

- Installing, live-updating, and live-upgrading a cluster
- Monitoring the lifecycle of top-level controllers, such as `virt-controller`,
`virt-handler`, `virt-launcher`, and managing their reconciliation
- Certain cluster-wide tasks, such as certificate rotation and infrastructure
management

## Impact

A cluster-level failure might occur. Critical cluster-wide management
functionalities, such as certification rotation, upgrade, and reconciliation of
controllers, might become unavailable. Such a state also triggers the
`NoReadyVirtOperator` alert.

The `virt-operator` is not directly responsible for virtual machines (VMs)
in the cluster. Therefore, its temporary unavailability does not significantly
affect VM workloads.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A \
     -o custom-columns="":.metadata.namespace)"
   ```

2. Obtain the name of the `virt-operator` deployment:

   ```bash
   $ oc -n $NAMESPACE get deploy virt-operator -o yaml
   ```

3. Obtain the details of the `virt-operator` deployment:

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

# LowVirtOperatorCount

## Meaning

This alert fires when only one `virt-operator` pod in a `Ready` state has been
running for the last 60 minutes.

The `virt-operator` is the first Operator to start in a cluster. Its primary
responsibilities include the following:

- Installing, live-updating, and live-upgrading a cluster
- Monitoring the lifecycle of top-level controllers, such as `virt-controller`,
`virt-handler`, `virt-launcher`, and managing their reconciliation
- Certain cluster-wide tasks, such as certificate rotation and infrastructure
management

## Impact

The `virt-operator` cannot provide high availability (HA) for the deployment. HA
requires two or more `virt-operator` pods in a `Ready` state. The default
deployment is two pods.

The `virt-operator` is not directly responsible for virtual machines (VMs) in
the cluster. Therefore, its decreased availability does not significantly affect
VM workloads.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o jsonpath='{.items[].metadata.namespace}')"
   ```

2. Check the states of the `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-operator
   ```

3. Review the logs of the affected `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE logs <virt-operator>
   ```

4. Obtain the details of the affected `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE describe pod <virt-operator>
   ```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
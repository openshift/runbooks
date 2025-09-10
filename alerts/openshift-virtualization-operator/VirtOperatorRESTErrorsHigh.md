# VirtOperatorRESTErrorsHigh [Deprecated]

This alert is deprecated. You can safely ignore or silence it.
erator` pods
failed in the last 60 minutes. This usually indicates the `virt-operator` pods
cannot connect to the API server.

This error is frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the
case, check the metrics of the API server, and view its response times and
overall calls.

- The `virt-operator` pod cannot reach the API server. This is commonly caused
by DNS issues on the node and networking connectivity issues.

## Impact

Cluster-level actions, such as upgrading and controller reconciliation, might be
delayed.

However, customer workloads, such as virtual machines (VMs) and VM instances
(VMIs), are not likely to be affected.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-operator
   ```

3. Check the `virt-operator` logs for error messages when connecting to the API
server:

   ```bash
   $ oc -n $NAMESPACE logs <virt-operator>
   ```

4. Obtain the details of the `virt-operator` pod:

   ```bash
   $ oc -n $NAMESPACE describe pod <virt-operator>
   ```

## Mitigation

- If the `virt-operator` pod cannot connect to the API server, delete the pod to
force a restart:

  ```bash
  $ oc delete -n <install-namespace> <virt-operator>
  ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
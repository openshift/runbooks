# LowVirtAPICount

## Meaning

This alert fires when only one available `virt-api` pod is detected during a
60-minute period, although at least two nodes are available for scheduling.

## Impact

An API call outage might occur during node eviction because the `virt-api` pod
becomes a single point of failure.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace | tr -d '\n')"
   ```

2. Check the number of available `virt-api` pods:

   ```bash
   $ oc get deployment -n $NAMESPACE virt-api -o jsonpath='{.status.readyReplicas}'
   ```

3. Check the status of the `virt-api` deployment for error conditions:

   ```bash
   $ oc -n $NAMESPACE get deploy virt-api -o yaml
   ```

4. Check the nodes for issues such as nodes in a `NotReady` state:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Try to identify the root cause and to resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
# HPPOperatorDown

## Meaning

This alert fires when the hostpath provisioner (HPP) Operator is down.

The HPP Operator deploys and manages the HPP infrastructure components, such as
the daemon set that provisions hostpath volumes.

## Impact

The HPP components might fail to deploy or to remain in the required state. As a
result, the HPP installation might not work correctly in the cluster.

## Diagnosis

1. Configure the `HPP_NAMESPACE` environment variable:

   ```bash
   $ HPP_NAMESPACE="$(oc get deployment -A | grep hostpath-provisioner-operator | awk '{print $1}')"
   ```

2. Check whether the `hostpath-provisioner-operator` pod is currently running:

   ```bash
   $ oc -n $HPP_NAMESPACE get pods -l name=hostpath-provisioner-operator
   ```

3. Obtain the details of the `hostpath-provisioner-operator` pod:

   ```bash
   $ oc -n $HPP_NAMESPACE describe pods -l name=hostpath-provisioner-operator
   ```

4. Check the log of the `hostpath-provisioner-operator` pod for errors:

   ```bash
   $ oc -n $HPP_NAMESPACE logs -l name=hostpath-provisioner-operator
   ```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
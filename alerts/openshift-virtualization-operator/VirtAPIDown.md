# VirtAPIDown
<!-- Edited by apinnick, Nov 2022-->

## Meaning

This alert fires when all the API Server pods are down.

## Impact

KubeVirt objects cannot send API calls.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A \
     -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-api` pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-api
   ```

3. Check the status of the `virt-api` deployment:

   ```bash
   $ oc -n $NAMESPACE get deploy virt-api -o yaml
   ```

4. Check the `virt-api` deployment details for issues such as crashing pods or
image pull failures:

   ```bash
   $ oc -n $NAMESPACE describe deploy virt-api
   ```

5. Check for issues such as nodes in a `NotReady` state:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Try to identify the root cause and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.

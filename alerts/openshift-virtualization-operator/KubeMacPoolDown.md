# KubeMacPoolDown
<!-- Edited by apinnick, Oct. 2022-->

## Meaning

`KubeMacPool` is down. `KubeMacPool` is responsible for allocating MAC
addresses and preventing MAC address conflicts.

## Impact

If `KubeMacPool` is down, `VirtualMachine` objects cannot be created.

## Diagnosis

1. Set the `KMP_NAMESPACE` environment variable:

   ```bash
   $ export KMP_NAMESPACE="$(oc get pod -A --no-headers -l \
     control-plane=mac-controller-manager | awk '{print $1}')"
   ```

2. Set the `KMP_NAME` environment variable:

   ```bash
   $ export KMP_NAME="$(oc get pod -A --no-headers -l \
     control-plane=mac-controller-manager | awk '{print $2}')"
   ```

3. Obtain the `KubeMacPool-manager` pod details:

   ```bash
   $ oc describe pod -n $KMP_NAMESPACE $KMP_NAME
   ```

4. Check the `KubeMacPool-manager` logs for error messages:

   ```bash
   $ oc logs -n $KMP_NAMESPACE $KMP_NAME
   ```

## Mitigation

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.

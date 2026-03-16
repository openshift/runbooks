# LowReadyVirtAPICount

## Meaning

This alert fires when one or more `virt-api` pods are running, but not
all of them have been in a `Ready` state for the last 10 minutes.

The `virt-api` serves the OpenShift Virtualization API. The deployment
typically runs two
replicas for high-availability.

## Impact

Reduced capacity or redundancy for the OpenShift Virtualization API. If the
condition
persists, it can lead to the `NoReadyVirtAPI` alert and API unavailability.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o jsonpath='{.items[].metadata.namespace}')"
   ```

2. Check the status of the `virt-api` pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-api
   ```

3. Check the `virt-api` deployment and its events:

   ```bash
   $ oc -n $NAMESPACE describe deploy virt-api
   ```

4. Check pod readiness and conditions for non-ready pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-api -o wide
   $ oc -n $NAMESPACE describe pod -l kubevirt.io=virt-api
   ```

5. If pods are in `CrashLoopBackOff` or to inspect runtime failures, check
   `virt-api` pod logs and look for errors:

   ```bash
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-api
   ```

6. Check for node issues, such as a `NotReady` state:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Identify why some `virt-api` pods are not ready (e.g. failed readiness probe,
resource pressure, image pull issues) and resolve the underlying cause.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
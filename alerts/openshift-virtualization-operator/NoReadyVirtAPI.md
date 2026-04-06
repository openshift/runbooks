# NoReadyVirtAPI

## Meaning

This alert fires when no `virt-api` pod in a `Ready` state has been detected
for 10 minutes.

The `virt-api` serves the OpenShift Virtualization API.
Without a ready `virt-api`, API requests for virtual machines
and other KubeVirt resources cannot be served.

## Impact

KubeVirt API is effectively unavailable. Users and controllers cannot perform
API operations such as creating, updating, or deleting virtual machine
instances (VMIs) or other KubeVirt resources.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-api` pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-api
   ```

3. Check the `virt-api` deployment and events:

   ```bash
   $ oc -n $NAMESPACE describe deploy virt-api
   ```

4. Review logs of any `virt-api` pod that is running but not ready:

   ```bash
   $ oc -n $NAMESPACE logs <virt-api-pod-name> --previous
   $ oc -n $NAMESPACE logs <virt-api-pod-name>
   ```

5. Check for node issues:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Identify the root cause (e.g. all replicas crashing, readiness probe failures,
node or resource issues) and restore at least one ready `virt-api` pod.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
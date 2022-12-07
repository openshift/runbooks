# VirtHandlerDaemonSetRolloutFailing
<!-- Edited by apinnick, Nov 2022-->

## Meaning

The `virt-handler` daemon set has failed to deploy on one or more worker
nodes after 15 minutes.

## Impact

This alert is a warning. It does not indicate that all `virt-handler` daemon
sets have failed to deploy. Therefore, the normal lifecycle of virtual
machines is not affected unless the cluster is overloaded.

## Diagnosis

Identify worker nodes that do not have a running `virt-handler` pod:

1. Export the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A \
     -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-handler` pods to identify pods that have
not deployed:

   ```bash
   $ oc get pods -n $NAMESPACE -l=kubevirt.io=virt-handler
   ```

3. Obtain the name of the worker node of the `virt-handler` pod:

   ```bash
   $ oc -n $NAMESPACE get pod <virt-handler> -o jsonpath='{.spec.nodeName}'
   ```

## Mitigation

If the `virt-handler` pods failed to deploy because of insufficient resources,
you can delete other pods on the affected worker node.

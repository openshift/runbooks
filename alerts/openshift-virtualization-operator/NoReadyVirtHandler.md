# NoReadyVirtHandler

## Meaning

This alert fires when no `virt-handler` pod in a `Ready` state has been
detected for 10 minutes.

The `virt-handler` runs on every node that can schedule VMIs (as a
DaemonSet). It is responsible for domain lifecycle and node-level operations
for virtual machine instances.

## Impact

No node has a ready `virt-handler`. Virtual machine instances cannot be
properly managed: domain updates, migrations, and node-level operations will
fail or be delayed until at least one `virt-handler` becomes ready.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the status of the `virt-handler` DaemonSet and pods:

   ```bash
   $ oc -n $NAMESPACE get daemonset virt-handler
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-handler -o wide
   ```

3. Check DaemonSet and pod events:

   ```bash
   $ oc -n $NAMESPACE describe daemonset virt-handler
   $ oc -n $NAMESPACE describe pod -l kubevirt.io=virt-handler
   ```

4. Review logs of any running but not ready `virt-handler` pod:

   ```bash
   $ oc -n $NAMESPACE logs <virt-handler-pod-name> --previous
   $ oc -n $NAMESPACE logs <virt-handler-pod-name>
   ```

5. Check for cluster-wide node or scheduling issues:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Identify the root cause (e.g. DaemonSet not scheduling, all pods crashing or
failing readiness, node or image issues) and restore at least one ready
`virt-handler` pod on a schedulable node.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
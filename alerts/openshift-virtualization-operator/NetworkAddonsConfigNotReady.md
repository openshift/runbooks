# NetworkAddonsConfigNotReady

## Meaning

This alert fires when the `NetworkAddonsConfig` custom resource (CR) of the
Cluster Network Addons Operator (CNAO) is not ready.

CNAO deploys additional networking components on the cluster. This alert
indicates that one of the deployed components is not ready.

## Impact

Network functionality is affected.

## Diagnosis

1. Check the status conditions of the `NetworkAddonsConfig` CR to identify the
deployment or daemon set that is not ready:

   ```bash
   $ oc get networkaddonsconfig -o custom-columns=':.status.conditions[*].message' | tr -d '\n'
   ```

   Example output:

   ```text
   DaemonSet "openshift-cnv/macvtap-cni" update is being processed...
   ```

2. Check the component's pod for errors:

   ```bash
   $ oc -n openshift-cnv get daemonset <daemonset> -o yaml
   ```

3. Check the component's logs:

   ```bash
   $ oc -n openshift-cnv logs <pod>
   ```

4. Check the component's details for error conditions:

   ```bash
   $ oc -n openshift-cnv describe pod <pod>
   ```

## Mitigation

Try to identify the root cause and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
# NoLeadingVirtController

## Meaning

This alert fires when no `virt-controller` pod with a leader lease has been
detected for 10 minutes, although the `virt-controller` pods are in a `Ready`
state. The alert indicates that no leader pod is available.

The `virt-controller` is responsible for monitoring the custom resource
definitions of virtual machine instances (VMIs) and managing the associated
pods. It creates pods for VMIs and manages their lifecycle.

The `virt-controller` deployment has a default replica of 2 pods, with one pod
holding a leader lease.

## Impact

This alert indicates a failure at the level of the cluster. As a result,
critical cluster-wide virtualization functionalities, such as VMI lifecycle
management, might not be available.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns=NAMESPACE:.metadata.namespace --no-headers | head -1)"
   ```

2. Obtain the status of the `virt-controller` pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-controller
   ```

3. Check the `virt-controller` pod logs to determine the leader status:

   ```bash
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-controller | grep -i lead
   ```

   Leader pod example:

   ```text
   I1130 12:15:18.635452       1 leaderelection.go:243] attempting to acquire leader lease <namespace>/virt-controller...
   I1130 12:15:19.216582       1 leaderelection.go:253] successfully acquired lease <namespace>/virt-controller
   STARTING controllers with following threads : node 3, vmi 3, replicaset 3, vm 3, migration 3, evacuation 3, disruptionBudget 3
   ```

   Non-leader pod example:

   ```text
   I1130 12:15:20.533792       1 leaderelection.go:243] attempting to acquire leader lease <namespace>/virt-controller...
   ```

4. Check the leader lease:

   ```bash
   $ oc -n $NAMESPACE get lease virt-controller -o yaml
   ```

   A healthy cluster has a `holderIdentity` set to the name of the leading pod.

5. Obtain the details of the affected `virt-controller` pods:

   ```bash
   $ oc -n $NAMESPACE describe pod <virt-controller>
   ```

6. Check for admission webhook rejections that block lease create or update
   operations:

   ```bash
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-controller \
       | grep -iE 'webhook|admission|forbidden|denied'
   ```

   If logs show a validating webhook rejecting `leases` operations, identify
   the webhook configuration and review its logs.

   ```bash
   $ oc get validatingwebhookconfiguration
   $ oc get mutatingwebhookconfiguration
   ```

## Mitigation

Based on the diagnosis results, apply the remediation that matches the root
cause:

1. **Stale or corrupted lease**: If the `virt-controller` lease exists but
   `holderIdentity` references a pod that is not running, or if pod logs show
   repeated lease acquisition failures, delete the lease to allow re-election:

   ```bash
   $ oc -n $NAMESPACE delete lease virt-controller
   ```

2. **RBAC permissions**: Verify that the `kubevirt-controller` service account
   can create and update the lease object in the OpenShift Virtualization
namespace:

   ```bash
   $ oc auth can-i create leases \
       --as=system:serviceaccount:$NAMESPACE:kubevirt-controller \
       -n $NAMESPACE
   $ oc auth can-i update leases \
       --as=system:serviceaccount:$NAMESPACE:kubevirt-controller \
       -n $NAMESPACE
   ```

   If either command returns `no`, inspect the `kubevirt-controller` Role and
   RoleBinding in `$NAMESPACE`. The Role must grant
   `coordination.k8s.io` `leases` verbs `get`, `list`, `watch`, `create`,
   `update`, and `patch`.

3. **Admission webhook rejection**: If `virt-controller` logs show admission
   webhook errors when creating or updating the `virt-controller` lease,
   identify the rejecting webhook and review its configuration and logs.
   Update or remove the webhook policy so `coordination.k8s.io` `leases`
   operations are allowed in the OpenShift Virtualization namespace.

4. **Resource constraints**: If pods show `OOMKilled`, high CPU throttling,
   or slow startup in `describe` output, increase the `virt-controller`
   deployment resource requests and limits, or relieve node resource pressure
   so lease acquisition can complete within the 15-second lease duration.

5. **Force leader re-election**: If the lease is healthy but no pod acquires
   leadership, restart the `virt-controller` deployment to trigger a new
   election:

   ```bash
   $ oc -n $NAMESPACE rollout restart deployment virt-controller
   ```

   Alternatively, delete the affected pods so the remaining replicas can
   acquire the lease:

   ```bash
   $ oc -n $NAMESPACE delete pod -l kubevirt.io=virt-controller
   ```

6. **Validate resolution**: Confirm that one pod holds the lease and the alert
   clears:

   ```bash
   $ oc -n $NAMESPACE get lease virt-controller \
       -o jsonpath='{.spec.holderIdentity}{"\n"}'
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-controller \
       | grep 'successfully acquired lease'
   ```

   The `holderIdentity` must match a running `virt-controller` pod name.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
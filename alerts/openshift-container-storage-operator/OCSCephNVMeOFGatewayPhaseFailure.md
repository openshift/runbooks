# OCSCephNVMeOFGatewayPhaseFailure

## Meaning

A `CephNVMeOFGateway` custom resource has been in a `Failed` or
`ReconcileFailed` phase for more than 10 minutes. `CephNVMeOFGateway`
deploys a Ceph NVMe-oF gateway, which exposes RBD images over the NVMe-oF
protocol to NVMe-oF-capable clients, created/managed by the
`StorageCluster` via the Rook operator.

Note: the Rook NVMe-oF controller currently only sets `status.phase` to
empty (on first creation) or `Ready` (on reconcile success); reconcile
errors are reported only as Kubernetes `Events`/operator log lines, not as
a `Failed`/`ReconcileFailed` phase. This alert is therefore defense-in-depth
for a future/alternate Rook behavior and will rarely fire today. If the
NVMe-oF gateway is unhealthy but this alert is not firing, rely on the
pod-level checks below instead of `status.phase`.

## Impact

NVMe-oF clients connected through the affected gateway can lose access to
their RBD volumes, and new NVMe-oF subsystem/namespace provisioning
through this gateway may fail. This feature is only enabled on Fusion Data
Foundation (FDF); it is not part of standard ODF deployments.

## Diagnosis

### Check the CephNVMeOFGateway status and events

```bash
oc get cephnvmeofgateway -n openshift-storage
oc describe cephnvmeofgateway <gateway_name> -n openshift-storage
```

### Check the NVMe-oF gateway pods

```bash
oc get pods -n openshift-storage -l app=rook-ceph-nvmeof
oc logs -n openshift-storage -l app=rook-ceph-nvmeof
```

Look for `CrashLoopBackOff`/restarts, or NVMe-oF subsystem/target errors
in the pod logs.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <gateway_name>
```

### Check overall Ceph cluster health

See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status`
and `ceph health detail`.

## Mitigation

1. Resolve any underlying Ceph cluster health issues identified above.
2. If the gateway pod(s) are crash-looping, check the pod logs for a
   configuration error in `spec` (for example, an invalid port or
   `nvmeOfConfig` setting) and correct it via the `StorageCluster` CR.
3. Once the root cause is resolved, the resource reconciles back to
   `Ready` automatically; no manual phase reset is needed.
4. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.

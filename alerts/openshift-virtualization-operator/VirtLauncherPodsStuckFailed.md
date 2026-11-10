# VirtLauncherPodsStuckFailed

## Meaning

A large number of virt-launcher Pods remain in the `Failed` phase.

Condition:
`sum(kube_pod_status_phase{phase="Failed", pod=~"virt-launcher-.*"}) >= 200`
holds for 10 minutes.
Virt-launcher Pods host virtual machine (VM) workloads and mass failures
can indicate migration loops, image/network/storage issues, or
control-plane regressions.

## Impact

VMs and the cluster control plane might be affected:

- API server and etcd pressure (large object lists/watches, increased latency)
- Controller and scheduler slowdown (reconciliation over huge Pod sets)
- Monitoring cardinality spikes (kube-state-metrics and Prometheus load)
- Operational churn (re-creation loops, CNI and storage attach/detach)
- Triage noise and SLO risk (timeouts on list operations, noisy dashboards)

## Diagnosis

1. Confirm scope and distribution:
    ```promql
    sum(kube_pod_status_phase{phase="Failed", pod=~"virt-launcher-.*"})
    ```
    ```promql
    count by (namespace) (
      max by (namespace, pod) (kube_pod_status_phase{phase="Failed", pod=~"virt-launcher-.*"}) == 1
    )
    ```
    ```promql
    count by (node) (
      (max by (namespace, pod) (kube_pod_status_phase{phase="Failed", pod=~"virt-launcher-.*"}) == 1)
      * on(namespace, pod) group_left(node)
        max by (namespace, pod, node) (kube_pod_info{pod=~"virt-launcher-.*", node!=""})
    )
    ```
    The `max by` aggregations collapse duplicate series for each Pod, including
    `kube_pod_info` series that differ only by `pod_ip`. Matching on both
    `namespace` and `pod` keeps Pods with the same name in different namespaces
    separate.

    ```promql
    topk(5, count by (reason) (kube_pod_container_status_last_terminated_reason{pod=~"virt-launcher-.*"} == 1))
    ```

2. Sample failed Pods and events:
    ```bash
    # List a few failed virt-launcher pods cluster-wide
    $ oc get pods -A -l kubevirt.io=virt-launcher --field-selector=status.phase=Failed --no-headers | head -n 20
    ```
    ```bash
    # Inspect events for a representative pod (image/CNI/storage/useful errors)
    $ oc -n <namespace> describe pod <virt-launcher-pod> | sed -n '/Events/,$p'
    ```

3. Check for migration storms:
    ```bash
    $ oc get vmim -A
    ```

4. In the control plane and component logs, look for spikes/errors:
    ```bash
    NAMESPACE="$(oc get kubevirt -A -o jsonpath='{.items[0].metadata.namespace}')"
    ```
    ```bash
    $ oc -n "$NAMESPACE" logs -l kubevirt.io=virt-controller --tail=200
    ```
    ```bash
    $ oc -n "$NAMESPACE" logs -l kubevirt.io=virt-handler --tail=200
    ```

5. Check the infrastructure for common causes:

- Image pulls: registry reachability/credentials; `ImagePullBackOff` events
- Network: CNI errors/timeouts in Pod events and node logs
- Storage: volume attach/mount errors in Pod events and CSI logs

## Mitigation

1. Reduce the blast radius:

    - Migration loop: cancel in-flight migrations (scope as needed).
      ```bash
      $ oc delete vmim -A
      ```

    - Coordinate with noisy tenants; pause offending workloads if necessary.

2. Clean up failed Pods (relieves API/etcd and monitoring):

    This deletes Pods across every namespace, so preview the Pods that
    match the selectors first. The listing includes the namespace of
    each Pod, so you can confirm that only `Failed` virt-launcher Pods
    match, and not running VM workloads:

    ```bash
    $ oc get pods -A -l kubevirt.io=virt-launcher \
        --field-selector=status.phase=Failed
    ```

    Only after verifying that listing, delete the Pods with the same
    selectors:

    ```bash
    $ oc delete pods -A -l kubevirt.io=virt-launcher \
        --field-selector=status.phase=Failed
    ```

3. Resolve root cause:

    - Image issues: fix registry access, credentials, or tags; re-run
      affected workloads.
    - Network/CNI: fix CNI/data-plane errors; confirm new Pods start cleanly.
    - Storage: resolve attach/mount failures; verify PVC/VolumeSnapshot health.
    - OpenShift Virtualization regression: roll forward/back to a known-good
version and re-try.

4. Validate resolution (alert clears):
    ```promql
    sum(kube_pod_status_phase{phase="Failed", pod=~"virt-launcher-.*"})
    ```

Ensure the failed count drops and stays below threshold and that new
virt-launcher Pods start successfully and VMIs are healthy.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
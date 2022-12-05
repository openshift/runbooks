# KubeVirtComponentExceedsRequestedCPU
<!-- Edited by apinnick, Nov 2022-->

## Meaning

This alert fires when a component's CPU usage exceeds the requested limit.

## Impact

Usage of CPU resources is not optimal and the node might be overloaded.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. Check the component's CPU request limit:

   ```bash
   $ oc -n $NAMESPACE get deployment <component> -o yaml | grep requests: -A 2
   ```

3. Check the actual CPU usage by using a PromQL query:

   ```text
   node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate
   {namespace="$NAMESPACE",container="<component>"}
   ```

See the
[Prometheus documentation](https://prometheus.io/docs/prometheus/latest/querying/basics/)
for more information.

## Mitigation

Update the CPU request limit in the `HCO` custom resource.
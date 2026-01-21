# NetObservNoFlows

## Meaning

The `NetObservNoFlows` alert is an operational alert that triggers when
Network Observability is not receiving any network flow data for a certain
period. This indicates that the flow collection pipeline has stopped working,
preventing network traffic visibility.

Unlike other Network Observability alerts, `NetObservNoFlows` does not have
multiple variants based on grouping or severity levels. It is a single
critical alert that indicates a complete failure in flow collection.

This alert monitors the `flowlogs-pipeline` component which processes and
exports network flows. When no flows are received, it means either:

- The eBPF agents on nodes are not collecting flows
- `flowlogs-pipeline` is not receiving data from agents
- `flowlogs-pipeline` has crashed or stopped processing
- `flowlogs-pipeline` fails to report its metrics

**Note:** This is an operational alert that monitors the health of Network
Observability itself, not the health of cluster network traffic.

### Configuration limitations

Like other Network Observability operational alerts, `NetObservNoFlows` cannot be configured, other than being disabled:

- It cannot be converted to recording mode - it is always an alert
- It does not support thresholds - it fires when Loki write errors occur
  (flow rate == 0) consistently after 10 minutes
- It does not support grouping - it is a global cluster-wide operational
  alert
- It cannot have variants - there is only one alert instance

The alert triggers with this hardcoded PromQL expression:

```promql
sum(rate(netobserv_ingest_flows_processed[1m])) == 0
```

This design is intentional because no flows being received indicates a
**complete failure** of Network Observability and should always generate
alerts.

### Disable this alert entirely

We do not recommend disabling this alert as it indicates a critical failure.
However, if needed:

```bash
oc edit flowcollector cluster
```

Add NetObservNoFlows to the disableAlerts list:

```yaml
spec:
  processor:
    metrics:
      disableAlerts:
      - NetObservNoFlows
```

For more information on Network Observability operational health and
troubleshooting, see the
[Network Observability documentation](https://docs.openshift.com/container-platform/latest/network_observability/troubleshooting-network-observability.html).

## Impact

When no flows are being observed, the impact is:

- Complete loss of network traffic visibility in the cluster
- Network Observability Console shows no data or stale data
- Network Health dashboard shows no metrics
- All network-based alerts will not fire (blind to network issues)
- Unable to troubleshoot network problems
- Loss of network traffic audit trail
- Compliance violations if network monitoring is required

This is a critical operational issue that requires immediate attention, as it
means Network Observability is completely non-functional.

## Diagnosis


If you are running in Kafka mode (`spec.deploymentModel: Kafka` in `FlowCollector`), make sure that Kafka is up and running, and correctly configured in `FlowCollector` via `spec.kafka`.

Check if all components are running and don't show any critical errors in their logs:

**eBPF agents:**

```bash
oc get pods -n netobserv-privileged
oc logs -n netobserv-privileged <POD>
```

**flowlogs-pipeline:**

```bash
oc get pods -n netobserv -l app=flowlogs-pipeline
oc logs -n netobserv <POD>
```

Check for any existing network policies, both in agent and `flowlogs-pipeline` namespaces (by default, `netobserv` and `netobserv-privileged`), that could be blocking the traffic.

Absence of flows might also be related to a misconfiguration in `FlowCollector`: the eBPF agents can be configured to include or exclude network interfaces from their listeners (respectively via `spec.agent.ebpf.interfaces` and `spec.agent.ebpf.excludeInterfaces`), or to define custom filtering rules (via `spec.agent.ebpf.flowFilter`). You should review these configurations and make sure they could not lead to exclude all the traffic.

If everything looks good, another possibility is that `flowlogs-pipeline` metrics are not correctly pulled by the Cluster Monitoring components (Prometheus). From the OpenShift Console, navigate to _Observe_ > _Metrics_, and search for any metric prefixed with `netobserv_`: if there is none, this is likely an issue in the monitoring configuration.

For additional troubleshooting resources, refer to the documentation links in
the Mitigation section below.

## Mitigation

If you find that a network policy is blocking the traffic:
Network policies are automatically installed when `spec.networkPolicy.enable` is `true` in `FlowCollector`. It should not block the traffic between agents and `flowlogs-pipeline`, however if you find that it is, you can disable it entirely and write your own policy instead. We would also kindly ask you to report the issue to the maintainers team.

If you find that all metrics prefixed with `netobserv_` are missing, review the monitoring configuration:

```bash
oc get servicemonitors -n netobserv
```

It should show several Service Monitors, including `flowlogs-pipeline-monitor`. If this is not the case, check the Network Observability operator logs (in namespace `openshift-netobserv-operator`) for any relevant errors. You can also search for the label `job=flowlogs-pipeline-prom` in the OpenShift Console under _Observe_ > _Targets_, to check the pull status of the metrics. More info on troubleshooting monitoring issues is available in the [Monitoring stack documentation](https://docs.redhat.com/en/documentation/monitoring_stack_for_red_hat_openshift/latest/html/troubleshooting_monitoring_issues/index)

For other mitigation strategies and solutions, refer to the [Troubleshooting Network
Observability](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/network_observability/installing-troubleshooting)
documentation.

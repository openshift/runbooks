# ClusterMonitoringOperatorDeprecatedConfig

## Meaning

This alert fires when a deprecated config is used in the
`openshift-monitoring/cluster-monitoring-config` config map.

## Impact

Avoid using a deprecated config because the config has no effect, and doing so
might cause the `Upgradeable` condition for the Cluster Monitoring Operator
to become `False` in a future OpenShift Container Platform release.

## Diagnosis

- Check the __Description__ and __Summary__ annotations of the alert to identify
the deprecated config as shown in the following example:

  __Description__

  `The configuration field k8sPrometheusAdapter in
  openshift-monitoring/cluster-monitoring-config was deprecated in version 4.16
  and has no effect`.

  __Summary__

  `Cluster Monitoring Operator is being used with a deprecated configuration.`

## Mitigation

* For the `k8sPrometheusAdapter.dedicatedServiceMonitors`
field, you can remove the block. For more information, see
`Monitoring deprecated and removed features` under
[Deprecated and removed features](https://docs.openshift.com/container-platform/4.16/release_notes/ocp-4-16-release-notes.html#ocp-4-16-deprecated-removed-features_release-notes).

* For the other `k8sPrometheusAdapter` fields, see `Monitoring deprecated and
removed features` under [Deprecated and removed features](https://docs.openshift.com/container-platform/4.16/release_notes/ocp-4-16-release-notes.html#ocp-4-16-deprecated-removed-features_release-notes).
You might need to migrate some of the fields under [metricsServer](https://docs.openshift.com/container-platform/latest/observability/monitoring/config-map-reference-for-the-cluster-monitoring-operator.html#metricsserverconfig).

The alert resolves itself when the deprecated config is not used.

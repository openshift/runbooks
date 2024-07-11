# ClusterMonitoringOperatorDeprecatedConfig

## Meaning

This alert fires when a deprecated config is used in the config map `openshift-monitoring/cluster-monitoring-config`.

## Impact

Using a deprecated config is not recommended because the config has no effect and
it may cause the Cluster Monitoring operator `Upgradeable` condition to become
`False` in a future OCP release.

## Diagnosis

- Check the __Description__ and __Summary__ annotations of the alert to identify
the
  deprecated config as in the following example:

  __Description__

  `The configuration field k8sPrometheusAdapter in
  openshift-monitoring/cluster-monitoring-config was deprecated in 4.16 and has
  no effect`.

  __Summary__

  `Cluster Monitoring Operator is being used with deprecated configuration.`

## Mitigation

* For the `k8sPrometheusAdapter.dedicatedServiceMonitors`
field, you can remove the block, see the release notes entry
`Monitoring deprecated and removed features` under [deprecated-removed-features](https://docs.openshift.com/container-platform/4.16/release_notes/ocp-4-16-release-notes.html#ocp-4-16-deprecated-removed-features_release-notes)

* For the other `k8sPrometheusAdapter` fields, check the
release notes entry `Monitoring deprecated and removed features` under [deprecated-removed-features](https://docs.openshift.com/container-platform/4.16/release_notes/ocp-4-16-release-notes.html#ocp-4-16-deprecated-removed-features_release-notes),
some of the fields may need to be migrated under [metricsServer](https://docs.openshift.com/container-platform/latest/observability/monitoring/config-map-reference-for-the-cluster-monitoring-operator.html#metricsserverconfig)

The alert resolves itself when the deprecated config is not used.

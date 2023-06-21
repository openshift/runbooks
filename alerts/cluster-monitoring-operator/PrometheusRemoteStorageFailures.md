# PrometheusRemoteStorageFailures

## Meaning

The `PrometheusRemoteStorageFailures` alert triggers when failures to send
samples to remote storage ave been constantly increasing for more than 15
minutes for either platform Prometheus pods or user-defined monitoring
Prometheus pods.

## Impact

Prometheus samples in remote storage might be missing. Depending on the
`remote_write` pipeline configuration, Prometheus memory usage might increase
while pending samples are queued.

## Diagnosis

* Check the `namespace` label in the alert message to determine if the alert was
  triggered for the instance of Prometheus used for default cluster monitoring
  or for the instance that monitors user-defined projects. The `namespace` value
  indicates the Prometheus instance: `openshift-monitoring`for default
  monitoring and `openshift-user-workload-monitoring` for user-defined
  monitoring.

* Review the logs for the affected Prometheus instance:

  ```console
  $ NAMESPACE='<value of namespace label from alert>'

  $ oc -n $NAMESPACE logs -l 'app.kubernetes.io/component=prometheus'
  level=error ... msg="Failed to send batch, retrying" ...
  ```

* Review the Prometheus logs and the remote storage logs.

## Mitigation

This alert fires when Prometheus has an issue communicating with the remote
system. The cause can be on either the Prometheus side or the remote side.

Common issues that can cause the alert to fire include the following:

* **Cause**: Failure to authenticate to the remote storage system
  * **Mitigation**: Verify that the authentication parameters in the Cluster
    Monitoring (CMO) or in the user workload are correct.

* **Cause**: Requests hitting rate-limits
  * **Mitigation**: Tweak the queue configuration in the CMO and user-workload
    config maps and/or limit the number of samples being sent.

* **Cause**: 5xx HTTP errors
  * **Mitigation**: Review the logs on for the remote storage system.

If the logs indicate a configuration error, troubleshoot the issue. Note
that the issue might be related to general networking issues or a bad configuration.

The cause might also be that the amount of data snet to the remote system is too
high for a given network speed. If so, minimize transfers by limiting which
metrics are sent to remote storage.

Additionally, you can check the `cluster-network-operator` configuration to help
debug possible networking issues.

# AuditLogError

## Meaning

This alert is triggered when an API Server instance in the cluster is unable
to write audit logs. It fires when there's any errors, which is calculated
by checking the error rate with `apiserver_audit_error_total` and `apiserver_audit_event_total`.

There might be many causes to this:

* This might have been caused by the node that host's that API Server instance
  running out of disk space.

* A malicious actor could be tampering with the audit log files or directory
  permissions

* The API server might be encountering an unexpected error.

## Impact

When there are errors writing audit logs, security events will not be logged
by that specific API Server instance. Security Incident Response teams use
these audit logs, amongst other artifacts, to determine the impact of
security breaches or events. Without these logs, it becomes very difficult
to assess a situation and do appropriate root cause analysis in such incidents.

However, this is not detrimental to the cluster's availability.

## Diagnosis

Verify if there are other alerts being triggered, e.g. the
`NodeFilesystemFillingUp` alert. This might indicate what's causing this error.

The metric `apiserver_audit_error_total` will only show up on instances
that are experiencing errors. There will be appropriate labels that indicate
the API Server type and the specific instance that's affected.

With this information, gather the runtime logs from that specific api server
pod, and verify if the logs indicate an unexpected error.

From the `instance` label, it's possible to determine what node is hosting
the aforementioned affected API Server.

Log into the appropriate node and Verify that the relevant audit log file
permissions are what's expected:
Owned by the `root` user and with a mode of `0600`. Make sure as well that
there aren't unexpected attributes for the files, such as immutability or append-only.

While logged into the appropriate node, also verify that the relevant audit
log directory permissions are what's expected:
Owned by the `root` user and with a mode of `0700`.

If you suspect tampering is happening, contact your incident response team.

## Mitigation

The appropriate mitigation will be very different depending on the organization
and the compliance requirements. A FedRAMP moderate deployment might need to
isolate the node and investigate, while deployment with more strict compliance
requirements would need to snapshot and shut down the system immediately.
A more usual deployment might just need to investigate, and since the causes
could be many. Regardless, investigate the deployment as described in the
diagnosis, and contact the incident response team in your organization if
necessary.

# KubeVirtDeprecatedAPIRequested

## Meaning

This alert fires when a deprecated `KubeVirt` API is used.

## Impact

Using a deprecated API is not recommended because the request will
fail when the API is removed in a future release.

## Diagnosis

- Check the __Description__ and __Summary__ sections of the alert to identify
the
  deprecated API as in the following example:

  __Description__

  `Detected requests to the deprecated virtualmachines.kubevirt.io/v1alpha3 API.`

  __Summary__

  `2 requests were detected in the last 10 minutes.`

## Mitigation

Use fully supported APIs. The alert resolves itself after 10 minutes if the
deprecated API is not used.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
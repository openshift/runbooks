# KubeVirtDeprecatedAPIsRequested

## Meaning

This alert fires when a deprecated KubeVirt API is used.

## Impact

Using a deprecated API is not recommended because the request will
fail when the API is removed in a future release.

## Diagnosis

- Check the `description` and `summary` alert annotations to identify the
deprecated API:

  ```text
                description: "Detected requests to the deprecated virtualmachines.kubevirt.io/v1alpha3 API."
                summary: "2 requests were detected in the last 10 minutes."
  ```

## Mitigation

Use fully supported APIs. The alert resolves itself after 10 minutes if the deprecated
API is not used.


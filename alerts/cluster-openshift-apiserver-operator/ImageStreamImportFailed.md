# ImageStreamImportFailed

## Meaning
The `ImageStreamImportFailed` alert fires when one or more ImageStreams
fail to import new image tags within the last 10 minutes.
This alert is based on the `openshift_imagestreamcontroller_error_count` metric,
which increments whenever the ImageStream controller reports an import failure (ImportSuccess=False)
for a registry.

## Impact
When the ImageStream controller fails to import image tags, workloads that depend
on ImageStreamTags might continue using stale or outdated image versions.

Without this alert:

* Cluster administrators and application owners may remain unaware that
image imports are failing.
* Pods referencing ImageStreamTags (for example, openshift/tools:latest) could
continue running an older image even after a cluster upgrade
(e.g., from version vA to vB).

This alert provides early visibility so administrators can quickly identify,
investigate, and resolve import issues.

## Diagnosis
Investigate failing ImageStreams:

1. List all ImageStreams with import failures:

    ```console
    $ oc get imagestreams --all-namespaces -o wide | grep -v True
    ```
    This filters out ImageStreams whose
    `.status.tags[*].conditions[].type=ImportSuccess` is False or missing.

2. Check status of a specific ImageStream:

    ```console
    $ oc get -o json imagestream <name> -n <namespace> | jq .status
    ```
    The `.status.tags[].conditions` will show specific import failure reasons.

## Mitigation

* Inspect the failure reason in `.status.tags[].conditions[].message`
and `.status.tags[].conditions[].reason`

    Typical causes include:
  * Invalid or expired image registry credentials.
  * Network or proxy configuration issues preventing access to the external registry.
  * Typo or deprecation in the referenced image tag or pullspec.
  * Missing trusted CA for HTTPS registries

* Re-trigger the import manually to verify the fix:

    ```console
    $ oc import-image <imagestream-name> -n <namespace> --confirm
    ```

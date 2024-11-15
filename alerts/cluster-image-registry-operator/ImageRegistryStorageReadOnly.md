# ImageRegistryStorageReadOnly

## Meaning

This alert is triggered when the image registry pvc storage is read-only.

## Impact

Builds and direct pushes to the image registry will fail.

Images imported via Image Streams or `oc import-image` will be served directly
from the upstream registry where they are being imported from, but caching to
local storage will fail imperceptibly to end-users, unless the upstream
registry becomes unavailable.

## Mitigation

Ensure the PVC volume is not mounted to the image registry pod in `readOnly` mode.

The underlying storage volume might also be configured as read-only. Ensure it's
writable.

Specific mitigation steps for this problem will differ depending on the underlying
storage solution. Please refer to your specific storage solution's documentation.

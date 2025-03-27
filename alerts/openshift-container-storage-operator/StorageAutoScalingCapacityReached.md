# StorageAutoScalingCapacityReached

## Meaning

The Storage Capacity Limit is the maximum storage that can be provisioned for
the system. `StorageAutoScalingCapacityReached` alert is fired when the complete
storage space has been utilized and hence additional storage cannot be allocated
for the device class in the particular storage cluster.

## Impact

The storage cannot be auto-scaled at the moment.

## Diagnosis

Retrieve status of `storageCapacityLimitReached` from the StorageAutoScaler CR.
The `resource_name` and `resource_namespace` would be specified in the alert itself.

```bash
oc get storageautoscaler <resource_name> -n <resource_namespace> -o json | jq -r '.status.storageCapacityLimitReached'
```

If status of `storageCapacityLimitReached` is `true` then maximum capacity of
the storage has been utilized.

## Mitigation

To address the issues mentioned above, you can apply one of the following solutions
to scale the storage:

- Increase storage capacity: Run the command below with the updated `storageCapacityLimit`:
  ```bash
  oc patch storageautoscaler <resource_name> -n <resource_namespace> --type=merge --patch '{"spec":{"storageCapacityLimit":<storageCapacityLimit>}}'
  ```
  Note: The `status.error.message` in the StorageAutoScaler CR will show the
  minimum required increase for `storageCapacityLimit`.Ensure that the specified
  capacity is at least that value or higher.

- If auto scaling is no longer desired, disable the feature. Additional storage
can be added when needed via the ODF console.To disable auto scaling:
  ```bash
  oc delete storageautoscaler <resource_name> -n <resource_namespace>
  ```
  


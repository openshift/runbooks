# StorageAutoScalerCRIsInvalid

## Meaning

This alert is fired when the StorageAutoScaler encounters an unsupported or invalid
configuration of the storage class. The specific conditions that trigger this
alert are detailed in the `Diagnosis` section below.

## Impact

The storage cannot be scaled at the moment.

## Diagnosis

The `StorageAutoScaler` status field contains a message field that explains why
the configuration is invalid. To view the status:
```bash
oc get storageautoscaler -n openshift-storage -o jsonpath='{.items[*].status}'
```

Common error messages are listed below:

- LSO storage class detected: The storage class has no provisioner specified,
making it an LSO (Local Storage Operator) storage class. Auto-scaling is not
supported for LSO storage classes. The error message will appear as follows:
```bash
"storage class has provisioner as no-provisioner, which is an lso storageclass, autoscaler does not support lso storageclass, delete the autoStorageScaler cr as scaling is not supported"
```

- Duplicate CR detected: While multiple StorageAutoScaler CRs can exist within
a namespace, they cannot have both the same device class and storage cluster name
simultaneously. The presence of duplicate CRs with identical device class and
storage cluster name results in this alert. The error message will appear as follows:
```bash
"duplicate cr detected, more than one storage autoscaler present with same device class and same storage cluster name, names are <CR1> and <CR2> delete any one of the autoStorageScaler cr"
```

- Lean Storage Profile detected: The storage autoscaler does not support lean
storage profiles. This will trigger an error message as follows:
```bash
"storage profile is lean, autoscaler does not support lean storage profile, delete the autoStorageScaler cr as scaling is not supported"
```

## Mitigation

Delete the appropriate StorageAutoScaler resources.
```bash
oc delete storageautoscaler <scaler_name> -n <resource_namespace>
```
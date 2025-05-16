# StorageAutoScalerCRIsInvalid

## Meaning

As a predefined rule for automatic storage scaling, only one StorageAutoScaler
CR should exist within a namespace. This alert is triggered when multiple
StorageAutoScaler CRs are present, which disrupts the scaling process.

## Impact

The storage cannot be scaled at the moment.

## Diagnosis

Retrieve all `StorageAutoScaler` CRs within the respective namespace. If more
than one is found, this indicates the source of the issue.

```bash
oc get storageautoscaler -n <resource_namespace> -o json | jq -r '[.items[] | {name: .metadata.name, storageCluster: .spec.storageCluster.name, deviceClass: .spec.deviceClass }] | group_by(.storageCluster+"-"+ .deviceClass) | map(select(length>1))[] | .[] | .name'
```

## Mitigation

Delete all StorageAutoScaler resources except the one that best fits the requirements.
```bash
oc delete storageautoscaler <scaler_name> -n <resource_namespace>
```
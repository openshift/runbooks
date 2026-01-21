# NooBaaSystemCapacityWarning100

## Meaning

NooBaa system capacity is calculated by aggregating the total storage capacity
of the backingstore resources. It does not include namespacestore resources.
This alert is triggered when the total used space of the backingstore resources
reaches 100% of the total capacity.

## Impact

The used storage of all backingstore resources is at 100% of the total capacity.
The system is out of space and cannot accept new writes.

## Diagnosis

Depending on the type of the backingstore resource, the diagnosis may vary.
- If the backingstore resource is based on a public cloud provider,
the free space is always considered to be 1 Petabyte.
In that case the alert is not relevant.
- If the backingstore resource is based on a compatible S3 service,
check the usage and available space of the target bucket used by backingstore.
- If the backingstore resource is a PV-Pool backingstore, check the usage and
available space of PVs used by the pv-pool pods.

## Mitigation

- For public cloud providers, there is nothing to do.
- For compatible S3 services, check the usage and available
space of the target bucket used by backingstore. If the usage is approaching
the capacity, consider increasing the capacity of the bucket.
- For PV-Pool backingstores, check the usage and available space of PVs used
by the pv-pool pods. If the usage is approaching the capacity, consider
increasing the number of volumes in the backingstore spec to add new PVs to
the pool. Notice that the number of volumes cannot be reduced later.

# ObcQuotaObjectsAlert

## Meaning

An ObjectBucketClaim object has crossed 80% of the size limit set by the quota(objects)
and will become read-only on reaching the quota limit.

## Impact

OBC has reached 80% of it's limit and soon will get exhausted once reaching
the quota limit.

## Diagnosis

Alert message will clearly indicate which OBC is being filled up fast.
Look at the deployments attached to the OBC and
see what all apps are using/filling-up the OBC.

## Mitigation

We can increase quota option on OBC by using the `maxObjects` and
`maxSize` options in the ObjectBucketClaim CRD

```yaml
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim
metadata:
  name: <obc_name>
  namespace: <namespace>
spec:
  bucketName: <name_of_the_backend_bucket>
  storageClassName: <name_of_storage_class>
  additionalConfig:
    maxObjects: "1000" # sets limit on no of objects this obc can hold
    maxSize: "2G" # sets max limit for the size of data this obc can hold
```

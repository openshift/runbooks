# ObcQuotaBytesExhausedAlert

## Meaning

This is the next stage once we have reached [ObcQuotaObjectsAlert](ObcQuotaObjectsAlert.md).
ObjectBucketClaim has crossed the limit set by the quota(bytes) and will be
read-only now. Increase the quota in the OBC custom resource immediately.

## Impact

OBC has exhausted and reached it's limit.

## Diagnosis

Alert message will clearly indicate which OBC has reached the quota bytes limit.
Look at the deployments attached to the OBC and see what all apps are
using/filling-up the OBC.

## Mitigation

Need to increase the quota limit immediately for the ObjectBucketClaim
custom resource. We can set quota option on OBC by using the `maxObjects`
and `maxSize` options in the ObjectBucketClaim CRD

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

